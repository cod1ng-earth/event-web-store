//go:generate ../../cmd/dev-tools/simba/simba
//go:generate gofmt -s -w context.go
//go:generate protoc --go_out=. pim.proto

package pim

import "log"

type model struct {
	products map[string]*Product
}

func newModel() *model {
	return &model{
		products: make(map[string]*Product),
	}
}

func updateModelProduct(m *model, offset int64, product *Product) error {
	log.Printf("updateModelProduct %v", offset)
	m.products[product.Id] = product
	return nil
}
