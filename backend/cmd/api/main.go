package main

import (
	"log"
	"net/http"

	"git.votum-media.net/event-web-store/event-web-store/backend/pkg/checkout"
	"git.votum-media.net/event-web-store/event-web-store/backend/pkg/products"
	"github.com/Shopify/sarama"
	"github.com/bsm/sarama-cluster"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	brokerList = kingpin.Flag("brokerList", "List of brokers to connect").Default("localhost:9092").OverrideDefaultFromEnvar("BROKER_LIST").Strings()
)

func main() {
	log.Println("Hello, world!")

	kingpin.Parse()

	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Group.Return.Notifications = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Flush.MaxMessages = 500
	config.Producer.Return.Successes = true

	shutdown := products.StartContext(brokerList, config)
	defer shutdown()
	http.HandleFunc("/product", products.PDPHandler)
	http.HandleFunc("/products", products.CatalogHandler)

	shutdown = checkout.StartContext(brokerList, config)
	defer shutdown()
	http.HandleFunc("/cart", checkout.CartHandler)
	http.HandleFunc("/orderCart", checkout.OrderHandler)

	log.Println("listening on port :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
