package main

import (
	"log"

	"github.com/Shopify/sarama"
	"github.com/cod1ng-earth/event-web-store/backend/pkg/fulfilment"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	brokerList  = kingpin.Flag("brokerList", "List of brokers to connect").Default("kafka:9092").OverrideDefaultFromEnvar("BROKER_LIST").Strings()
	verbose     = kingpin.Flag("verbose", "Verbosity").Default("true").Bool()
	currentPath = kingpin.Arg("path", "path to import file").Required().String()
)

func main() {

	kingpin.Parse()

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Flush.MaxMessages = 500
	config.Producer.Return.Successes = true

	if *verbose {
		log.Printf("import file = %s", *currentPath)
	}

	fulfilment := fulfilment.NewContext(brokerList, config)
	go fulfilment.Start()
	defer fulfilment.Stop()
	fulfilment.AwaitLastOffset()

	if *verbose {
		log.Printf("context started up")
	}

	fulfilment.ImportFile(*currentPath, *verbose)
}
