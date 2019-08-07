package catalog

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
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
	keys     []string
	mux      sync.Mutex
)

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func StartHandler(brokers *[]string, cfg *cluster.Config) (http.HandlerFunc, func()) {

	consumer, err := cluster.NewConsumer(*brokers, "catalog-consumer-group", []string{"products"}, cfg)
	if err != nil {
		log.Panicf("failed to setup kafka consumer: %s", err)
	}

	agent := simba.NewConsumer(consumer, processor)
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

	pageParam := r.URL.Query().Get("page")
	if pageParam == "" {
		pageParam = "0"
	}

	page, err := strconv.Atoi(pageParam)
	if err != nil {
		page = 0
		log.Printf("failed to parse page param: %v\n", err)
	}

	itemsPerPage := 100
	pages := len(keys) / itemsPerPage
	page = Max(Min(page, pages-1), 0)
	startIdx := page * itemsPerPage
	endIdx := Min(startIdx+itemsPerPage, len(keys))

	var pp []*pb.Product = make([]*pb.Product, itemsPerPage)
	for i, k := range keys[startIdx:endIdx] {
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
	} else {
		products[UUID] = &pb.Product{
			Uuid:  p.New.Uuid,
			Title: p.New.Title,
			Price: p.New.Price,
		}
	}
	keys = nil

	return nil
}
