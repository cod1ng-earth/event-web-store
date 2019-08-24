//go:generate sh -c "../../cmd/dev-tools/simba/simba | gofmt -s > context.go"
//go:generate protoc --go_out=. warehouse.proto

package warehouse

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
