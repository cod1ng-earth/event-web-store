package checkout

import (
	"git.votum-media.net/event-web-store/event-web-store/backend/pkg/pb"
)

func productsProcessor(p *pb.ProductUpdate) error {
	mut.Lock()
	defer mut.Unlock()

	if p.New == nil {
		delete(products, p.Old.Uuid)
	} else {
		products[p.New.Uuid] = p.New
	}

	return nil
}
