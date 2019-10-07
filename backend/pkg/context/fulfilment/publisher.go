package fulfilment

import (
	public "github.com/cod1ng-earth/event-web-store/backend/pkg/context/fulfilment/public"
)

func (p publisher) publishStockCorrected(internalOffset int64, fact *StockCorrected) error {
	i := &public.StockCorrected{
		ProductID:      fact.ProductID,
		QuantityChange: fact.QuantityChange,
	}
	_, _, err := p.logStockCorrected(internalOffset, i)
	return err
}
