package main

import (
	"log"
	"net/http"

	"git.votum-media.net/event-web-store/event-web-store/backend/pkg/catalog"
	"git.votum-media.net/event-web-store/event-web-store/backend/pkg/checkout"
	"git.votum-media.net/event-web-store/event-web-store/backend/pkg/product"
	"github.com/Shopify/sarama"
	"github.com/bsm/sarama-cluster"
	"gopkg.in/alecthomas/kingpin.v2"
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
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Flush.MaxMessages = 500

	log.Println("Hello, world!")

	catalogHandler, catalogShutdown := catalog.StartHandler(brokerList, config)
	defer catalogShutdown()
	http.HandleFunc("/products", catalogHandler)

	productHandler, productShutdown := product.StartHandler(brokerList, config)
	defer productShutdown()
	http.HandleFunc("/product", productHandler)

	cartHandler, cartShutdown := checkout.StartCartHandler(brokerList, config)
	defer cartShutdown()
	http.HandleFunc("/cart", cartHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
