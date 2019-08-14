package checkout

import (
	"fmt"
	"log"
	"sync"

	"git.votum-media.net/event-web-store/event-web-store/backend/pkg/pb"
	"git.votum-media.net/event-web-store/event-web-store/backend/pkg/simba"
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

	products map[string]*pb.Product
	stock    map[string]int64
	carts    map[string]map[string]int64
	orders   map[string]map[string]int64
)

func init() {
	offsetChanged = sync.NewCond(&sync.Mutex{})

	products = map[string]*pb.Product{}
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
	cc := pb.CheckoutContext{}
	err := proto.Unmarshal(msg.Value, &cc)
	if err != nil {
		return fmt.Errorf("failed to unmarshal kafka cart massage %d: %v", msg.Offset, err)
	}

	defer func() { offset = msg.Offset }()

	switch x := cc.GetCheckoutContextMsg().(type) {

	case *pb.CheckoutContext_CartChange:
		if err := cartProcessor(cc.GetCartChange(), offset); err != nil {
			return err
		}

	case *pb.CheckoutContext_Stock:
		if err := stockProcessor(cc.GetStock()); err != nil {
			return err
		}

	case *pb.CheckoutContext_ProductUpdate:
		if err := productsProcessor(cc.GetProductUpdate()); err != nil {
			return err
		}

	case *pb.CheckoutContext_CartOrder:
		ordersProcessor(cc.GetCartOrder(), offset)

	case nil:
		panic(fmt.Sprintf("checkout context message is empty"))

	default:
		panic(fmt.Sprintf("unexpected type %T in oneof", x))
	}
	return nil
}
