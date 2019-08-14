package products

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"git.votum-media.net/event-web-store/event-web-store/backend/pkg/pb"
)

type catalogPayload struct {
	Data []*pb.Product      `json:"data"`
	Meta catalogPayloadMeta `json:"meta"`
}

type catalogPayloadMeta struct {
	TotalItems   int    `json:"total_items"`
	TotalPages   int    `json:"total_pages"`
	CurrentPage  int    `json:"current_page"`
	SetPageTo    int    `json:"set_page_to"`
	ItemsPerPage int    `json:"items_per_page"`
	Sorting      string `json:"sorting"`
	Filtering    string `json:"filtering"`
}

func (c *context) NewCatalogHandler() http.HandlerFunc {
	return newCatalogHandler(c.locklessModel)
}

func newCatalogHandler(modelAccess locklessModel) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8000")
		w.Header().Set("Content-Type", "application/json")
		//	fmt.Fprintf(w, "offset: %d", offset)
		pageParam := r.URL.Query().Get("page")
		if pageParam == "" {
			pageParam = "0"
		}

		itemsPerPage, err := strconv.Atoi(r.URL.Query().Get("itemsPerPage"))
		if err != nil {
			itemsPerPage = 100
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

		model, close := modelAccess.read()
		defer close()
		pp := loadProducts(sortParam, prefixParam, model)

		totalItems := len(pp)
		totalPages := calculateTotalPages(totalItems, itemsPerPage)
		newPage := calculatePage(page, totalPages)

		startIdx := newPage * itemsPerPage
		endIdx := min(startIdx+itemsPerPage, totalItems)
		pp = pp[startIdx:endIdx]

		ppCompact := []*pb.Product{}
		for _, p := range pp {
			ppCompact = append(ppCompact, &pb.Product{
				Uuid:  p.Uuid,
				Title: p.Title,
				Price: p.Price,
			})
		}

		payload := catalogPayload{
			Data: ppCompact,
			Meta: catalogPayloadMeta{
				TotalItems:   totalItems,
				TotalPages:   totalPages,
				CurrentPage:  page,
				SetPageTo:    newPage,
				Sorting:      sortParam,
				Filtering:    prefixParam,
				ItemsPerPage: itemsPerPage,
			},
		}
		bytes, err := json.Marshal(payload)
		if err != nil {
			log.Printf("failed to serialize products: %v", err)
		}

		_, err = w.Write(bytes)
		if err != nil {
			log.Printf("failed to send result: %v", err)
		}
	}
}

func loadProducts(sorting string, prefix string, m *model) []*pb.Product {
	var pp []*pb.Product
	switch sorting {
	case "uuid":
		pp = m.sortedByUUID
	case "price":
		pp = m.sortedByPrice
	case "name":
		pp = m.sortedByTitle
	default:
		log.Printf("sorting %s unknown", sorting)
		return pp
	}

	if prefix != "" {
		pp = filterProducts(pp, func(p *pb.Product) bool {
			return strings.HasPrefix(p.Title, prefix)
		})
	}

	return pp
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

func calculatePage(page int, totalPages int) int {
	return max(min(page, totalPages-1), 0)
}

func calculateTotalPages(totalItems, itemsPerPage int) int {
	n := totalItems / itemsPerPage
	if (totalItems % itemsPerPage) != 0 {
		n++
	}
	return n
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
