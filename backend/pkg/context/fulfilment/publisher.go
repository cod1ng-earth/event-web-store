package fulfilment

import (
	public "github.com/cod1ng-earth/event-web-store/backend/pkg/context/fulfilment/public"
)

func publishStockCorrected(c *context, offset int64, fact *StockCorrected) error {
	p := &public.StockCorrected{
		ProductID:      fact.ProductID,
		QuantityChange: fact.QuantityChange,
	}
	_, _, err := c.logPublicStockCorrected(p)
	return err
}
