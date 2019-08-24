//go:generate sh -c "../../cmd/dev-tools/simba/simba -readLock=wait-free -batch=true | gofmt -s > context.go"
//go:generate protoc --go_out=. catalog.proto

package catalog

import (
	"sort"
)

type model struct {
	products map[string]*Product

	sortedByUUID  productsByUUID
	sortedByPrice productsByPrice
	sortedByName  productsByName
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

func updateModelProduct(m *model, offset int64, new *Product) error {
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

func batchUpdateModelProduct(m *model, offset int64, new *Product) error {
	if new.Disabled {
		delete(m.products, new.Id)
	} else {
		m.products[new.Id] = new
	}
	return nil
}

func batchFinalizeModel(m *model) error {
	m.sortedByUUID = nil
	m.sortedByPrice = nil
	m.sortedByName = nil

	for _, v := range m.products {
		m.sortedByUUID = append(m.sortedByUUID, v)
		m.sortedByPrice = append(m.sortedByPrice, v)
		m.sortedByName = append(m.sortedByName, v)
	}
	sort.Sort(m.sortedByUUID)
	sort.Sort(m.sortedByPrice)
	sort.Sort(m.sortedByName)

	return nil
}
