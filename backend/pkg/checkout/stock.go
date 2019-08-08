package checkout

import (
	"git.votum-media.net/event-web-store/event-web-store/backend/pkg/pb"
)

func stockProcessor(s *pb.Stock) error {
	mux.Lock()
	defer mux.Unlock()

	stock[s.Uuid] += s.Quantity

	return nil
}
