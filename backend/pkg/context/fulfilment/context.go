// Code generated by simba. DO NOT EDIT.

package fulfilment

import (
	"fmt"
	"log"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"

	checkout "github.com/cod1ng-earth/event-web-store/backend/pkg/context/checkout/public"

	public "github.com/cod1ng-earth/event-web-store/backend/pkg/context/fulfilment/public"
)

const (
	Topic     = "fulfilment_internal"
	Partition = 0
)

type context struct {
	doneCh chan struct{}

	doneChCheckout chan struct{}

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

		doneChCheckout: make(chan struct{}, 1),

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

	updateModel(c, msg, m)

}

func (c *context) bridgeCheckout() {

	c.AwaitLastOffset()

	model, free := c.read()
	checkoutOffset := model.getCheckoutOffset()
	free()
	partition, err := c.consumer.ConsumePartition(checkout.Topic, 0, checkoutOffset)
	if err != nil {
		log.Panicf("failed to setup kafka partition: %s", err)
	}

	log.Printf("starting bridge %v", Topic)

	for {
		select {
		case err := <-partition.Errors():
			log.Printf("failure from kafka consumer: %s", err)

		case msg := <-partition.Messages():
			//			log.Printf("recieved message with offset %v", msg.Offset)

			cc := checkout.TopicMessage{}
			err := proto.Unmarshal(msg.Value, &cc)
			if err != nil {
				log.Fatalf("failed to unmarshal kafka massage %s/%d: %v", Topic, msg.Offset, err)
			}

			switch x := cc.GetMessages().(type) {

			case *checkout.TopicMessage_OrderCreated:
				if err := translateCheckoutOrderCreated(c, model, msg.Offset, cc.GetOrderCreated()); err != nil {
					log.Fatalf("failed to translate kafka message $bridge.Name/%v: %s", msg.Offset, err)
				}

			case nil:
				panic(fmt.Sprintf("context message is empty"))

			default:
				panic(fmt.Sprintf("unexpected type %T in oneof", x))
			}

		}
	}
}

func (c *context) Start() {

	log.Printf("starting context %v", Topic)

	writes := make(chan *sarama.ConsumerMessage, 32768)
	go c.updateLoop(writes)

	partition, err := c.consumer.ConsumePartition(Topic, 0, 0)
	if err != nil {
		log.Panicf("failed to setup kafka partition: %s", err)
	}

	go c.bridgeCheckout()

	for {
		select {
		case err := <-partition.Errors():
			log.Printf("failure from kafka consumer: %s", err)

		case msg := <-partition.Messages():
			//			log.Printf("recieved message with offset %v", msg.Offset)
			writes <- msg

		case <-c.doneCh:
			log.Print("interrupt is detected")

			c.doneChCheckout <- struct{}{}

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

func updateModel(c *context, msg *sarama.ConsumerMessage, model *model) error {
	cc := TopicMessage{}
	err := proto.Unmarshal(msg.Value, &cc)
	if err != nil {
		return fmt.Errorf("failed to unmarshal kafka massage %s/%d: %v", Topic, msg.Offset, err)
	}

	switch x := cc.GetMessages().(type) {

	case *TopicMessage_StockCorrected:
		fact := cc.GetStockCorrected()
		err = updateModelStockCorrected(model, msg.Offset, fact)
		if err != nil {
			return fmt.Errorf("failed to update kafka massage %s/%d: %v", Topic, msg.Offset, err)
		}

		err = publishStockCorrected(c, msg.Offset, fact)
		if err != nil {
			return fmt.Errorf("failed to publish kafka massage %s/%d: %v", Topic, msg.Offset, err)
		}

	case nil:
		panic(fmt.Sprintf("context message is empty"))

	default:
		panic(fmt.Sprintf("unexpected type %T in oneof", x))
	}

	return nil
}

type asyncProducer struct {
	producer sarama.AsyncProducer
	wg       sync.WaitGroup
}

func (c *context) newSyncProducer(f func(error)) (asyncProducer, error) {
	producer, err := sarama.NewAsyncProducerFromClient(c.client)
	if err != nil {
		return asyncProducer{}, fmt.Errorf("failed to create async producer: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for err := range producer.Errors() {
			f(err)
		}
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		for range producer.Successes() {
		}
		wg.Done()
	}()

	return asyncProducer{
		producer: producer,
		wg:       wg,
	}, nil
}

func (p *asyncProducer) Close() {
	p.producer.AsyncClose()
	p.wg.Wait()
}

func (c *context) logStockCorrected(msg *StockCorrected) (int32, int64, error) {

	topicMsg := &TopicMessage{
		Messages: &TopicMessage_StockCorrected{
			StockCorrected: msg,
		},
	}

	bytes, err := proto.Marshal(topicMsg)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to serialize stockCorrected change massage: %v", err)
	}

	producerMsg := &sarama.ProducerMessage{
		Topic: Topic,
		Value: sarama.ByteEncoder(bytes),
	}
	return c.producer.SendMessage(producerMsg)
}

func (p asyncProducer) logStockCorrected(msg *StockCorrected) error {

	topicMsg := &TopicMessage{
		Messages: &TopicMessage_StockCorrected{
			StockCorrected: msg,
		},
	}

	bytes, err := proto.Marshal(topicMsg)
	if err != nil {
		return fmt.Errorf("failed to serialize stockCorrected change massage: %v", err)
	}

	producerMsg := &sarama.ProducerMessage{
		Topic: Topic,
		Value: sarama.ByteEncoder(bytes),
	}
	p.producer.Input() <- producerMsg

	return nil
}

func (c *context) logPublicStockCorrected(msg *public.StockCorrected) (int32, int64, error) {

	topicMsg := &public.TopicMessage{
		Messages: &public.TopicMessage_StockCorrected{
			StockCorrected: msg,
		},
	}

	bytes, err := proto.Marshal(topicMsg)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to serialize stockCorrected change massage: %v", err)
	}

	producerMsg := &sarama.ProducerMessage{
		Topic: public.Topic,
		Value: sarama.ByteEncoder(bytes),
	}
	return c.producer.SendMessage(producerMsg)
}
