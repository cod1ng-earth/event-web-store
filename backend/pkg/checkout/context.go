package checkout

import (
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Shopify/sarama"
)

const (
	Topic     = "checkout"
	Partition = 0
)

type context struct {
	doneCh    chan struct{}
	client    sarama.Client
	consumer  sarama.Consumer
	partition sarama.PartitionConsumer

	batchOffset int64


	model  *model
	lock   sync.RWMutex
	writes chan *sarama.ConsumerMessage


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

	context := context{
		doneCh:    make(chan struct{}, 1),
		client:    client,
		consumer:  consumer,
		partition: partition,

		batchOffset: batchOffset,

	
		model:  newModel(),
		lock:   &sync.RWMutex{},
		writes: make(chan *sarama.ConsumerMessage, 32768),
	

		offset:        0,
		offsetChanged: sync.NewCond(&sync.Mutex{}),
	}
	return context
}

func (c *context) Stop() {
	c.doneCh <- struct{}{}
}

func (c *context) AwaitLastOffset() {
	c.offsetChanged.L.Lock()
	for c.offset < c.batchOffset {
		c.offsetChanged.Wait()
	}
	c.offsetChanged.L.Unlock()
}

func (c *context) updateLoop() {

	var offset int64
	var start time.Time

	for {
	
		applyChange(msg, c.model, c)
	
	}

}

func applyChange(msg *sarama.ConsumerMessage, m *model, c *context) {


	c.lock.Lock()
	defer c.lock.Unlock()



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

	if msg.Offset < c.batchOffset {
		c.AwaitLastOffset()
	}


	c.lock.RLock()
	return model, c.lock.Unlock

}


func batchUpdateModel(msg *sarama.ConsumerMessage, model *model) error {
	cc := CheckoutMessages{}
	err := proto.Unmarshal(msg.Value, &cc)
	if err != nil {
		return fmt.Errorf("failed to unmarshal kafka massage %s/%d: %v", Topic, msg.Offset, err)
	}

	defer func() { offset = msg.Offset }()

	switch x := cc.GetCheckoutMessage().(type) {

	
	case *CheckoutMessages_ChangeProductQuantity:
		if err := batchUpdateModelChangeProductQuantity(cc.GetChangeProductQuantity(), offset); err != nil {
			return err
		}
	
	case *CheckoutMessages_Stock:
		if err := batchUpdateModelStock(cc.GetStock(), offset); err != nil {
			return err
		}
	
	case *CheckoutMessages_Product:
		if err := batchUpdateModelProduct(cc.GetProduct(), offset); err != nil {
			return err
		}
	
	case *CheckoutMessages_CartOrder:
		if err := batchUpdateModelCartOrder(cc.GetCartOrder(), offset); err != nil {
			return err
		}
	

	case nil:
		panic(fmt.Sprintf("context message is empty"))

	default:
		panic(fmt.Sprintf("unexpected type %T in oneof", x))
	}
	return nil
}


func updateModel(msg *sarama.ConsumerMessage, model *model) error {
	cc := CheckoutMessages{}
	err := proto.Unmarshal(msg.Value, &cc)
	if err != nil {
		return fmt.Errorf("failed to unmarshal kafka massage %s/%d: %v", Topic, msg.Offset, err)
	}

	defer func() { offset = msg.Offset }()

	switch x := cc.GetCheckoutMessage().(type) {

	
	case *CheckoutMessages_ChangeProductQuantity:
		if err := updateModelChangeProductQuantity(cc.GetChangeProductQuantity(), offset); err != nil {
			return err
		}
	
	case *CheckoutMessages_Stock:
		if err := updateModelStock(cc.GetStock(), offset); err != nil {
			return err
		}
	
	case *CheckoutMessages_Product:
		if err := updateModelProduct(cc.GetProduct(), offset); err != nil {
			return err
		}
	
	case *CheckoutMessages_CartOrder:
		if err := updateModelCartOrder(cc.GetCartOrder(), offset); err != nil {
			return err
		}
	

	case nil:
		panic(fmt.Sprintf("checkout context message is empty"))

	default:
		panic(fmt.Sprintf("unexpected type %T in oneof", x))
	}
	return nil
}
