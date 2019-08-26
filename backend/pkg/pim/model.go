//go:generate sh -c "../../cmd/dev-tools/simba/simba > context.go"
///go:generate gofmt -s w context.go
//go:generate protoc --go_out=. pim.proto

package pim

type model struct {
	products map[string]*Product
}

func newModel() *model {
	return &model{
		products: make(map[string]*Product),
	}
}

func updateModelProduct(m *model, offset int64, product *Product) error {
	m.products[product.Id] = product
	return nil
}
