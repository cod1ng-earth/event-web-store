// Code generated by simba. DO NOT EDIT.

package pim

import (
	"fmt"
	"log"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
)

const (
	Topic     = "pim"
	Partition = 0
)

type context struct {
	doneCh chan struct{}

	client   sarama.Client
	consumer sarama.Consumer
	producer sarama.SyncProducer

	batchOffset int64

	model *model
	lock  *sync.RWMutex

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

	batchOffset, err := client.GetOffset(Topic, Partition, sarama.OffsetNewest)
	if err != nil {
		log.Printf("failed to get last offset for topic %v partition %v", Topic, Partition)
	}
	batchOffset--

	context := context{
		doneCh: make(chan struct{}, 1),

		client:   client,
		consumer: consumer,
		producer: producer,

		batchOffset: batchOffset,

		model: newModel(),
		lock:  &sync.RWMutex{},

		offset:        0,
		offsetChanged: sync.NewCond(&sync.Mutex{}),
	}
	return context
}

func (c *context) Stop() {
	c.doneCh <- struct{}{}
}

func (c *context) await(offset int64) {
	if offset == -1 {
		return
	}
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
	c.await(c.batchOffset)
}

func (c *context) updateLoop(writes <-chan *sarama.ConsumerMessage) {

	for {

		for msg := range writes {
			applyChange(msg, c.model, c)
		}

	}
}

func applyChange(msg *sarama.ConsumerMessage, m *model, c *context) {

	//	log.Printf("applying message with offset %v", msg.Offset)

	c.lock.Lock()
	defer c.lock.Unlock()
	defer func() {
		c.offset = msg.Offset
		c.offsetChanged.Broadcast()
	}()

	updateModel(msg, m)

}

func (c *context) Start() {

	log.Printf("starting context %v", Topic)

	writes := make(chan *sarama.ConsumerMessage, 32768)
	go c.updateLoop(writes)

	partition, err := c.consumer.ConsumePartition(Topic, 0, 0)
	if err != nil {
		log.Panicf("failed to setup kafka partition: %s", err)
	}

	for {
		select {
		case err := <-partition.Errors():
			log.Printf("failure from kafka consumer: %s", err)

		case msg := <-partition.Messages():
			//			log.Printf("recieved message with offset %v", msg.Offset)
			writes <- msg

		case <-c.doneCh:
			log.Print("interrupt is detected")

			if err := partition.Close(); err != nil {
				log.Panicf("failed to close kafka partition: %s", err)
			}
			close(writes)
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

	c.lock.RLock()
	return c.model, c.lock.RUnlock

}

func updateModel(msg *sarama.ConsumerMessage, model *model) error {
	cc := PimMessages{}
	err := proto.Unmarshal(msg.Value, &cc)
	if err != nil {
		return fmt.Errorf("failed to unmarshal kafka massage %s/%d: %v", Topic, msg.Offset, err)
	}

	switch x := cc.GetPimMessage().(type) {

	case *PimMessages_Product:
		return updateModelProduct(model, msg.Offset, cc.GetProduct())

	case nil:
		panic(fmt.Sprintf("context message is empty"))

	default:
		panic(fmt.Sprintf("unexpected type %T in oneof", x))
	}
}

func (c *context) logProduct(logMsg *Product) (int32, int64, error) {

	//log.Printf("logProduct");

	change := &PimMessages{
		PimMessage: &PimMessages_Product{
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
