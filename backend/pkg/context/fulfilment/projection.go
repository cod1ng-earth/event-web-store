//go:generate protoc --go_out=. topic.proto
//go:generate ../../../cmd/dev-tools/simba/simba --bridge=github.com/cod1ng-earth/event-web-store/backend/pkg/context/checkout
//go:generate gofmt -s -w context.go

package fulfilment

type model struct {
	inventory map[string]int64
}

func newModel() *model {
	return &model{
		inventory: make(map[string]int64),
	}
}

func updateModelStockCorrected(m *model, offset int64, correction *StockCorrected) error {
	m.inventory[correction.ProductID] += correction.QuantityChange
	return nil
}
