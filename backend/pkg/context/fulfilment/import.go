package fulfilment

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
)

type stock struct {
	productId string
	quantity  int64
}

func (c *context) ImportFile(path string, verbose bool) {

	model, free := c.aggregator.read()
	defer free()

	inventory := model.inventory

	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("failed to open import file: %s", err)
	}
	defer f.Close()

	newStock := make(chan stock)
	go parseStocks(f, newStock)

	for stock := range newStock {

		correction := &StockCorrected{
			ProductID:      stock.productId,
			QuantityChange: stock.quantity,
		}

		quantity, found := inventory[stock.productId]
		if found {
			correction.QuantityChange -= quantity
		}

		if correction.QuantityChange == 0 {
			continue
		}

		_, _, err = c.internalTopic.logStockCorrected(correction)
		if err != nil {
			log.Fatalf("failed to store stock correction %v: %v", stock, err)
		}
	}

	if verbose {
		log.Printf("updated stock")
	}
}

func parseStocks(r io.Reader, ch chan stock) {

	csv := csv.NewReader(r)
	for {
		row, err := csv.Read()
		if err == io.EOF {
			close(ch)
			return
		}

		quantity, err := strconv.ParseInt(row[1], 10, 64)
		if err != nil {
			log.Printf("stock quantity can not be parsed '%v': %v", row, err)
			continue
		}
		ch <- stock{
			productId: row[0],
			quantity:  quantity,
		}
	}
}
