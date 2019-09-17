package checkout

import (
	"crypto/sha256"
	"encoding/binary"
	fmt "fmt"

	public "github.com/cod1ng-earth/event-web-store/backend/pkg/context/checkout/public"
)

func (p *publisher) publishChangeProductQuantity(offset int64, fact *ChangeProductQuantity) error {
	return nil
}

func (p *publisher) publishStockCorrected(offset int64, fact *StockCorrected) error {
	return nil
}

func (p *publisher) publishProduct(offset int64, fact *Product) error {
	return nil
}

func (p *publisher) publishOrderCart(offset int64, fact *OrderCart) error {
	bb := make([]byte, 8)
	binary.LittleEndian.PutUint64(bb, uint64(offset)) // TODO add a secret component || create random string and store it in kafka
	orderID := sha256.Sum256(bb)
	_, _, err := p.logOrderCreated(&public.OrderCreated{OrderID: fmt.Sprintf("%x", orderID)})
	return err
}
