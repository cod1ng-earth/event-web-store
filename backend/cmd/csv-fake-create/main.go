package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	lorem "github.com/drhodes/golorem"
	"github.com/satori/go.uuid"
)

func main() {

	seed := flag.Int64("seed", time.Now().UnixNano(), "seed for random number generator")
	rows := flag.Int("rows", 100000, "number of rows")
	flag.Parse()

	log.Printf("seed: %d\n", *seed)
	rand.Seed(*seed)

	w := csv.NewWriter(os.Stdout)

	for i := 0; i < *rows; i++ {
		err := w.Write(newRow())
		if err != nil {
			log.Panicf("error adding row: %s\n", err)
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		log.Panicf("error writing csv: %s\n", err)
	}
}

func newRow() []string {
	UUID := uuid.NewV4()
	price := float64(rand.Intn(10000)) / 100
	return []string{
		UUID.String(),
		fmt.Sprintf("%s %s %s", lorem.Word(4, 13), lorem.Word(4, 13), lorem.Word(4, 13)),
		lorem.Sentence(12, 24),
		lorem.Paragraph(3, 6),
		fmt.Sprintf("%s/%s", lorem.Word(4, 13), lorem.Word(4, 13)),
		lorem.Url(),
		lorem.Url(),
		fmt.Sprintf("%.2f", price),
	}
}
