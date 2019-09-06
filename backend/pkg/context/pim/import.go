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

	log.Printf("importing %s", path)
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

	producer, err := c.newSyncProducer(func(err error) {
		log.Fatalf("failure to write to kafka: %s", err)
	})
	if err != nil {
		log.Fatalf("failed send messages to kafka: %s", err)
	}

	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("failed to open import file: %s", err)
	}
	defer f.Close()
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
		producer.logProduct(newProduct)
	}
	if verbose {
		log.Printf("updated and inserted products")
	}

	for _, oldProduct := range oldProducts {
		oldProduct.Disabled = true
		producer.logProduct(oldProduct)
	}
	if verbose {
		log.Printf("disabled old products")
	}

	producer.Close()
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
		if err != nil {
			log.Fatalf("failed to import csv file: %v", err)
		}

		price, err := strconv.ParseInt(row[7], 10, 64)
		if err != nil {
			log.Fatalf("failed to parse price in row %v: %v", row, err)
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
