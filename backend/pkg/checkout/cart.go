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
	offset        int64
	offsetChanged *sync.Cond
	carts         map[string]map[string]int
	mux           sync.Mutex
	producer      sarama.SyncProducer
)

func StartCartHandler(brokers *[]string, cfg *cluster.Config) (http.HandlerFunc, func()) {

	offsetChanged = sync.NewCond(&sync.Mutex{})
	offset = 0
	carts = make(map[string]map[string]int)

	consumer, err := cluster.NewConsumer(*brokers, "productdetail-cart-group", []string{"cart"}, cfg)
	if err != nil {
		log.Panicf("failed to setup kafka consumer: %s", err)
	}

	producer, err = sarama.NewSyncProducer(*brokers, &cfg.Config)
	if err != nil {
		log.Panicf("failed to setup the kafka producer: %s", err)
	}

	agent := simba.NewConsumer(consumer, cartProcessor)
	go agent.Start()

	return cartHandler, func() {
		agent.Stop()
		if err := producer.Close(); err != nil {
			log.Panicf("failed to close the kafka producer: %s", err)
		}
	}
}

type cartChange struct {
	Action string
	Uuid   string
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

	decoder := json.NewDecoder(r.Body)
	var cc pb.CartChange
	err := decoder.Decode(&cc)
	if err != nil && err != io.EOF {
		log.Printf("failed to decode cart change: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	bytes, err := proto.Marshal(&cc)
	if err != nil {
		log.Printf("failed to serialize cart change massage: %v", err)
		w.WriteHeader(http.StatusBadRequest)
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

	cartId := cookie.Value

	msg := &sarama.ProducerMessage{
		Topic: "cart",
		Key:   sarama.StringEncoder(cartId),
		Value: sarama.ByteEncoder(bytes),
	}
	_, msgOffset, err := producer.SendMessage(msg)
	if err != nil {
		log.Printf("failed to send cart change to kafka: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	offsetChanged.L.Lock()
	for offset < msgOffset {
		offsetChanged.Wait()
	}
	offsetChanged.L.Unlock()

	mux.Lock()
	defer mux.Unlock()

	cart := []cartItem{}
	for uuid, count := range carts[cartId] {
		cart = append(cart, cartItem{
			Product: &pb.Product{
				Uuid:          uuid,
				Title:         "mock",
				Description:   "mock",
				Longtext:      "mock",
				Category:      "mock",
				SmallImageURL: "mock",
				LargeImageURL: "mock",
				Price:         0.0,
			},
			Quantity: count,
		})
	}

	bytes, err = json.Marshal(cart)
	if err != nil {
		log.Printf("failed to serialize cart: %v", err)
	}

	_, err = w.Write(bytes)
	if err != nil {
		log.Printf("failed to send result: %v", err)
	}
}

func cartProcessor(msg *sarama.ConsumerMessage) error {
	cc := pb.CartChange{}
	err := proto.Unmarshal(msg.Value, &cc)
	if err != nil {
		return fmt.Errorf("failed to unmarshal kafka cart massage %d: %v", msg.Offset, err)
	}

	offset = msg.Offset
	defer offsetChanged.Broadcast()

	cartId := string(msg.Key)

	mux.Lock()
	defer mux.Unlock()

	if _, ok := carts[cartId]; !ok {
		carts[cartId] = make(map[string]int)
	}

	switch cc.Action {
	case pb.CartChangeAction_add:
		carts[cartId][cc.Uuid] += 1
	case pb.CartChangeAction_remove:
		carts[cartId][cc.Uuid] -= 1
	}

	return nil
}
