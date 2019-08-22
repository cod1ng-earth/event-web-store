package checkout

func updateModelProduct(m *model, offset int64, p *Product) error {
	m.products[p.ProductID] = p
	return nil
}
