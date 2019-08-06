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

	var keys []string
	mux.Lock()
	for k := range products {
		keys = append(keys, k)
	}
	defer mux.Unlock()

	sort.Strings(keys)

	var pp []*pb.Product
	for i, k := range keys {
		if i > 100 {
			break
		}
		pp = append(pp, products[k])
	}
	bytes, err := json.Marshal(pp)
	if err != nil {
		log.Printf("failed to serialize: %v", err)
	}
	_, err = w.Write(bytes)
	if err != nil {
		log.Printf("failed to send result: %v", err)
	}

	//	bytes, err := json.Marshal(keys)
	//	if err != nil {
	//		log.Printf("failed to serialize: %v", err)
	//	}
	//	_, err = w.Write(bytes)
	//	if err != nil {
	//		log.Printf("failed to send result: %v", err)
	//	}

	//	fmt.Fprintf(w, "{")
	//	for i, k := range keys {
	//		fmt.Fprintf(w, "\"%s\"", k)
	//		if i != len(keys)-1 {
	//			fmt.Fprintf(w, ",")
	//		}
	//	}
	//	fmt.Fprintf(w, "}")
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
