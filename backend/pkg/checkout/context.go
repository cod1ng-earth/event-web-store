package checkout

import (
	"fmt"
	"log"
	"sync"

	"github.com/cod1ng-earth/event-web-store/backend/pkg/simba"
	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
)

const (
	Topic = "checkout"
)

var (
	mut           sync.RWMutex
	offset        int64
	offsetChanged *sync.Cond

	producer sarama.SyncProducer

	products map[string]*Product
	stock    map[string]int64
	carts    map[string]map[string]int64
	orders   map[string]map[string]int64
)

func init() {
	offsetChanged = sync.NewCond(&sync.Mutex{})

	products = map[string]*Product{}
	stock = map[string]int64{}
	carts = make(map[string]map[string]int64)
	orders = make(map[string]map[string]int64)
}

func StartContext(brokers *[]string, cfg *sarama.Config) func() {

	consumer, err := sarama.NewConsumer(*brokers, cfg)
	if err != nil {
		log.Panicf("failed to setup kafka consumer: %s", err)
	}
	partition, err := consumer.ConsumePartition(Topic, 0, 0)
	if err != nil {
		log.Panicf("failed to setup kafka patition: %s", err)
	}

	producer, err = sarama.NewSyncProducer(*brokers, cfg)
	if err != nil {
		log.Panicf("failed to setup the kafka producer: %s", err)
	}

	agent := simba.NewConsumer(consumer, partition, checkoutProcessor)
	go agent.Start()

	return func() {
		agent.Stop()
		if err := producer.Close(); err != nil {
			log.Printf("failed to close the kafka producer: %s", err)
		}
	}
}

func checkoutProcessor(msg *sarama.ConsumerMessage) error {
	cc := CheckoutContext{}
	err := proto.Unmarshal(msg.Value, &cc)
	if err != nil {
		return fmt.Errorf("failed to unmarshal kafka cart massage %d: %v", msg.Offset, err)
	}

	defer func() { offset = msg.Offset }()

	switch x := cc.GetCheckoutContextMsg().(type) {

	case *CheckoutContext_ChangeProductQuantity:
		if err := cartProcessor(cc.GetChangeProductQuantity(), offset); err != nil {
			return err
		}

	case *CheckoutContext_Stock:
		if err := stockProcessor(cc.GetStock()); err != nil {
			return err
		}

	case *CheckoutContext_Product:
		if err := productsProcessor(cc.GetProduct()); err != nil {
			return err
		}

	case *CheckoutContext_CartOrder:
		ordersProcessor(cc.GetCartOrder(), offset)

	case nil:
		panic(fmt.Sprintf("checkout context message is empty"))

	default:
		panic(fmt.Sprintf("unexpected type %T in oneof", x))
	}
	return nil
}
