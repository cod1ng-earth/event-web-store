
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

lint:
	golangci-lint run -e ./vendor ./...

pkg/pb/products.pb.go: pkg/pb/products.proto
	protoc --go_out=. pkg/pb/products.proto

benchmark:
	echo "GET http://localhost:8080/product?uuid=1d9d1c51-cc0a-48e0-80b5-00d2ded65791" | vegeta attack -workers=1 -max-workers=1 -name=1workers -rate=5000 -duration=15s > ~/results.bin

hdrplot:
	vegeta report --type hdrplot ~/results.bin

plot:
	vegeta plot ~/results.bin > plot.html
	google-chrome plot.html
