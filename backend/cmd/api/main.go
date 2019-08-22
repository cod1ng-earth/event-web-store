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

	cat := catalog.NewContext(brokerList, config)
	go cat.Start()
	defer cat.Stop()
	cat.AwaitLastOffset()
	http.HandleFunc("/product", cat.NewPDPHandler())
	http.HandleFunc("/products", cat.NewCatalogHandler())

	ckt := checkout.NewContext(brokerList, config)
	go ckt.Start()
	defer ckt.Stop()
	ckt.AwaitLastOffset()
	http.HandleFunc("/cart", ckt.NewCartHandler())
	http.HandleFunc("/orderCart", ckt.NewOrderHandler())

	http.Handle("/metrics", promhttp.Handler())

	log.Println("listening on port :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
