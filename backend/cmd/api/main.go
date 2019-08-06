package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	brokerList = kingpin.Flag("brokerList", "List of brokers to connect").Default("localhost:9092").Strings()
	topic      = kingpin.Flag("topic", "Topic name").Default("products").String()
)

func main() {
	kingpin.Parse()

	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Group.Return.Notifications = true
	topics := []string{*topic}
	consumer, err := cluster.NewConsumer(*brokerList, "inventory-categories-v1", topics, config)
	if err != nil {
		log.Panicf("failed to setup kafka consumer: %s", err)
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Panicf("failed to close kafka consumer: %s", err)
		}
	}()

	log.Println("Hello, world!")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8000")
		fmt.Fprintf(w, "Hello, you have requested: %s %d", r.URL.Path, t.UnixNano())
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
