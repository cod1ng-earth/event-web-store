package catalog

import (
	"fmt"
	"log"
	"sync"

	"sync/atomic"
	"time"

	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
)

const (
	Topic     = "catalog"
	Partition = 0
)

type context struct {
	doneCh    chan struct{}
	client    sarama.Client
	consumer  sarama.Consumer
	producer  sarama.SyncProducer
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
	producer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		log.Panicf("failed to setup kafka producer: %s", err)
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

	context := context{
		doneCh:    make(chan struct{}, 1),
		client:    client,
		consumer:  consumer,
		producer:  producer,
		partition: partition,

		batchOffset: batchOffset,

		readerAChanged: sync.NewCond(&sync.Mutex{}),
		readerBChanged: sync.NewCond(&sync.Mutex{}),
		aIsReading:     true,
		readersA:       0,
		readersB:       0,
		modelA:         newModel(),
		modelB:         newModel(),

		writes:     make(chan *sarama.ConsumerMessage, 32768),
		writesRedo: make(chan *sarama.ConsumerMessage, 32768),

		offset:        0,
		offsetChanged: sync.NewCond(&sync.Mutex{}),
	}
	return context
}

func (c *context) Stop() {
	c.doneCh <- struct{}{}
}

func (c *context) await(offset int64) {
	if c.offset >= offset {
		return
	}
	c.offsetChanged.L.Lock()
	for c.offset < offset {
		c.offsetChanged.Wait()
	}
	c.offsetChanged.L.Unlock()
}

func (c *context) AwaitLastOffset() {
	c.offsetChanged.L.Lock()
	for c.offset < c.batchOffset {
		c.offsetChanged.Wait()
	}
	c.offsetChanged.L.Unlock()
}

func (c *context) updateLoop() {

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
		start := time.Now()
		applyChange(msg, model, c)
		offset := msg.Offset
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

		c.aIsReading = !c.aIsReading
		c.offset = offset
		c.offsetChanged.Broadcast()

	}
}

func applyChange(msg *sarama.ConsumerMessage, m *model, c *context) {

	//	log.Printf("applying message with offset %v", msg.Offset)

	if msg.Offset < c.batchOffset {
		batchUpdateModel(msg, m)
	} else if msg.Offset == c.batchOffset {
		batchUpdateModel(msg, m)
		batchFinalizeModel(m)
	} else {
		updateModel(msg, m)
	}

}

func (c *context) Start() {

	log.Printf("starting context %v", Topic)

	go c.updateLoop()

	for {
		select {
		case err := <-c.partition.Errors():
			log.Printf("failure from kafka consumer: %s", err)

		case msg := <-c.partition.Messages():
			//			log.Printf("recieved message with offset %v", msg.Offset)
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
			if err := c.producer.Close(); err != nil {
				log.Panicf("failed to close kafka producer: %s", err)
			}
			if err := c.client.Close(); err != nil {
				log.Panicf("failed to close kafka client: %s", err)
			}
			return
		}
	}
}

func (c *context) read() (*model, func()) {

	c.offsetChanged.L.Lock()
	for c.offset < c.batchOffset {
		c.offsetChanged.Wait()
	}
	c.offsetChanged.L.Unlock()

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

func batchUpdateModel(msg *sarama.ConsumerMessage, model *model) error {
	cc := CatalogMessages{}
	err := proto.Unmarshal(msg.Value, &cc)
	if err != nil {
		return fmt.Errorf("failed to unmarshal kafka massage %s/%d: %v", Topic, msg.Offset, err)
	}

	switch x := cc.GetCatalogMessage().(type) {

	case *CatalogMessages_Product:
		return batchUpdateModelProduct(model, msg.Offset, cc.GetProduct())

	case nil:
		panic(fmt.Sprintf("context message is empty"))

	default:
		panic(fmt.Sprintf("unexpected type %T in oneof", x))
	}
	return nil
}

func updateModel(msg *sarama.ConsumerMessage, model *model) error {
	cc := CatalogMessages{}
	err := proto.Unmarshal(msg.Value, &cc)
	if err != nil {
		return fmt.Errorf("failed to unmarshal kafka massage %s/%d: %v", Topic, msg.Offset, err)
	}

	switch x := cc.GetCatalogMessage().(type) {

	case *CatalogMessages_Product:
		return updateModelProduct(model, msg.Offset, cc.GetProduct())

	case nil:
		panic(fmt.Sprintf("context message is empty"))

	default:
		panic(fmt.Sprintf("unexpected type %T in oneof", x))
	}
	return nil
}

func (c *context) logProduct(logMsg *Product) (int32, int64, error) {

	change := &CatalogMessages{
		CatalogMessage: &CatalogMessages_Product{
			Product: logMsg,
		},
	}

	bytes, err := proto.Marshal(change)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to serialize cart change massage: %v", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: Topic,
		Value: sarama.ByteEncoder(bytes),
	}
	return c.producer.SendMessage(msg)
}
