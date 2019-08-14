package products

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
	locklessModel locklessModel

	doneCh    chan struct{}
	client    sarama.Client
	consumer  sarama.Consumer
	partition sarama.PartitionConsumer

	batchUpdates      func(*sarama.ConsumerMessage, locklessModel) error
	batchFinalizer    func(*model)
	realtimeUpdates   func(*sarama.ConsumerMessage, locklessModel) error
	realtimeFinalizer func(*model)

	realtimeOffset int64
}

type locklessModel struct {
	readerAChanged *sync.Cond
	readerBChanged *sync.Cond
	aIsReading     bool
	readersA       int32
	readersB       int32
	modelA         *model
	modelB         *model

	doneCh     chan struct{}
	writes     chan func(*model) int64
	writesRedo chan func(*model) int64

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

	realtimeOffset, err := client.GetOffset(Topic, Partition, sarama.OffsetNewest)
	if err != nil {
		log.Printf("failed to get last offset for topic %v partition %v", Topic, Partition)
	}

	context := context{
		locklessModel: locklessModel{
			aIsReading:     true,
			modelA:         newModel(),
			modelB:         newModel(),
			offsetChanged:  sync.NewCond(&sync.Mutex{}),
			readerAChanged: sync.NewCond(&sync.Mutex{}),
			readerBChanged: sync.NewCond(&sync.Mutex{}),
			writes:         make(chan func(*model) int64, 32768),
			writesRedo:     make(chan func(*model) int64, 32768),
		},

		doneCh:            make(chan struct{}, 1),
		client:            client,
		consumer:          consumer,
		partition:         partition,
		batchUpdates:      modelBatchUpdates,
		batchFinalizer:    modelBatchFinalizer,
		realtimeUpdates:   modelRealtimeUpdates,
		realtimeFinalizer: modelRealtimeFinalizer,
		realtimeOffset:    realtimeOffset,
	}
	return context
}

func (l *locklessModel) Stop() {
	close(l.writes)
}

func (c *context) Stop() {
	c.doneCh <- struct{}{}
}

func (c *context) AwaitLastOffset() {
	log.Printf("waiting for offset %v", c.realtimeOffset-1)
	c.locklessModel.offsetChanged.L.Lock()
	for c.locklessModel.offset < c.realtimeOffset-1 {
		c.locklessModel.offsetChanged.Wait()
	}
	c.locklessModel.offsetChanged.L.Unlock()
	log.Printf("watermark %v for topic %v reached", c.realtimeOffset-1, Topic)
}

func (m *locklessModel) Start() {

	var offset int64
	var start time.Time

	for {
		writeDelay := 0 * time.Second

		model := m.modelA
		readers := &m.readersA
		waiter := m.readerAChanged
		if m.aIsReading {
			model = m.modelB
			readers = &m.readersB
			waiter = m.readerBChanged
		}

		waiter.L.Lock()
		if atomic.LoadInt32(readers) != 0 {
			waiter.Wait()
		}
		waiter.L.Unlock()

		for len(m.writesRedo) > 0 {
			w, ok := <-m.writesRedo
			if ok {
				w(model)
			}
		}

		w, ok := <-m.writes
		if !ok {
			return
		}
		m.writesRedo <- w
		start = time.Now()
		offset = w(model)
		writeDelay += time.Since(start)

		for (writeDelay < 10*time.Millisecond) && len(m.writes) > 0 && len(m.writesRedo) < 32768 {
			w, ok := <-m.writes
			if !ok {
				return
			}
			m.writesRedo <- w
			start := time.Now()
			offset = w(model)
			writeDelay += time.Since(start)
		}

		start = time.Now()
		if msg.Offset < c.realtimeOffset {
			modelBatchFinalizer(model)
		} else {
			modelUpdateFinalizer(model)
		}
		finalizeDelay := time.Since(start)
		if finalizeDelay > 200*time.Millisecond {
			log.Printf("finalize took %v to execute", finalizeDelay)
		}

		m.aIsReading = !m.aIsReading
		m.offset = offset
		m.offsetChanged.Broadcast()
	}

}

// Start listens for events from kafka
func (c *context) Start() {

	log.Printf("starting context %v", Topic)

	go c.locklessModel.Start()

	for {
		select {
		case err := <-c.partition.Errors():
			log.Printf("failure from kafka consumer: %s", err)

		case msg := <-c.partition.Messages():
			log.Printf("processing %v:%v:%v", Topic, msg.Partition, msg.Offset)
			if msg.Offset < c.realtimeOffset {
				err := c.batchUpdates(msg, c.locklessModel)
				if err != nil {
					log.Panicf("processing kafka message failed: %s", err)
					c.Stop()
				}
			} else {
				err := c.realtimeUpdates(msg, c.locklessModel)
				if err != nil {
					log.Panicf("processing kafka message failed: %s", err)
					c.Stop()
				}
			}

		case <-c.doneCh:
			c.locklessModel.Stop()
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

func (c *locklessModel) read() (*model, func()) {
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
