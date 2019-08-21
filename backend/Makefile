
include ./etc/help.mk
include ./etc/benchmark.mk

export UID=$(shell id -u)
export GID=$(shell id -g)

start: ##@docker Start the backend
	docker-compose up --detach --build
.PHONY: start

stop: ##@docker Stop the backend
	docker-compose stop --timeout 1
.PHONY: stop

clean: ##@docker Stop and remove the backend
	docker-compose down --remove-orphans --volumes --timeout 1
.PHONY: clean

logs: ##@docker Display the backends logs
	docker-compose logs --follow
.PHONY: logs

cli: ##@docker Access cli
	docker-compose exec backend bash
.PHONY: cli

import: ##@docker Access cli
	docker-compose exec backend make import-products
	docker-compose exec backend make import-stock
.PHONY: import



export GO111MODULE=on

run: ##@host Start the backend
	air -c air.conf

import-products: tmp/products-1.csv ##@import Import products
	go run ./cmd/import-products/main.go tmp/products-1.csv

update-products: tmp/products-1.csv products-2.csv  ##@import Update products
	go run ./cmd/import-products/main.go tmp/products-2.csv tmp/products-1.csv

import-stock: tmp/products-1-stock.csv ##@import Import stock
	go run ./cmd/import-stock/main.go tmp/products-1-stock.csv

tmp/products-1.csv:
	go run ./cmd/dev-tools/csv-fake-products/main.go -seed=0 > tmp/products-1.csv

tmp/products-1-stock.csv: tmp/products-1.csv
	cat tmp/products-1.csv | go run ./cmd/dev-tools/csv-fake-stock/main.go -seed=2 > tmp/products-1-stock.csv

tmp/products-2.csv: tmp/products-1.csv
	cat tmp/products-1.csv | go run ./cmd/dev-tools/csv-fake-product-updates/main.go -seed=1 > tmp/products-2.csv

lint: ##@host Run the go linter
	golangci-lint run -e ./vendor ./...


