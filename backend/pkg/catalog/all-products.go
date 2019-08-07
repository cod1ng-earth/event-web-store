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

type ProductsByUUID []*pb.Product
type ProductsByPrice []*pb.Product

var (
	offset          int64
	products        map[string]*pb.Product
	productsByUUID  ProductsByUUID
	productsByPrice ProductsByPrice
	mux             sync.Mutex
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

func (a ProductsByUUID) Len() int           { return len(a) }
func (a ProductsByUUID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ProductsByUUID) Less(i, j int) bool { return a[i].Uuid < a[j].Uuid }

func (a ProductsByPrice) Len() int           { return len(a) }
func (a ProductsByPrice) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ProductsByPrice) Less(i, j int) bool { return a[i].Price < a[j].Price }

func GetProductsByUUID() []*pb.Product {
	if productsByUUID == nil {
		for _, v := range products {
			productsByUUID = append(productsByUUID, v)
		}
		sort.Sort(productsByUUID)
	}
	return productsByUUID
}

func GetProductsByPrice() []*pb.Product {
	if productsByPrice == nil {
		for _, v := range products {
			productsByPrice = append(productsByPrice, v)
		}
		sort.Sort(productsByPrice)
	}
	return productsByPrice
}

func GetProducts(page int, sorting string) []*pb.Product {
	itemsPerPage := 100
	pages := len(products) / itemsPerPage
	page = Max(Min(page, pages-1), 0)
	startIdx := page * itemsPerPage
	endIdx := Min(startIdx+itemsPerPage, len(products))

	switch sorting {
	case "uuid":
		return GetProductsByUUID()[startIdx:endIdx]
	case "price":
		return GetProductsByPrice()[startIdx:endIdx]
	}

	return GetProductsByUUID()[startIdx:endIdx]
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
	pageParam := r.URL.Query().Get("page")
	if pageParam == "" {
		pageParam = "0"
	}

	sortParam := r.URL.Query().Get("sort")
	if sortParam == "" {
		sortParam = "uuid"
	}

	page, err := strconv.Atoi(pageParam)
	if err != nil {
		page = 0
		log.Printf("failed to parse page param: %v\n", err)
	}

	mux.Lock()
	defer mux.Unlock()
	pp := GetProducts(page, sortParam)
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
	productsByUUID = nil
	productsByPrice = nil

	return nil
}
