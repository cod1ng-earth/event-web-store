package main

import (
	"log"
	"net/http"

	"github.com/cod1ng-earth/event-web-store/backend/pkg/context/catalog"
	"github.com/cod1ng-earth/event-web-store/backend/pkg/context/checkout"

	"github.com/Shopify/sarama"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	log.Println("Hello, world!")

	brokers := kingpin.Flag("broker", "kafka broker to connect").Default("kafka:9092").OverrideDefaultFromEnvar("BROKER").Strings()
	//contexts := kingpin.Flag("context", "contexts to run").Strings()
	kingpin.Parse()

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Flush.MaxMessages = 500
	config.Producer.Return.Successes = true

	cat := catalog.NewContext(brokers, config)
	go cat.Start()
	defer cat.Stop()
	cat.AwaitLastOffset()
	http.HandleFunc("/product", cat.NewPDPHandler())
	http.HandleFunc("/products", cat.NewCatalogHandler())

	ckt := checkout.NewContext(brokers, config)
	go ckt.Start()
	defer ckt.Stop()
	ckt.AwaitLastOffset()
	http.HandleFunc("/cart", ckt.NewCartHandler())
	http.HandleFunc("/orderCart", ckt.NewOrderHandler())

	http.Handle("/metrics", promhttp.Handler())

	// TODO wait for brides to be up to date
	// TODO wait for updater to be up to date
	// TODO do not leak pim offset to frontend

	log.Println("listening on port :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
