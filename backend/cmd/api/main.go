package main

import (
	"log"
	"net/http"

	"github.com/cod1ng-earth/event-web-store/backend/pkg/catalog"
	"github.com/cod1ng-earth/event-web-store/backend/pkg/checkout"

	"github.com/Shopify/sarama"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	brokerList  = kingpin.Flag("brokerList", "List of brokers to connect").Default("kafka:9092").OverrideDefaultFromEnvar("BROKER_LIST").Strings()
	contextList = kingpin.Flag("contexts", "List of contexts to run").Default("all").Strings()
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

	prd := catalog.NewContext(brokerList, config)
	go prd.Start()
	defer prd.Stop()
	prd.AwaitLastOffset()
	http.HandleFunc("/product", prd.NewPDPHandler())
	http.HandleFunc("/products", prd.NewCatalogHandler())

	shutdown := checkout.StartContext(brokerList, config)
	defer shutdown()
	http.HandleFunc("/cart", checkout.CartHandler)
	http.HandleFunc("/orderCart", checkout.OrderHandler)

	http.Handle("/metrics", promhttp.Handler())

	log.Println("listening on port :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
