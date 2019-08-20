package main

import (
	"encoding/csv"
	"flag"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func main() {

	seed := flag.Int64("seed", time.Now().UnixNano(), "seed for random number generator")
	outOfStock := flag.Int("out-of-stock", 10, "probability a product is out of stock, in %")
	maxStock := flag.Int("max-stock", 30, "maximum stock quantity")
	flag.Parse()

	log.Printf("seed: %d\n", *seed)
	rand.Seed(*seed)

	r := csv.NewReader(os.Stdin)
	w := csv.NewWriter(os.Stdout)

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}

		if rand.Intn(100) < *outOfStock {
			continue
		}

		record = createStock(record, *maxStock)

		err = w.Write(record)
		if err != nil {
			log.Panicf("error keeping row: %s\n", err)
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		log.Fatalf("error writing csv: %s\n", err)
	}
}

func createStock(in []string, maxStock int) []string {
	return []string{
		in[0],
		strconv.Itoa(rand.Intn(maxStock-1) + 1),
	}
}
