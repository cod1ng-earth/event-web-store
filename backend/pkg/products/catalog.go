package products

import (
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"git.votum-media.net/event-web-store/event-web-store/backend/pkg/pb"
	"github.com/golang/protobuf/proto"
)

type catalogPage struct {
	Data []*pb.Product       `json:"data"`
	Meta catalogPageMetadata `json:"meta"`
}

type catalogPageMetadata struct {
	TotalItems   int    `json:"total_items"`
	TotalPages   int    `json:"total_pages"`
	CurrentPage  int    `json:"current_page"`
	SetPageTo    int    `json:"set_page_to"`
	ItemsPerPage int    `json:"items_per_page"`
	Sorting      string `json:"sorting"`
	Filtering    string `json:"filtering"`
}

func (c *context) NewCatalogHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8000")

		itemsPerPageInt, err := strconv.Atoi(r.URL.Query().Get("itemsPerPage"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		itemsPerPage := int64(itemsPerPageInt)
		sort := r.URL.Query().Get("sort")
		prefix := r.URL.Query().Get("prefix")
		pageInt, err := strconv.Atoi(r.URL.Query().Get("page"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		page := int64(pageInt)

		model, close := c.read()
		defer close()
		pp := loadProducts(sort, prefix, model)

		totalItems := int64(len(pp))
		totalPages := calculateTotalPages(totalItems, itemsPerPage)
		newPage := calculatePage(page, totalPages)

		startIdx := newPage * itemsPerPage
		endIdx := startIdx + itemsPerPage
		//		log.Printf("startIdx %v, endIdx %v, totalItems %v", startIdx, endIdx, totalItems)
		if endIdx > totalItems {
			endIdx = totalItems
		}
		pp = pp[startIdx:endIdx]

		products := []*pb.Product{}
		for _, p := range pp {
			products = append(products, &pb.Product{
				Uuid:  p.Uuid,
				Title: p.Title,
				Price: p.Price,
			})
		}

		payload := &pb.CatalogPage{
			Products:     products,
			TotalItems:   totalItems,
			TotalPages:   totalPages,
			CurrentPage:  page,
			SetPageTo:    newPage,
			Sorting:      sort,
			Filtering:    prefix,
			ItemsPerPage: itemsPerPage,
		}

		//		respondJson(w, payload)
		w.Header().Set("Content-Type", "application/protobuf")

		bytes, err := proto.Marshal(payload)
		if err != nil {
			log.Printf("failed to serialize: %v", err)
		}

		_, err = w.Write(bytes)
		if err != nil {
			log.Printf("failed to send result: %v", err)
		}
	}
}

func loadProducts(sorting string, prefix string, m *model) []*pb.Product {
	if prefix != "" {
		pp := m.sortedByTitle
		startIdx := sort.Search(len(pp), func(i int) bool { return pp[i].Title >= prefix })
		pp = pp[startIdx:]
		//		endrunes := []rune(prefix)
		//		endrunes[len(endrunes)-1]++
		//		end := string(endrunes)
		endIdx := sort.Search(len(pp), func(i int) bool { return !strings.HasPrefix(pp[i].Title, prefix) })
		pp = pp[:endIdx]
		//	pp := m.sortedByTitle
		//	startIdx := sort.Search(len(pp), func(i int) bool { return pp[i].Title >= prefix })
		//	//pp = pp[startIdx:]
		//	endIdx := sort.Search(len(pp), func(i int) bool { return pp[i].Title > prefix })
		//	//pp = pp[:endIdx]

		//log.Printf("startIdx: %v endIdx: %v, prefix: %v", startIdx, startIdx+endIdx, prefix)
		//		log.Printf("%s", pp[startIdx])
		//		log.Printf("%v", pp[endIdx].Title >= prefix)
		//		log.Printf("%v", pp[endIdx+1].Title >= prefix)
		//		log.Printf("%v", pp[endIdx].Title >= prefix)
		//		log.Printf("%v", pp[endIdx+1].Title >= prefix)
		//

		//log.Printf("%s - %s - %s", m.sortedByTitle[startIdx+endIdx-1].Title, m.sortedByTitle[startIdx+endIdx].Title, m.sortedByTitle[startIdx+endIdx+1].Title)

		switch sorting {
		case "uuid":
			ppp := productsByUUID{}
			for _, v := range pp {
				log.Println(v.Title)
				ppp = append(ppp, v)
			}
			sort.Sort(ppp)
		case "price":
			ppp := productsByPrice{}
			for _, v := range pp {
				log.Println(v.Title)
				ppp = append(ppp, v)
			}
			sort.Sort(ppp)
		case "name":
			return pp
		default:
			log.Printf("sorting %s unknown", sorting)
		}
	}

	switch sorting {
	case "uuid":
		return m.sortedByUUID
	case "price":
		return m.sortedByPrice
	case "name":
		return m.sortedByTitle
	default:
		log.Printf("sorting %s unknown", sorting)
	}

	return []*pb.Product{}
}

func calculatePage(page, totalPages int64) int64 {
	if page > totalPages-1 {
		page = totalPages - 1
	}
	if page < 0 {
		page = 0
	}
	return page
}

func calculateTotalPages(totalItems, itemsPerPage int64) int64 {
	n := totalItems / itemsPerPage
	if (totalItems % itemsPerPage) != 0 {
		n++
	}
	return n
}
