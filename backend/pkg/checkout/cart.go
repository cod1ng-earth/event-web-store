package checkout

import (
	"encoding/json"
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
	offset   int64
	products map[string]*pb.Product
	mux      sync.Mutex
	producer sarama.AsyncProducer
)

func StartCartHandler(brokers *[]string, cfg *cluster.Config) (http.HandlerFunc, func()) {

	consumer, err := cluster.NewConsumer(*brokers, "productdetail-cart-group", []string{"cart", "products"}, cfg)
	if err != nil {
		log.Panicf("failed to setup kafka consumer: %s", err)
	}

	producer, err = sarama.NewAsyncProducer(*brokers, &cfg.Config)
	if err != nil {
		log.Panicf("failed to setup the kafka producer: %s", err)
	}

	agent := simba.NewConsumer(consumer, cartProcessor)
	go agent.Start()

	offset = 0
	products = map[string]*pb.Product{}

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

	log.Printf("cookie: %v", cookie)
	log.Printf("cookie value: %v", cookie.Value)

	producer.Input() <- &sarama.ProducerMessage{
		Topic: "cart",
		Key:   sarama.StringEncoder(cookie.Value),
		Value: sarama.ByteEncoder(bytes),
	}
}

func cartProcessor(msg *sarama.ConsumerMessage) error {
	return nil
}
