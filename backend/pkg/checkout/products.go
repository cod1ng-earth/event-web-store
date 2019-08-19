package checkout

func productsProcessor(p *Product) error {
	mut.Lock()
	defer mut.Unlock()

	products[p.ProductID] = p

	return nil
}
