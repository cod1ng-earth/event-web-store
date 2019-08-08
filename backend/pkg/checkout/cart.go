package checkout

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"git.votum-media.net/event-web-store/event-web-store/backend/pkg/pb"
	"git.votum-media.net/event-web-store/event-web-store/backend/pkg/simba"
	"github.com/Shopify/sarama"
	"github.com/bsm/sarama-cluster"
	"github.com/golang/protobuf/proto"
)

var (
	cartOffset    int64
	offsetChanged *sync.Cond
	carts         map[string]map[string]int
	mux           sync.Mutex
	producer      sarama.SyncProducer
)

func StartCartHandler(brokers *[]string, cfg *cluster.Config) (http.HandlerFunc, func()) {

	offsetChanged = sync.NewCond(&sync.Mutex{})
	carts = make(map[string]map[string]int)
	products = map[string]*pb.Product{}

	cartConsumer, err := cluster.NewConsumer(*brokers, "checkout-cart-group", []string{"cart"}, cfg)
	if err != nil {
		log.Panicf("failed to setup kafka consumer: %s", err)
	}

	productsConsumer, err := cluster.NewConsumer(*brokers, "checkout-products-group", []string{"products"}, cfg)
	if err != nil {
		log.Panicf("failed to setup kafka consumer: %s", err)
	}

	producer, err = sarama.NewSyncProducer(*brokers, &cfg.Config)
	if err != nil {
		log.Panicf("failed to setup the kafka producer: %s", err)
	}

	cartAgent := simba.NewConsumer(cartConsumer, cartProcessor)
	go cartAgent.Start()

	productsAgent := simba.NewConsumer(productsConsumer, productsProcessor)
	go productsAgent.Start()

	return cartHandler, func() {
		cartAgent.Stop()
		productsAgent.Stop()
		if err := producer.Close(); err != nil {
			log.Printf("failed to close the kafka producer: %s", err)
		}
	}
}

type cartItem struct {
	Product  *pb.Product `json:"product,omitempty"`
	Quantity int         `json:"quantity,omitempty"`
}

func cartHandler(w http.ResponseWriter, r *http.Request) {
	//	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8000|http://localhost:8080")
	//	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Set-Cookie")

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

	cart := []cartItem{}
	for uuid, count := range carts[cartID] {

		if _, ok := products[uuid]; ok {
			cart = append(cart, cartItem{
				Product:  products[uuid],
				Quantity: count,
			})
		}

	}

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

	bytes, err := proto.Marshal(&cc)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return fmt.Errorf("failed to serialize cart change massage: %v", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: "cart",
		Key:   sarama.StringEncoder(cartID),
		Value: sarama.ByteEncoder(bytes),
	}
	_, msgOffset, err := producer.SendMessage(msg)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return fmt.Errorf("failed to send cart change to kafka: %v", err)
	}

	offsetChanged.L.Lock()
	for cartOffset < msgOffset {
		offsetChanged.Wait()
	}
	offsetChanged.L.Unlock()

	return nil
}

func cartProcessor(msg *sarama.ConsumerMessage) error {
	cc := pb.CartChange{}
	err := proto.Unmarshal(msg.Value, &cc)
	if err != nil {
		return fmt.Errorf("failed to unmarshal kafka cart massage %d: %v", msg.Offset, err)
	}

	cartOffset = msg.Offset
	defer offsetChanged.Broadcast()

	cartID := string(msg.Key)

	mux.Lock()
	defer mux.Unlock()

	if _, ok := carts[cartID]; !ok {
		carts[cartID] = make(map[string]int)
	}

	switch cc.Action {
	case pb.CartChangeAction_add:
		carts[cartID][cc.Uuid] += 1
	case pb.CartChangeAction_remove:
		carts[cartID][cc.Uuid] -= 1
	}

	return nil
}
