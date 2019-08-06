
export GO111MODULE=on

start:
	air -c air.conf

import: products-1.csv
	go run ./cmd/csv-import/main.go products-1.csv

update: products-1.csv products-2.csv
	go run ./cmd/csv-import/main.go products-2.csv products-1.csv

products-1.csv:
	go run ./cmd/csv-fake-create/main.go > products-1.csv

products-2.csv: products-1.csv
	cat products-1.csv | go run ./cmd/csv-fake-alternate/main.go > products-2.csv
