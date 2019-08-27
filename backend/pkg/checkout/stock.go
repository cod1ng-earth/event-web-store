package checkout

func updateModelStockCorrected(m *model, offset int64, s *StockCorrected) error {
	m.stock[s.ProductID] += s.QuantityChange
	m.warehouseOffset = s.WarehouseOffset
	return nil
}
