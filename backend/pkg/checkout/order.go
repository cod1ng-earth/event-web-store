package checkout

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"git.votum-media.net/event-web-store/event-web-store/backend/pkg/pb"
	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
)

func OrderHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8000")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Set-Cookie, *")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		return
	}

	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cookie, err := r.Cookie("cart")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("failed to find cart cookie id: %v", err)
		return
	}

	cartID := cookie.Value

	err = orderCart(w, r, cartID)
	if err != nil {
		log.Printf("failed to create order request: %v", err)
		return
	}

	mut.RLock()
	defer mut.RUnlock()

	_, isOrdered := orders[cartID]

	status := map[string]string{}
	if isOrdered {
		status["status"] = "success"

		expiration := time.Now().Add(365 * 24 * time.Hour)
		cookie = &http.Cookie{
			Name:    "cart",
			Value:   randomString(32),
			Expires: expiration,
		}
		http.SetCookie(w, cookie)

	} else {
		status["status"] = "failure"
	}

	bytes, err := json.Marshal(status)
	if err != nil {
		log.Printf("failed to serialize order status: %v", err)
	}

	_, err = w.Write(bytes)
	if err != nil {
		log.Printf("failed to send result: %v", err)
	}
}

func orderCart(w http.ResponseWriter, r *http.Request, cartID string) error {

	change := &pb.CheckoutContext{
		CheckoutContextMsg: &pb.CheckoutContext_CartOrder{
			CartOrder: &pb.CartOrder{
				CartID: cartID,
			},
		},
	}
	bytes, err := proto.Marshal(change)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return fmt.Errorf("failed to serialize cart change massage: %v", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: "checkout",
		Key:   sarama.StringEncoder(cartID),
		Value: sarama.ByteEncoder(bytes),
	}
	_, msgOffset, err := producer.SendMessage(msg)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return fmt.Errorf("failed to send cart change to kafka: %v", err)
	}

	offsetChanged.L.Lock()
	for offset < msgOffset {
		offsetChanged.Wait()
	}
	offsetChanged.L.Unlock()

	return nil
}

func ordersProcessor(p *pb.CartOrder, msgOffset int64) {
	mut.Lock()
	defer mut.Unlock()

	offset = msgOffset
	defer offsetChanged.Broadcast()

	cartID := p.CartID

	if _, found := orders[cartID]; found {
		return
	}

	if _, found := carts[cartID]; !found {
		return
	}

	cart := carts[cartID]

	for uuid, quantity := range cart {
		stock, found := stock[uuid]
		if !found {
			return
		}
		if quantity > stock {
			return
		}
	}

	for uuid, quantity := range cart {
		stock[uuid] = stock[uuid] - quantity
	}

	orders[p.CartID] = carts[p.CartID]

	delete(carts, p.CartID)

	log.Printf("order %s was created", p.CartID)
}
