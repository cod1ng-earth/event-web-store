package main

import (
	"log"
	"net/http"

	"git.votum-media.net/event-web-store/event-web-store/backend/pkg/checkout"
	"git.votum-media.net/event-web-store/event-web-store/backend/pkg/products"
	"github.com/Shopify/sarama"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	brokerList = kingpin.Flag("brokerList", "List of brokers to connect").Default("kafka:9092").OverrideDefaultFromEnvar("BROKER_LIST").Strings()
)

func main() {
	log.Println("Hello, world!")

	kingpin.Parse()

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Flush.MaxMessages = 500
	config.Producer.Return.Successes = true

	prd := products.NewContext(brokerList, config)
	go prd.Start()
	defer prd.Stop()
	prd.AwaitLastOffset()
	http.HandleFunc("/product", prd.NewPDPHandler())
	http.HandleFunc("/products", prd.NewCatalogHandler())

	shutdown := checkout.StartContext(brokerList, config)
	defer shutdown()
	http.HandleFunc("/cart", checkout.CartHandler)
	http.HandleFunc("/orderCart", checkout.OrderHandler)

	log.Println("listening on port :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
