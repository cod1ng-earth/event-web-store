package checkout

import public "github.com/cod1ng-earth/event-web-store/backend/pkg/context/checkout/public"

func publishChangeProductQuantity(c *context, offset int64, fact *ChangeProductQuantity) error {
	return nil
}

func publishStockCorrected(c *context, offset int64, fact *StockCorrected) error {
	return nil
}

func publishProduct(c *context, offset int64, fact *Product) error {
	return nil
}

func publishOrderCart(c *context, offset int64, fact *OrderCart) error {
	p := &public.OrderCreated{
		OrderID: fact.OrderID, // TODO
	}
	_, _, err := c.logPublicOrderCreated(p)
	return err
}
