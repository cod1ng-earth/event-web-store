package pim

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"reflect"
	"strconv"
)

func (c *context) ImportFile(path string, verbose bool) {

	log.Print("ImportFile")

	model, free := c.read()

	log.Print("accessing model")

	oldProducts := make(map[string]*Product)
	for _, product := range model.products {
		if product.Disabled {
			continue
		}
		oldProducts[product.Id] = product
	}
	free()
	if verbose {
		log.Printf("loaded current products")
	}

	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("failed to open import file: %s", err)
	}
	defer f.Close()

	producer, err := c.newSyncProducer()
	if err != nil {
		log.Fatalf("failed send messages to kafka: %s", err)
	}

	closed := producer.awaitClose(func(err error) {
		log.Fatalf("failure to write to kafka: %s", err)
	})

	newProducts := make(chan *Product, 1000)
	go parseProducts(f, newProducts)
	for newProduct := range newProducts {

		oldProduct, found := oldProducts[newProduct.Id]

		// remove from "existed before" list
		if found {
			delete(oldProducts, oldProduct.Id)
		}

		// if the same -> skip
		if found && reflect.DeepEqual(oldProduct, newProduct) {
			continue
		}

		// if new or changed -> upsert
		producer.appendProduct(newProduct)
	}
	if verbose {
		log.Printf("updated and inserted products")
	}

	for _, oldProduct := range oldProducts {
		oldProduct.Disabled = true
		producer.appendProduct(oldProduct)
	}
	if verbose {
		log.Printf("disabled old products")
	}

	producer.AsyncClose()

	closed.Wait()
	if verbose {
		log.Printf("received ACKs for all messages from kafka")
	}
}

func parseProducts(r io.Reader, ch chan *Product) {

	csv := csv.NewReader(r)
	for {
		row, err := csv.Read()
		if err == io.EOF {
			close(ch)
			return
		}

		price, err := strconv.ParseInt(row[7], 10, 64)
		if err != nil {
			log.Printf("failed to parse price in row %v: %v", row, err)
			continue
		}

		ch <- &Product{
			Id:            row[0],
			Name:          row[1],
			Description:   row[2],
			Longtext:      row[3],
			Category:      row[4],
			SmallImageURL: row[5],
			LargeImageURL: row[6],
			Price:         price,
		}
	}
}
