package checkout

import (
	"log"

	"github.com/cod1ng-earth/event-web-store/backend/pkg/pim"
	"github.com/cod1ng-earth/event-web-store/backend/pkg/warehouse"
)

func (m *model) getPimOffset() int64 {
	return m.pimOffset
}

func (m *model) getWarehouseOffset() int64 {
	return m.warehouseOffset
}

func translatePimProduct(c *context, m *model, offset int64, msg *pim.Product) error {
	log.Printf("translate PimProduct %v to checkout", offset)
	_, _, err := c.logProduct(&Product{
		ProductID:     msg.Id,
		Price:         msg.Price,
		Name:          msg.Name,
		SmallImageURL: msg.SmallImageURL,
		PimOffset:     offset,
	})
	return err
}

func translateWarehouseStockCorrected(c *context, m *model, offset int64, msg *warehouse.StockCorrected) error {
	log.Printf("translate stock corrected %v to checkkout", offset)
	_, _, err := c.logStockCorrected(&StockCorrected{
		ProductID:       msg.ProductID,
		QuantityChange:  msg.QuantityChange,
		WarehouseOffset: offset,
	})
	return err
}
