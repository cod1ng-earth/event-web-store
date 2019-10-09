package main

import (
	"log"
	"time"

	"github.com/Shopify/sarama"
	"github.com/cod1ng-earth/event-web-store/backend/pkg/context/pim"
	"github.com/pkg/profile"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	brokerList  = kingpin.Flag("brokerList", "List of brokers to connect").Default("kafka:9092").OverrideDefaultFromEnvar("BROKER_LIST").Strings()
	verbose     = kingpin.Flag("verbose", "Verbosity").Default("true").Bool()
	currentPath = kingpin.Arg("path", "path to import file").Required().String()
)

func main() {

	_ = profile.MemProfile

	//defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()
	// defer profile.Start(profile.MemProfile, profile.ProfilePath(".")).Stop()
	//	defer profile.Start(profile.MemProfile, profile.MemProfileRate(1), profile.ProfilePath(".")).Stop()
	// defer profile.Start(profile.TraceProfile, profile.ProfilePath(".")).Stop()

	kingpin.Parse()

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.IsolationLevel = sarama.ReadCommitted

	// config.Producer.Flush.Messages = 5
	//	config.Producer.Flush.MaxMessages = 1000
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 10
	config.Producer.Retry.BackoffFunc = func(retries, maxRetries int) time.Duration {
		return time.Duration(retries) * 500 * time.Millisecond
	}
	config.Producer.Idempotent = true
	config.Version = sarama.V2_3_0_0
	config.Net.MaxOpenRequests = 1

	if *verbose {
		log.Printf("import file = %s", *currentPath)
	}

	pim := pim.NewContext(brokerList, config)
	go pim.Start()
	defer pim.Stop()

	pim.ImportFile(*currentPath, *verbose)
}
