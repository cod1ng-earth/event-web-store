package catalog

import (
	"sort"
)

type model struct {
	products map[string]*Product

	sortedByUUID  productsByUUID
	sortedByPrice productsByPrice
	sortedByTitle productsByTitle
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

type productsByTitle []*Product

func (pp productsByTitle) Len() int           { return len(pp) }
func (pp productsByTitle) Swap(i, j int)      { pp[i], pp[j] = pp[j], pp[i] }
func (pp productsByTitle) Less(i, j int) bool { return pp[i].Title < pp[j].Title }

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

		idx = sort.Search(len(m.sortedByTitle), func(i int) bool { return m.sortedByTitle[i].Title >= old.Title })
		for m.sortedByTitle[idx].Id != old.Id {
			idx++
		}
		m.sortedByTitle = remove(m.sortedByTitle, idx)

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

		idx = sort.Search(len(m.sortedByTitle), func(i int) bool { return m.sortedByTitle[i].Title >= old.Title })
		for m.sortedByTitle[idx].Id != old.Id {
			idx++
		}
		m.sortedByTitle = remove(m.sortedByTitle, idx)
	}

	idx = sort.Search(len(m.sortedByPrice), func(i int) bool { return m.sortedByPrice[i].Price >= new.Price })
	m.sortedByPrice = insert(m.sortedByPrice, idx, new)

	idx = sort.Search(len(m.sortedByTitle), func(i int) bool { return m.sortedByTitle[i].Title >= new.Title })
	m.sortedByTitle = insert(m.sortedByTitle, idx, new)

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
	m.sortedByTitle = nil

	for _, v := range m.products {
		m.sortedByUUID = append(m.sortedByUUID, v)
		m.sortedByPrice = append(m.sortedByPrice, v)
		m.sortedByTitle = append(m.sortedByTitle, v)
	}
	sort.Sort(m.sortedByUUID)
	sort.Sort(m.sortedByPrice)
	sort.Sort(m.sortedByTitle)

	return nil
}
