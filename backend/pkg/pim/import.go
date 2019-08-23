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

	model, free := c.read()
	defer free()

	oldProducts := make(map[string]*Product)
	for _, product := range model.products {
		if product.Disabled {
			continue
		}
		oldProducts[product.Id] = product
	}
	if verbose {
		log.Printf("loaded current products")
	}

	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("failed to open import file: %s", err)
	}
	defer f.Close()

	newProducts := make(chan *Product)
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
		_, _, err = c.logProduct(newProduct)
		if err != nil {
			log.Printf("store %v", oldProduct.Id)
			log.Fatalf("failed to store product %v: %v", newProduct, err)
		}
	}

	if verbose {
		log.Printf("updated and inserted products")
	}

	// disable Products that existed before but are missing from new
	for _, oldProduct := range oldProducts {
		log.Printf("disable %v", oldProduct.Id)
		oldProduct.Disabled = true
		_, _, err = c.logProduct(oldProduct)
		if err != nil {
			log.Fatalf("failed to disable product %v: %v", oldProduct, err)
		}
	}
	if verbose {
		log.Printf("disabled old products")
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
