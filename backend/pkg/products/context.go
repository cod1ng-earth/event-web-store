package products

import (
	"fmt"
	"log"
	"sync"

	"git.votum-media.net/event-web-store/event-web-store/backend/pkg/pb"
	"git.votum-media.net/event-web-store/event-web-store/backend/pkg/simba"
	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	"github.com/golang/protobuf/proto"
)

var (
	offset   int64
	products map[string]*pb.Product
	mut      sync.RWMutex
)

func init() {
	products = map[string]*pb.Product{}

}

func StartContext(brokers *[]string, cfg *cluster.Config) func() {

	consumer, err := cluster.NewConsumer(*brokers, "products-consumer-group", []string{"products"}, cfg)
	if err != nil {
		log.Panicf("failed to setup kafka consumer: %s", err)
	}

	agent := simba.NewConsumer(consumer, processor)
	go agent.Start()

	return agent.Stop
}

func processor(msg *sarama.ConsumerMessage) error {
	p := pb.ProductUpdate{}
	err := proto.Unmarshal(msg.Value, &p)
	if err != nil {
		return fmt.Errorf("failed to unmarshal kafka product massage %d: %v", msg.Offset, err)
	}

	offset = msg.Offset
	//	log.Printf("offset: %d", offset)

	UUID := string(msg.Key)

	mut.Lock()
	defer mut.Unlock()
	if p.New == nil {
		delete(products, UUID)
	} else {
		products[UUID] = p.New
	}

	sortedByUUID = nil
	sortedByPrice = nil
	sortedByName = nil

	return nil
}
