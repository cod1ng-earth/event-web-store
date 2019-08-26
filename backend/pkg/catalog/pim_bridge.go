package catalog

import (
	"github.com/cod1ng-earth/event-web-store/backend/pkg/pim"
)

func (m *model) getPimOffset() int64 {
	return 0
}

func translatePimProduct(c *context, m *model, offset int64, msg *pim.Product) error {
	return nil
}
