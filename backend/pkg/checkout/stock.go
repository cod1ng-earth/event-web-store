package checkout

import (
	"git.votum-media.net/event-web-store/event-web-store/backend/pkg/pb"
)

func stockProcessor(s *pb.Stock) error {
	mut.Lock()
	defer mut.Unlock()

	stock[s.Uuid] += s.Quantity

	return nil
}
