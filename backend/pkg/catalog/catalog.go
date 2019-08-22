package catalog

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

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
		if sort == "" {
			sort = "uuid"
		}
		prefix := r.URL.Query().Get("prefix")
		pageInt, err := strconv.Atoi(r.URL.Query().Get("page"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		page := int64(pageInt)

		model, close := c.read()
		defer close()
		pp, err := loadProducts(sort, prefix, model)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

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

		products := []*Product{}
		for _, p := range pp {
			products = append(products, &Product{
				Id:    p.Id,
				Name:  p.Name,
				Price: p.Price,
			})
		}

		payload := &CatalogPage{
			Products:     products,
			TotalItems:   totalItems,
			TotalPages:   totalPages,
			CurrentPage:  page,
			SetPageTo:    newPage,
			Sorting:      sort,
			Filtering:    prefix,
			ItemsPerPage: itemsPerPage,
		}

		respond(w, payload)
	}
}

func loadProducts(sorting string, prefix string, m *model) ([]*Product, error) {
	if prefix != "" {
		pp := m.sortedByName
		startIdx := sort.Search(len(pp), func(i int) bool { return pp[i].Name >= prefix })
		pp = pp[startIdx:]
		endIdx := sort.Search(len(pp), func(i int) bool { return !strings.HasPrefix(pp[i].Name, prefix) })
		pp = pp[:endIdx]

		switch sorting {
		case "uuid":
			ppp := make(productsByUUID, len(pp))
			copy(ppp, pp)
			sort.Sort(ppp)
			return ppp, nil
		case "price":
			ppp := make(productsByPrice, len(pp))
			copy(ppp, pp)
			sort.Sort(ppp)
			return ppp, nil
		case "name":
			return pp, nil
		}
	}

	switch sorting {
	case "uuid":
		return m.sortedByUUID, nil
	case "price":
		return m.sortedByPrice, nil
	case "name":
		return m.sortedByName, nil
	}

	return []*Product{}, fmt.Errorf("sorting %s unknown", sorting)
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
