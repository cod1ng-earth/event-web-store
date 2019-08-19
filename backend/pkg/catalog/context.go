package catalog

import (
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Shopify/sarama"
)

const (
	Topic     = "products"
	Partition = 0
)

type context struct {
	doneCh    chan struct{}
	client    sarama.Client
	consumer  sarama.Consumer
	partition sarama.PartitionConsumer

	batchOffset int64

	readerAChanged *sync.Cond
	readerBChanged *sync.Cond
	aIsReading     bool
	readersA       int32
	readersB       int32
	modelA         *model
	modelB         *model

	writes     chan *sarama.ConsumerMessage
	writesRedo chan *sarama.ConsumerMessage

	offset        int64
	offsetChanged *sync.Cond
}

func NewContext(brokers *[]string, cfg *sarama.Config) context {

	client, err := sarama.NewClient(*brokers, cfg)
	if err != nil {
		log.Panicf("failed to setup kafka client: %s", err)
	}
	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		log.Panicf("failed to setup kafka consumer: %s", err)
	}
	partition, err := consumer.ConsumePartition(Topic, 0, 0)
	if err != nil {
		log.Panicf("failed to setup kafka partition: %s", err)
	}

	batchOffset, err := client.GetOffset(Topic, Partition, sarama.OffsetNewest)
	if err != nil {
		log.Printf("failed to get last offset for topic %v partition %v", Topic, Partition)
	}
	batchOffset--

	//	batchOffset -= 50000

	context := context{
		doneCh:    make(chan struct{}, 1),
		client:    client,
		consumer:  consumer,
		partition: partition,

		batchOffset: batchOffset,

		aIsReading:     true,
		modelA:         newModel(),
		modelB:         newModel(),
		offsetChanged:  sync.NewCond(&sync.Mutex{}),
		readerAChanged: sync.NewCond(&sync.Mutex{}),
		readerBChanged: sync.NewCond(&sync.Mutex{}),
		writes:         make(chan *sarama.ConsumerMessage, 32768),
		writesRedo:     make(chan *sarama.ConsumerMessage, 32768),
	}
	return context
}

func (c *context) Stop() {
	c.doneCh <- struct{}{}
}

func (c *context) AwaitLastOffset() {
	log.Printf("waiting for offset %v", c.batchOffset)
	c.offsetChanged.L.Lock()
	for c.offset < c.batchOffset {
		c.offsetChanged.Wait()
	}
	c.offsetChanged.L.Unlock()
	log.Printf("watermark %v for topic %v reached", c.batchOffset, Topic)

	m, cl := c.read()
	defer cl()
	log.Printf("products count: %v", len(m.products))

}

func (c *context) updateLoop() {

	var offset int64
	var start time.Time

	for {
		writeDelay := 0 * time.Second

		model := c.modelA
		readers := &c.readersA
		waiter := c.readerAChanged
		if c.aIsReading {
			model = c.modelB
			readers = &c.readersB
			waiter = c.readerBChanged
		}

		waiter.L.Lock()
		if atomic.LoadInt32(readers) != 0 {
			waiter.Wait()
		}
		waiter.L.Unlock()

		for len(c.writesRedo) > 0 {
			msg, ok := <-c.writesRedo
			if ok {
				applyChange(msg, model, c)
			}
		}

		msg, ok := <-c.writes
		if !ok {
			return
		}
		c.writesRedo <- msg
		start = time.Now()
		applyChange(msg, model, c)
		offset = msg.Offset
		writeDelay += time.Since(start)

		for writeDelay < 10*time.Millisecond && len(c.writes) > 0 && len(c.writesRedo) < 32768 {
			msg, ok := <-c.writes
			if !ok {
				return
			}
			c.writesRedo <- msg
			start := time.Now()
			applyChange(msg, model, c)
			offset = msg.Offset
			writeDelay += time.Since(start)
		}

		start = time.Now()
		if msg.Offset >= c.batchOffset {
			modelRealtimeFinalizer(model)
		}
		finalizeDelay := time.Since(start)
		if finalizeDelay > 200*time.Millisecond {
			log.Printf("finalize took %v to execute", finalizeDelay)
		}

		c.aIsReading = !c.aIsReading
		c.offset = offset
		c.offsetChanged.Broadcast()
	}

}

func applyChange(msg *sarama.ConsumerMessage, m *model, c *context) {
	if msg.Offset < c.batchOffset {
		modelBatchUpdater(msg, m)
	} else if msg.Offset == c.batchOffset {
		modelBatchUpdater(msg, m)
		modelBatchFinalizer(m)
	} else {
		modelRealtimeUpdater(msg, m)
	}
}

// Start listens for events from kafka
func (c *context) Start() {

	log.Printf("starting context %v", Topic)

	go c.updateLoop()

	for {
		select {
		case err := <-c.partition.Errors():
			log.Printf("failure from kafka consumer: %s", err)

		case msg := <-c.partition.Messages():
			//			log.Printf("processing %v:%v:%v", Topic, msg.Partition, msg.Offset)
			c.writes <- msg

		case <-c.doneCh:
			close(c.writes)
			log.Print("interrupt is detected")
			if err := c.partition.Close(); err != nil {
				log.Panicf("failed to close kafka partition: %s", err)
			}
			if err := c.consumer.Close(); err != nil {
				log.Panicf("failed to close kafka consumer: %s", err)
			}
			if err := c.client.Close(); err != nil {
				log.Panicf("failed to close kafka client: %s", err)
			}
			return
		}
	}
}

func (c *context) read() (*model, func()) {
	//	log.Printf("A: %v, B: %v, Offset: %v", atomic.LoadInt32(&c.readersA), atomic.LoadInt32(&c.readersB), c.offset)

	atomic.AddInt32(&c.readersA, 1)
	atomic.AddInt32(&c.readersB, 1)

	if c.aIsReading {
		atomic.AddInt32(&c.readersB, -1)
		c.readerBChanged.Signal()
		return c.modelA, func() {
			atomic.AddInt32(&c.readersA, -1)
			c.readerAChanged.Signal()
		}
	}

	atomic.AddInt32(&c.readersA, -1)
	c.readerAChanged.Signal()
	return c.modelB, func() {
		atomic.AddInt32(&c.readersB, -1)
		c.readerBChanged.Signal()
	}

}
