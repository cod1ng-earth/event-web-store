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
