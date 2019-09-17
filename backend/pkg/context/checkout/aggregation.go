//go:generate protoc --go_out=. api.proto
//go:generate protoc --go_out=. topic.proto
//go:generate ../../../cmd/dev-tools/simba/simba --bridge=github.com/cod1ng-earth/event-web-store/backend/pkg/context/pim --bridge=github.com/cod1ng-earth/event-web-store/backend/pkg/context/fulfilment
//go:generate gofmt -s -w context.go

package checkout

import "log"

type model struct {
	products map[string]*Product
	stock    map[string]int64
	carts    map[string]map[string]int64
	orders   map[string]map[string]int64

	pimOffset        int64
	fulfilmentOffset int64
}

func newModel() *model {
	return &model{
		products: make(map[string]*Product),
		stock:    make(map[string]int64),
		carts:    make(map[string]map[string]int64),
		orders:   make(map[string]map[string]int64),
	}
}

func updateModelStockCorrected(m *model, offset int64, s *StockCorrected) error {
	m.stock[s.ProductID] += s.QuantityChange
	m.fulfilmentOffset = s.FulfilmentOffset
	return nil
}

func updateModelProduct(m *model, offset int64, p *Product) error {
	m.products[p.ProductID] = p
	m.pimOffset = p.PimOffset
	return nil
}

func updateModelOrderCart(m *model, offset int64, p *OrderCart) error {

	cartID := p.CartID

	if _, found := m.orders[cartID]; found {
		return nil
	}

	if _, found := m.carts[cartID]; !found {
		return nil
	}

	cart := m.carts[cartID]

	for uuid, quantity := range cart {
		stock, found := m.stock[uuid]
		if !found {
			return nil
		}
		if quantity > stock {
			return nil
		}
	}

	for uuid, quantity := range cart {
		m.stock[uuid] = m.stock[uuid] - quantity
	}

	m.orders[p.CartID] = m.carts[p.CartID]

	delete(m.carts, p.CartID)

	log.Printf("order %s was created", p.CartID)

	return nil
}

func updateModelChangeProductQuantity(m *model, offset int64, cc *ChangeProductQuantity) error {

	cartID := cc.CartID

	if _, ok := m.carts[cartID]; !ok {
		m.carts[cartID] = make(map[string]int64)
	}

	m.carts[cartID][cc.ProductID] = cc.Quantity

	if m.carts[cartID][cc.ProductID] == 0 {
		delete(m.carts[cartID], cc.ProductID)
	}

	return nil
}
