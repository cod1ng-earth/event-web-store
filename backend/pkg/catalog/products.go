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

type productsByUUID []*pb.Product
type productsByPrice []*pb.Product
type productsByName []*pb.Product

type catalogPayloadMeta struct {
	TotalItems   int `json:"total_items"`
	TotalPages   int `json:"total_pages"`
	CurrentPage  int `json:"current_page"`
	ItemsPerPage int `json:"items_per_page"`
}

type catalogPayload struct {
	Data []*pb.Product      `json:"data"`
	Meta catalogPayloadMeta `json:"meta"`
}

var (
	offset        int64
	products      map[string]*pb.Product
	sortedByUUID  productsByUUID
	sortedByPrice productsByPrice
	sortedByName  productsByName
	mux           sync.Mutex
)

const (
	itemsPerPage = 100
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

func (a productsByUUID) Len() int           { return len(a) }
func (a productsByUUID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a productsByUUID) Less(i, j int) bool { return a[i].Uuid < a[j].Uuid }

func (a productsByPrice) Len() int           { return len(a) }
func (a productsByPrice) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a productsByPrice) Less(i, j int) bool { return a[i].Price < a[j].Price }

func (a productsByName) Len() int           { return len(a) }
func (a productsByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a productsByName) Less(i, j int) bool { return a[i].Title < a[j].Title }

func getProductsByUUID() []*pb.Product {
	if sortedByUUID == nil {
		for _, v := range products {
			sortedByUUID = append(sortedByUUID, v)
		}
		sort.Sort(sortedByUUID)
	}
	return sortedByUUID
}

func getProductsByPrice() []*pb.Product {
	if sortedByPrice == nil {
		for _, v := range products {
			sortedByPrice = append(sortedByPrice, v)
		}
		sort.Sort(sortedByPrice)
	}
	return sortedByPrice
}

func getProductsByName() []*pb.Product {
	if sortedByName == nil {
		for _, v := range products {
			sortedByName = append(sortedByName, v)
		}
		sort.Sort(sortedByName)
	}
	return sortedByName
}

func getPages() int {
	return len(products) / itemsPerPage
}

func clampPage(page int) int {
	return Max(Min(page, getPages()-1), 0)
}

func getProducts(page int, sorting string) ([]*pb.Product, error) {
	page = clampPage(page)
	startIdx := page * itemsPerPage
	endIdx := Min(startIdx+itemsPerPage, len(products))

	switch sorting {
	case "uuid":
		return getProductsByUUID()[startIdx:endIdx], nil
	case "price":
		return getProductsByPrice()[startIdx:endIdx], nil
	case "name":
		return getProductsByName()[startIdx:endIdx], nil
	default:
		return nil, fmt.Errorf("sorting %s unknown", sorting)
	}
}

func StartHandler(brokers *[]string, cfg *cluster.Config) (http.HandlerFunc, func()) {

	consumer, err := cluster.NewConsumer(*brokers, "catalog-consumer-group", []string{"products"}, cfg)
	if err != nil {
		log.Panicf("failed to setup kafka consumer: %s", err)
	}

	agent := simba.NewConsumer(consumer, processor)
	go agent.Start()

	products = map[string]*pb.Product{}

	return Handler, agent.Stop
}

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8000")
	w.Header().Set("Content-Type", "application/json")
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
		log.Printf("failed to parse page param: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	mux.Lock()
	defer mux.Unlock()
	pp, err := getProducts(page, sortParam)
	if err != nil {
		log.Printf("failed to get products: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	payload := catalogPayload{pp, catalogPayloadMeta{len(products), getPages(), clampPage(page), itemsPerPage}}
	bytes, err := json.Marshal(payload)
	if err != nil {
		log.Printf("failed to serialize products: %v", err)
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
		return fmt.Errorf("failed to unmarshal kafka product massage %d: %v", msg.Offset, err)
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
	sortedByUUID = nil
	sortedByPrice = nil
	sortedByName = nil

	return nil
}
