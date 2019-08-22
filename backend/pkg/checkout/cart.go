package checkout

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
)

type positions []*Position

func (a positions) Len() int           { return len(a) }
func (a positions) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a positions) Less(i, j int) bool { return a[i].Title < a[j].Title }

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

	mut.RLock()
	defer mut.RUnlock()

	pp := positions{}
	for uuid, count := range carts[cartID] {
		if _, ok := products[uuid]; ok {
			pp = append(pp, &Position{
				ProductID:     uuid,
				Price:         products[uuid].Price,
				Title:         products[uuid].Title,
				SmallImageURL: products[uuid].SmallImageURL,

				Quantity:    count,
				MoreInStock: stock[uuid] > count,
				InStock:     stock[uuid] >= count,
			})
		}
	}
	sort.Sort(pp)

	cart := &Cart{
		Positions: pp,
	}

	bytes, err := proto.Marshal(cart)
	if err != nil {
		log.Printf("failed to serialize cart: %v", err)
	}

	_, err = w.Write(bytes)
	if err != nil {
		log.Printf("failed to send result: %v", err)
	}
}

func addToCart(w http.ResponseWriter, r *http.Request, cartID string) error {
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return fmt.Errorf("failed to read from http body: %v", err)
	}

	var cc ChangeProductQuantity
	err = proto.Unmarshal(bytes, &cc)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return fmt.Errorf("failed to decode cart change: %v", err)
	}

	cc.CartID = cartID

	change := &CheckoutContext{
		CheckoutContextMsg: &CheckoutContext_ChangeProductQuantity{
			ChangeProductQuantity: &cc,
		},
	}
	bytes, err = proto.Marshal(change)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return fmt.Errorf("failed to serialize cart change massage: %v", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: Topic,
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

func cartProcessor(cc *ChangeProductQuantity, msgOffset int64) error {

	offset = msgOffset
	defer offsetChanged.Broadcast()

	cartID := cc.CartID

	mut.Lock()
	defer mut.Unlock()

	if _, ok := carts[cartID]; !ok {
		carts[cartID] = make(map[string]int64)
	}

	carts[cartID][cc.ProductID] = cc.Quantity

	if carts[cartID][cc.ProductID] == 0 {
		delete(carts[cartID], cc.ProductID)
	}

	return nil
}
