package catalog

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"sync"

	"git.votum-media.net/event-web-store/event-web-store/backend/pkg/pb"
	"git.votum-media.net/event-web-store/event-web-store/backend/pkg/simba"
	"github.com/Shopify/sarama"
	"github.com/bsm/sarama-cluster"
	"github.com/golang/protobuf/proto"
)

var (
	offset   int64
	products map[string]*pb.Product
	keys     []string
	mux      sync.Mutex
)

func StartHandler(c *cluster.Consumer) (http.HandlerFunc, func()) {

	agent := simba.NewConsumer(c, processor)
	go agent.Start()

	offset = 0
	products = map[string]*pb.Product{}

	return Handler, agent.Stop
}

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8000")
	//	fmt.Fprintf(w, "offset: %d", offset)

	mux.Lock()
	defer mux.Unlock()
	if keys == nil {
		for k := range products {
			keys = append(keys, k)
		}
		sort.Strings(keys)
	}

	max := 100
	if len(keys) < max {
		max = len(keys)
	}
	var pp []*pb.Product = make([]*pb.Product, max-1)
	for i, k := range keys[0 : max-1] {
		pp[i] = products[k]
	}

	bytes, err := json.Marshal(pp)
	if err != nil {
		log.Printf("failed to serialize: %v", err)
	}

	_, err = w.Write(bytes)
	if err != nil {
		log.Printf("failed to send result: %v", err)
	}
}

func processor(msg *sarama.ConsumerMessage) error {
	p := pb.ProductUpdate{}
	err := proto.Unmarshal(msg.Value, &p)
	if err != nil {
		return fmt.Errorf("failed to unmarshal kafka massage %d: %v", msg.Offset, err)
	}

	offset = msg.Offset
	//	log.Printf("offset: %d", offset)

	UUID := string(msg.Key)

	mux.Lock()
	defer mux.Unlock()
	if p.New == nil {
		delete(products, UUID)
		keys = nil
	} else {
		products[UUID] = &pb.Product{
			Uuid:  p.New.Uuid,
			Title: p.New.Title,
			Price: p.New.Price,
		}
		keys = nil
	}

	return nil
}
