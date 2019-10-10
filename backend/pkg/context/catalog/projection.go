//go:generate protoc --go_out=. api.proto
//go:generate protoc --go_out=. topic.proto
//go:generate ../../../cmd/dev-tools/simba/simba --lock=wait-free --bridge=github.com/cod1ng-earth/event-web-store/backend/pkg/context/pim
//go:generate gofmt -s -w context.go

package catalog

import (
	"log"
	"sort"
)

type model struct {
	products map[string]*Product

	sortedByUUID  productsByUUID
	sortedByPrice productsByPrice
	sortedByName  productsByName

	pimOffset int64
}

func newModel() *model {
	return &model{
		products: make(map[string]*Product),
	}
}

type productsByUUID []*Product

func (pp productsByUUID) Len() int           { return len(pp) }
func (pp productsByUUID) Swap(i, j int)      { pp[i], pp[j] = pp[j], pp[i] }
func (pp productsByUUID) Less(i, j int) bool { return pp[i].Id < pp[j].Id }

type productsByPrice []*Product

func (pp productsByPrice) Len() int           { return len(pp) }
func (pp productsByPrice) Swap(i, j int)      { pp[i], pp[j] = pp[j], pp[i] }
func (pp productsByPrice) Less(i, j int) bool { return pp[i].Price < pp[j].Price }

type productsByName []*Product

func (pp productsByName) Len() int           { return len(pp) }
func (pp productsByName) Swap(i, j int)      { pp[i], pp[j] = pp[j], pp[i] }
func (pp productsByName) Less(i, j int) bool { return pp[i].Name < pp[j].Name }

//func (p *PimProduct) toProduct() *Product {
//	return &Product{
//		Id:            p.Id,
//		Price:         p.Price,
//		Name:          p.Name,
//		Description:   p.Description,
//		Longtext:      p.Longtext,
//		Category:      p.Category,
//		SmallImageURL: p.SmallImageURL,
//		LargeImageURL: p.LargeImageURL,
//		Disabled:      p.Disabled,
//	}
//}

func updateModelProduct(m *model, offset int64, new *Product) error {
	log.Printf("updateModelPimProduct %v", offset)

	m.pimOffset = offset

	old, oldFound := m.products[new.Id]

	if new.Disabled {

		if !oldFound {
			return nil
		}

		delete(m.products, new.Id)

		idx := sort.Search(len(m.sortedByUUID), func(i int) bool { return m.sortedByUUID[i].Id >= old.Id })
		m.sortedByUUID = remove(m.sortedByUUID, idx)

		idx = sort.Search(len(m.sortedByPrice), func(i int) bool { return m.sortedByPrice[i].Price >= old.Price })
		for m.sortedByPrice[idx].Id != old.Id {
			idx++
		}
		m.sortedByPrice = remove(m.sortedByPrice, idx)

		idx = sort.Search(len(m.sortedByName), func(i int) bool { return m.sortedByName[i].Name >= old.Name })
		for m.sortedByName[idx].Id != old.Id {
			idx++
		}
		m.sortedByName = remove(m.sortedByName, idx)

		return nil
	}

	m.products[new.Id] = new

	idx := sort.Search(len(m.sortedByUUID), func(i int) bool { return m.sortedByUUID[i].Id >= new.Id })
	if oldFound {
		m.sortedByUUID[idx] = new
	} else {
		m.sortedByUUID = insert(m.sortedByUUID, idx, new)
	}

	if oldFound {
		idx = sort.Search(len(m.sortedByPrice), func(i int) bool { return m.sortedByPrice[i].Price >= old.Price })
		for m.sortedByPrice[idx].Id != old.Id {
			idx++
		}
		m.sortedByPrice = remove(m.sortedByPrice, idx)

		idx = sort.Search(len(m.sortedByName), func(i int) bool { return m.sortedByName[i].Name >= old.Name })
		for m.sortedByName[idx].Id != old.Id {
			idx++
		}
		m.sortedByName = remove(m.sortedByName, idx)
	}

	idx = sort.Search(len(m.sortedByPrice), func(i int) bool { return m.sortedByPrice[i].Price >= new.Price })
	m.sortedByPrice = insert(m.sortedByPrice, idx, new)

	idx = sort.Search(len(m.sortedByName), func(i int) bool { return m.sortedByName[i].Name >= new.Name })
	m.sortedByName = insert(m.sortedByName, idx, new)

	return nil
}

func remove(slice []*Product, i int) []*Product {
	return append(slice[:i], slice[i+1:]...)
}

func insert(slice []*Product, i int, p *Product) []*Product {
	return append(slice[:i], append([]*Product{p}, slice[i:]...)...)
}
