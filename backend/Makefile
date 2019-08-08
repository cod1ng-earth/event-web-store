
export GO111MODULE=on

start:
	air -c air.conf

import: products-1.csv
	go run ./cmd/csv-import/main.go products-1.csv

update: products-1.csv products-2.csv
	go run ./cmd/csv-import/main.go products-2.csv products-1.csv

products-1.csv:
	go run ./cmd/csv-fake-create/main.go -seed=0 > products-1.csv

products-2.csv: products-1.csv
	cat products-1.csv | go run ./cmd/csv-fake-alternate/main.go -seed=1 > products-2.csv

lint:
	golangci-lint run -e ./vendor ./...

pkg/pb/products.pb.go: pkg/pb/products.proto
	protoc --go_out=. pkg/pb/products.proto

ab:
	ab -c 4 -n 200000 -k http://localhost:8080/product?uuid=e412bd86-6381-4c90-8365-3219d6115710

results.bin: $(shell find cmd pkg Makfile -type f)
	echo "GET http://localhost:8080/product?uuid=e412bd86-6381-4c90-8365-3219d6115710" | vegeta attack -workers=4 -max-workers=4 -rate=0 -duration=5s > results.bin

hdrplot: results.bin
	vegeta report --type hdrplot results.bin

plot: results.bin
	vegeta plot results.bin > plot.html
	google-chrome plot.html
