package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"git.votum-media.net/event-web-store/event-web-store/backend/pkg/pb"
	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	brokerList  = kingpin.Flag("brokerList", "List of brokers to connect").Default("localhost:9092").OverrideDefaultFromEnvar("BROKER_LIST").Strings()
	currentPath = kingpin.Arg("current", "path to current import file").Required().String()
)

func main() {

	kingpin.Parse()

	log.Printf("current import file %s", *currentPath)

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

	newStock, err := rows(*currentPath)
	if err != nil {
		log.Panicf("failed to load current import file: %s", err)
	}

	importStock(newStock, producer.Input())
}

func importStock(stocks map[string][]string, ch chan<- *sarama.ProducerMessage) {
	for _, row := range stocks {

		stock, err := row2stock(row)
		if err != nil {
			log.Panicf("failed to parse stock: %s", err)
		}

		err = sendUpdate(ch, stock)
		if err != nil {
			log.Panicf("failed to send update massage: %s", err)
		}
	}
}

func sendUpdate(ch chan<- *sarama.ProducerMessage, msg *pb.Stock) error {
	change := &pb.CheckoutContext{
		CheckoutContext: &pb.CheckoutContext_Stock{
			Stock: msg,
		},
	}
	bytes, err := proto.Marshal(change)
	if err != nil {
		return fmt.Errorf("failed to serialize product delete massage: %s", err)
	}
	ch <- &sarama.ProducerMessage{
		Topic: "checkout",
		Key:   sarama.StringEncoder(msg.Uuid),
		Value: sarama.ByteEncoder(bytes),
	}
	return nil
}

func rows(path string) (map[string][]string, error) {

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open import file: %s", err)
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
		m[uuid] = row
	}
	return m, nil
}

func row2stock(row []string) (*pb.Stock, error) {

	if row == nil {
		return nil, nil
	}

	quantity, err := strconv.ParseInt(row[1], 10, 64)
	if err != nil {
		return &pb.Stock{}, fmt.Errorf("stock quantity can not be parsed '%v': %v", row, err)
	}
	return &pb.Stock{
		Uuid:     row[0],
		Quantity: quantity,
	}, nil
}
