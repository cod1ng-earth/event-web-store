package checkout

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"time"

	"git.votum-media.net/event-web-store/event-web-store/backend/pkg/pb"
	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
)

type cart []cartItem

type cartItem struct {
	Product     *pb.Product `json:"product"`
	Quantity    int64       `json:"quantity"`
	MoreInStock bool        `json:"moreInStock"`
	InStock     bool        `json:"inStock"`
}

func (a cart) Len() int           { return len(a) }
func (a cart) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a cart) Less(i, j int) bool { return a[i].Product.Title < a[j].Product.Title }

func CartHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8000")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Set-Cookie, *")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		return
	}

	cookie, err := r.Cookie("cart")
	if err != nil {
		expiration := time.Now().Add(365 * 24 * time.Hour)
		cookie = &http.Cookie{
			Name:    "cart",
			Value:   randomString(32),
			Expires: expiration,
		}
		http.SetCookie(w, cookie)
	}

	cartID := cookie.Value

	if r.Method == "POST" {
		err := addToCart(w, r, cartID)
		if err != nil {
			log.Printf("failed to add a product to the cart: %v", err)
			return
		}
	}

	mux.Lock()
	defer mux.Unlock()

	cart := cart{}
	for uuid, count := range carts[cartID] {
		if _, ok := products[uuid]; ok {
			cart = append(cart, cartItem{
				Product:     products[uuid],
				Quantity:    count,
				MoreInStock: stock[uuid] > count,
				InStock:     stock[uuid] >= count,
			})
		}
	}
	sort.Sort(cart)

	bytes, err := json.Marshal(cart)
	if err != nil {
		log.Printf("failed to serialize cart: %v", err)
	}

	_, err = w.Write(bytes)
	if err != nil {
		log.Printf("failed to send result: %v", err)
	}
}

func addToCart(w http.ResponseWriter, r *http.Request, cartID string) error {
	decoder := json.NewDecoder(r.Body)
	var cc pb.CartChange
	err := decoder.Decode(&cc)
	if err != nil && err != io.EOF {
		w.WriteHeader(http.StatusBadRequest)
		return fmt.Errorf("failed to decode cart change: %v", err)
	}

	cc.CartID = cartID

	change := &pb.CheckoutContext{
		CheckoutContext: &pb.CheckoutContext_CartChange{
			CartChange: &cc,
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

func cartProcessor(cc *pb.CartChange, msgOffset int64) error {

	offset = msgOffset
	defer offsetChanged.Broadcast()

	cartID := cc.CartID

	mux.Lock()
	defer mux.Unlock()

	if _, ok := carts[cartID]; !ok {
		carts[cartID] = make(map[string]int64)
	}

	switch cc.Action {
	case pb.CartChangeAction_add:
		carts[cartID][cc.Uuid] += 1

	case pb.CartChangeAction_remove:
		carts[cartID][cc.Uuid] -= 1
	}

	if carts[cartID][cc.Uuid] == 0 {
		delete(carts[cartID], cc.Uuid)
	}

	return nil
}
