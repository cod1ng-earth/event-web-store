package catalog

import (
	"log"
	"sort"

	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
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

func modelRealtimeUpdater(msg *sarama.ConsumerMessage, m *model) {

	new := &Product{}
	err := proto.Unmarshal(msg.Value, new)
	if err != nil {
		log.Panicf("failed to unmarshal kafka product massage %d: %v", msg.Offset, err)
	}

	old, oldFound := m.products[new.Id]

	if new.Disabled {

		if !oldFound {
			return
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

		return
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
}

func remove(slice []*Product, i int) []*Product {
	return append(slice[:i], slice[i+1:]...)
}

func insert(slice []*Product, i int, p *Product) []*Product {
	return append(slice[:i], append([]*Product{p}, slice[i:]...)...)
}

func modelRealtimeFinalizer(m *model) {
}

func modelBatchUpdater(msg *sarama.ConsumerMessage, m *model) {

	new := &Product{}
	err := proto.Unmarshal(msg.Value, new)
	if err != nil {
		log.Panicf("failed to unmarshal kafka product massage %d: %v", msg.Offset, err)
	}

	if new.Disabled {
		delete(m.products, new.Id)
	} else {
		m.products[new.Id] = new
	}
}

func modelBatchFinalizer(m *model) {

	m.sortedByUUID = nil
	m.sortedByPrice = nil
	m.sortedByTitle = nil

	for _, v := range m.products {
		m.sortedByUUID = append(m.sortedByUUID, v)
	}
	sort.Sort(m.sortedByUUID)

	for _, v := range m.products {
		m.sortedByPrice = append(m.sortedByPrice, v)
	}
	sort.Sort(m.sortedByPrice)

	for _, v := range m.products {
		m.sortedByTitle = append(m.sortedByTitle, v)
	}
	sort.Sort(m.sortedByTitle)
}
