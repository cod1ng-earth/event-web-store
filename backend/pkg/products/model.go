package products

import (
	"log"
	"sort"

	"git.votum-media.net/event-web-store/event-web-store/backend/pkg/pb"
	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
)

type model struct {
	products map[string]*pb.Product

	sortedByUUID  productsByUUID
	sortedByPrice productsByPrice
	sortedByTitle productsByTitle
}

func newModel() *model {
	return &model{
		products: make(map[string]*pb.Product),
	}
}

type productsByUUID []*pb.Product

func (pp productsByUUID) Len() int           { return len(pp) }
func (pp productsByUUID) Swap(i, j int)      { pp[i], pp[j] = pp[j], pp[i] }
func (pp productsByUUID) Less(i, j int) bool { return pp[i].Uuid < pp[j].Uuid }

type productsByPrice []*pb.Product

func (pp productsByPrice) Len() int           { return len(pp) }
func (pp productsByPrice) Swap(i, j int)      { pp[i], pp[j] = pp[j], pp[i] }
func (pp productsByPrice) Less(i, j int) bool { return pp[i].Price < pp[j].Price }

type productsByTitle []*pb.Product

func (pp productsByTitle) Len() int           { return len(pp) }
func (pp productsByTitle) Swap(i, j int)      { pp[i], pp[j] = pp[j], pp[i] }
func (pp productsByTitle) Less(i, j int) bool { return pp[i].Title < pp[j].Title }

func modelRealtimeUpdater(msg *sarama.ConsumerMessage, m *model) {

	p := pb.ProductUpdate{}
	err := proto.Unmarshal(msg.Value, &p)
	if err != nil {
		log.Panicf("failed to unmarshal kafka product massage %d: %v", msg.Offset, err)
	}

	if p.New == nil {
		delete(m.products, p.Old.Uuid)

		idx := sort.Search(len(m.sortedByUUID), func(i int) bool { return m.sortedByUUID[i].Uuid >= p.Old.Uuid })
		m.sortedByUUID = remove(m.sortedByUUID, idx)

		idx = sort.Search(len(m.sortedByPrice), func(i int) bool { return m.sortedByPrice[i].Price >= p.Old.Price })
		for m.sortedByPrice[idx].Uuid != p.Old.Uuid {
			idx++
		}
		m.sortedByPrice = remove(m.sortedByPrice, idx)

		idx = sort.Search(len(m.sortedByTitle), func(i int) bool { return m.sortedByTitle[i].Title >= p.Old.Title })
		for m.sortedByTitle[idx].Uuid != p.Old.Uuid {
			idx++
		}
		m.sortedByTitle = remove(m.sortedByTitle, idx)
	} else if p.Old == nil {
		m.products[p.New.Uuid] = p.New

		idx := sort.Search(len(m.sortedByUUID), func(i int) bool { return m.sortedByUUID[i].Uuid >= p.New.Uuid })
		m.sortedByUUID = insert(m.sortedByUUID, idx, p.New)

		idx = sort.Search(len(m.sortedByPrice), func(i int) bool { return m.sortedByPrice[i].Price >= p.New.Price })
		m.sortedByPrice = insert(m.sortedByPrice, idx, p.New)

		idx = sort.Search(len(m.sortedByTitle), func(i int) bool { return m.sortedByTitle[i].Title >= p.New.Title })
		m.sortedByTitle = insert(m.sortedByTitle, idx, p.New)
	} else {
		m.products[p.New.Uuid] = p.New

		idx := sort.Search(len(m.sortedByUUID), func(i int) bool { return m.sortedByUUID[i].Uuid >= p.New.Uuid })
		m.sortedByUUID[idx] = p.New

		idx = sort.Search(len(m.sortedByPrice), func(i int) bool { return m.sortedByPrice[i].Price >= p.Old.Price })
		for m.sortedByPrice[idx].Uuid != p.Old.Uuid {
			idx++
		}
		m.sortedByPrice = remove(m.sortedByPrice, idx)

		idx = sort.Search(len(m.sortedByPrice), func(i int) bool { return m.sortedByPrice[i].Price >= p.New.Price })
		m.sortedByPrice = insert(m.sortedByPrice, idx, p.New)

		idx = sort.Search(len(m.sortedByTitle), func(i int) bool { return m.sortedByTitle[i].Title >= p.Old.Title })
		for m.sortedByTitle[idx].Uuid != p.Old.Uuid {
			idx++
		}
		m.sortedByTitle = remove(m.sortedByTitle, idx)

		idx = sort.Search(len(m.sortedByTitle), func(i int) bool { return m.sortedByTitle[i].Title >= p.New.Title })
		m.sortedByTitle = insert(m.sortedByTitle, idx, p.New)
	}
}

func remove(slice []*pb.Product, i int) []*pb.Product {
	return append(slice[:i], slice[i+1:]...)
}

func insert(slice []*pb.Product, i int, p *pb.Product) []*pb.Product {
	return append(slice[:i], append([]*pb.Product{p}, slice[i:]...)...)
}

func modelRealtimeFinalizer(m *model) {
}

func modelBatchUpdater(msg *sarama.ConsumerMessage, m *model) {

	p := pb.ProductUpdate{}
	err := proto.Unmarshal(msg.Value, &p)
	if err != nil {
		log.Panicf("failed to unmarshal kafka product massage %d: %v", msg.Offset, err)
	}

	if p.New == nil {
		delete(m.products, p.Old.Uuid)
	} else {
		m.products[p.New.Uuid] = p.New
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
