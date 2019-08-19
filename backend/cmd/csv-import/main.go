package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"git.votum-media.net/event-web-store/event-web-store/backend/pkg/catalog"
	"git.votum-media.net/event-web-store/event-web-store/backend/pkg/checkout"

	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	brokerList   = kingpin.Flag("brokerList", "List of brokers to connect").Default("kafka:9092").OverrideDefaultFromEnvar("BROKER_LIST").Strings()
	verbose      = kingpin.Flag("verbose", "Verbosity").Default("false").Bool()
	currentPath  = kingpin.Arg("current", "path to current import file").Required().String()
	previousPath = kingpin.Arg("previous", "path to previous import file").Default("/dev/null").String()
)

func main() {

	kingpin.Parse()

	log.Printf("current import file %s", *currentPath)
	log.Printf("previous import file %s", *previousPath)

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Flush.MaxMessages = 500
	producer, err := sarama.NewAsyncProducer(*brokerList, config)
	if err != nil {
		log.Panicf("failed to setup the kafka producer: %s", err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Panicf("failed to close the kafka producer: %s", err)
		}
	}()

	go func() {
		for err := range producer.Errors() {
			log.Panicf("failed to send msg (key %s): %s", err.Msg.Key, err.Err)
		}
	}()

	checkoutProducer, err := sarama.NewAsyncProducer(*brokerList, config)
	if err != nil {
		log.Panicf("failed to setup the kafka checkoutProducer: %s", err)
	}
	defer func() {
		if err := checkoutProducer.Close(); err != nil {
			log.Panicf("failed to close the kafka checkoutProducer: %s", err)
		}
	}()

	go func() {
		for err := range checkoutProducer.Errors() {
			log.Panicf("failed to send msg (key %s): %s", err.Msg.Key, err.Err)
		}
	}()

	prevProducts, err := rows(*previousPath)
	if err != nil {
		log.Panicf("failed to load previous import file: %s", err)
	}
	currentProducts, err := rows(*currentPath)
	if err != nil {
		log.Panicf("failed to load current import file: %s", err)
	}

	upsert(prevProducts, currentProducts, producer.Input())
	remove(prevProducts, producer.Input())

}

func upsert(prevProducts, currentProducts map[string][]string, ch chan<- *sarama.ProducerMessage) {
	for _, currentRow := range currentProducts {

		UUID := currentRow[0]

		prevRow, prevFound := prevProducts[UUID]

		if equal(prevRow, currentRow) {
			if *verbose {
				log.Printf("skip unchanged product %s\n", UUID)
			}
			continue
		}

		if *verbose {
			ll := "insert product %s\n"
			if prevFound {
				ll = "update product %s\n"
			}
			log.Printf(ll, UUID)
		}

		curr, err := row2product(currentRow)
		if err != nil {
			log.Panicf("failed to serialize previous product %s: %s", UUID, err)
		}

		err = sendUpdate(ch, UUID, curr)
		if err != nil {
			log.Panicf("failed to send update massage: %s", err)
		}

		delete(prevProducts, UUID)
	}
}

func remove(prevProducts map[string][]string, ch chan<- *sarama.ProducerMessage) {
	for _, prevRow := range prevProducts {
		UUID := prevRow[0]
		if *verbose {
			log.Printf("delete product %s\n", UUID)
		}

		prev, err := row2product(prevRow)
		if err != nil {
			log.Panicf("failed to serialize previous product %s: %s", UUID, err)
		}
		prev.Disabled = true

		err = sendUpdate(ch, UUID, prev)
		if err != nil {
			log.Panicf("failed to send update massage: %s", err)
		}
	}
}

func sendUpdate(ch chan<- *sarama.ProducerMessage, UUID string, msg *catalog.Product) error {

	change := &checkout.CheckoutContext{
		CheckoutContextMsg: &checkout.CheckoutContext_Product{
			Product: &checkout.Product{
				Price:         msg.Price,
				ProductID:     msg.Id,
				SmallImageURL: msg.SmallImageURL,
				Title:         msg.Title,
			},
		},
	}
	bytes, err := proto.Marshal(change)
	if err != nil {
		return fmt.Errorf("failed to serialize product massage for checkout topic: %s", err)
	}
	ch <- &sarama.ProducerMessage{
		Topic: checkout.Topic,
		Key:   sarama.StringEncoder(UUID),
		Value: sarama.ByteEncoder(bytes),
	}

	bytes, err = proto.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to serialize product massage: %s", err)
	}
	ch <- &sarama.ProducerMessage{
		Topic: catalog.Topic,
		Key:   sarama.StringEncoder(UUID),
		Value: sarama.ByteEncoder(bytes),
	}
	return nil
}

func equal(a, b []string) bool {

	// If one is nil, the other must also be nil.
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func rows(path string) (map[string][]string, error) {

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open previous import file: %s", err)
	}
	defer f.Close()

	m := make(map[string][]string)
	r := csv.NewReader(f)
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		uuid := row[0]

		if _, ok := m[uuid]; ok {
			return nil, fmt.Errorf("product doublication in import list")
		}

		m[uuid] = row
	}
	return m, nil
}

func row2product(row []string) (*catalog.Product, error) {

	if row == nil {
		return nil, nil
	}

	price, err := strconv.ParseInt(row[7], 10, 64)
	if err != nil {
		return &catalog.Product{}, err
	}
	return &catalog.Product{
		Id:            row[0],
		Title:         row[1],
		Description:   row[2],
		Longtext:      row[3],
		Category:      row[4],
		SmallImageURL: row[5],
		LargeImageURL: row[6],
		Price:         price,
	}, nil
}
