package fulfilment

import (
	"log"

	checkout "github.com/cod1ng-earth/event-web-store/backend/pkg/context/checkout/public"
)

func translateCheckoutOrderCreated(c *context, m *model, offset int64, msg *checkout.OrderCreated) error {
	log.Printf("translate OrderCreated %v to fulfilment", offset)
	// TODO process order
	return nil
}
