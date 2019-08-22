package checkout

func updateModelStock(m *model, offset int64, s *Stock) error {
	m.stock[s.ProductID] += s.Quantity
	return nil
}
