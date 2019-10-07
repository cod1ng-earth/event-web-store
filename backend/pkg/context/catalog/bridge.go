package catalog

import (
	"log"

	pim "github.com/cod1ng-earth/event-web-store/backend/pkg/context/pim/public"
)

func (m *model) getPimOffset() int64 {
	return m.pimOffset
}

func translatePimProduct(c *context, m *model, offset int64, msg *pim.Product) error {
	log.Printf("translatePimProduct %v", offset)
	_, _, err := c.internalTopic.logProduct(&Product{
		Id:            msg.Id,
		Price:         msg.Price,
		Name:          msg.Name,
		Description:   msg.Description,
		Longtext:      msg.Longtext,
		Category:      msg.Category,
		SmallImageURL: msg.SmallImageURL,
		LargeImageURL: msg.LargeImageURL,
		Disabled:      msg.Disabled,
		PimOffset:     offset,
	})
	return err
}
