package products

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"git.votum-media.net/event-web-store/event-web-store/backend/pkg/pb"
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
	sortedByUUID  productsByUUID
	sortedByPrice productsByPrice
	sortedByName  productsByName
)

const (
	itemsPerPage = 100
)

func CatalogHandler(w http.ResponseWriter, r *http.Request) {
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

	prefixParam := r.URL.Query().Get("prefix")

	page, err := strconv.Atoi(pageParam)
	if err != nil {
		log.Printf("failed to parse page param: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	mux.Lock()
	defer mux.Unlock()
	pp, totalItems, err := getProducts(page, sortParam, prefixParam)
	if err != nil {
		log.Printf("failed to get products: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	totalPages := getPageNumber(totalItems)
	currentPage := clampPage(page, totalPages)

	ppMini := []*pb.Product{}
	for _, p := range pp {
		ppMini = append(ppMini, &pb.Product{
			Uuid:  p.Uuid,
			Title: p.Title,
			Price: p.Price,
		})
	}

	payload := catalogPayload{ppMini, catalogPayloadMeta{totalItems, totalPages, currentPage, itemsPerPage}}
	bytes, err := json.Marshal(payload)
	if err != nil {
		log.Printf("failed to serialize products: %v", err)
	}

	_, err = w.Write(bytes)
	if err != nil {
		log.Printf("failed to send result: %v", err)
	}
}

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

func getPageNumber(nItems int) int {
	n := nItems / itemsPerPage
	if (nItems % itemsPerPage) != 0 {
		n++
	}
	return n
}

func clampPage(page int, nPages int) int {
	return Max(Min(page, nPages-1), 0)
}

func filterProducts(vs []*pb.Product, f func(*pb.Product) bool) []*pb.Product {
	vsf := make([]*pb.Product, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func getProducts(page int, sorting string, prefix string) ([]*pb.Product, int, error) {
	var getter func() []*pb.Product
	switch sorting {
	case "uuid":
		getter = getProductsByUUID
	case "price":
		getter = getProductsByPrice
	case "name":
		getter = getProductsByName
	default:
		return nil, 0, fmt.Errorf("sorting %s unknown", sorting)
	}
	pp := getter()
	if prefix != "" {
		pp = filterProducts(pp, func(p *pb.Product) bool {
			return strings.HasPrefix(p.Title, prefix)
		})
	}

	page = clampPage(page, getPageNumber(len(pp)))
	startIdx := page * itemsPerPage
	endIdx := Min(startIdx+itemsPerPage, len(pp))

	return pp[startIdx:endIdx], len(pp), nil
}
