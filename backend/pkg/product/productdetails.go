package product

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
	mux      sync.Mutex
)

func StartHandler(brokers *[]string, cfg *cluster.Config) (http.HandlerFunc, func()) {

	consumer, err := cluster.NewConsumer(*brokers, "productdetail-consumer-group", []string{"products"}, cfg)
	if err != nil {
		log.Panicf("failed to setup kafka consumer: %s", err)
	}

	agent := simba.NewConsumer(consumer, processor)
	go agent.Start()

	offset = 0
	products = map[string]*pb.Product{}

	return handler, agent.Stop
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8000")
	//	fmt.Fprintf(w, "offset: %d", offset)

	mux.Lock()
	defer mux.Unlock()

	uuid := r.FormValue("uuid")
	if uuid == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("uuid is missing"))
		if err != nil {
			log.Printf("failed to send result: %v", err)
		}
		return
	}

	p, found := products[uuid]
	if !found {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	bytes, err := json.Marshal(p)
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
	} else {
		products[UUID] = p.New
	}

	return nil
}
