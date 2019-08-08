package checkout

import (
	"fmt"

	"git.votum-media.net/event-web-store/event-web-store/backend/pkg/pb"
	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
)

var (
	productsOffset int64
	products       map[string]*pb.Product
)

func productsProcessor(msg *sarama.ConsumerMessage) error {
	p := pb.ProductUpdate{}
	err := proto.Unmarshal(msg.Value, &p)
	if err != nil {
		return fmt.Errorf("failed to unmarshal kafka product massage %d: %v", msg.Offset, err)
	}

	productsOffset = msg.Offset
	//	log.Printf("offset: %d", offset)

	UUID := string(msg.Key)

	mux.Lock()
	defer mux.Unlock()
	if p.New == nil {
		delete(products, UUID)
	} else {
		products[UUID] = p.New
	}

	return nil
}
