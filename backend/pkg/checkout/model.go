//go:generate sh -c "../../cmd/dev-tools/simba/simba > context.go"
///go:generate gofmt -s w context.go
//go:generate protoc --go_out=. checkout.proto

package checkout

type model struct {
	products map[string]*Product
	stock    map[string]int64
	carts    map[string]map[string]int64
	orders   map[string]map[string]int64
}

func newModel() *model {
	return &model{
		products: make(map[string]*Product),
		stock:    make(map[string]int64),
		carts:    make(map[string]map[string]int64),
		orders:   make(map[string]map[string]int64),
	}
}
