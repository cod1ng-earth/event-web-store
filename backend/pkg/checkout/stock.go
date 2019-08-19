package checkout

func stockProcessor(s *Stock) error {
	mut.Lock()
	defer mut.Unlock()

	stock[s.ProductID] += s.Quantity

	return nil
}
