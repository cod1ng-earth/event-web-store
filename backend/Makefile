# thanks to https://gist.github.com/prwhite/8168133#gistcomment-1727513 for 'make help'

#COLORS
GREEN  := $(shell tput -Txterm setaf 2)
WHITE  := $(shell tput -Txterm setaf 7)
YELLOW := $(shell tput -Txterm setaf 3)
RESET  := $(shell tput -Txterm sgr0)

# Add the following 'help' target to your Makefile
# And add help text after each target name starting with '\#\#'
# A category can be added with @category
HELP_FUN = \
	%help; \
	while(<>) { push @{$$help{$$2 // 'options'}}, [$$1, $$3] if /^([a-zA-Z\-]+)\s*:.*\#\#(?:@([a-zA-Z\-]+))?\s(.*)$$/ }; \
	print "usage: make [target]\n\n"; \
	for (sort keys %help) { \
	print "${WHITE}$$_:${RESET}\n"; \
	for (@{$$help{$$_}}) { \
	$$sep = " " x (32 - length $$_->[0]); \
	print "  ${YELLOW}$$_->[0]${RESET}$$sep${GREEN}$$_->[1]${RESET}\n"; \
	}; \
	print "\n"; }

.PHONY: help
help: ##@other Show this help.
	@perl -e '$(HELP_FUN)' $(MAKEFILE_LIST)




export GO111MODULE=on

URL = "http://localhost:8080/product?uuid=4aa3dfe0-0f9a-4a37-8efc-52f0e6a4f56b"

start: ##@develop Start the backend
	air -c air.conf

import: products-1.csv ##@import Import products
	go run ./cmd/csv-import/main.go ./tmp/products-1.csv

update: products-1.csv products-2.csv  ##@import Update products
	go run ./cmd/csv-import/main.go ./tmp/products-2.csv ./tmp/products-1.csv

stock: products-1-stock.csv ##@import Import stock
	go run ./cmd/stock-import/main.go ./tmp/products-1-stock.csv

products-1.csv: cmd/csv-fake-create/main.go
	go run ./cmd/csv-fake-create/main.go -seed=0 > ./tmp/products-1.csv

products-1-stock.csv: products-1.csv cmd/csv-fake-stock/main.go
	cat ./tmp/products-1.csv | go run ./cmd/csv-fake-stock/main.go -seed=2 > ./tmp/products-1-stock.csv

products-2.csv: products-1.csv cmd/csv-fake-alternate/main.go
	cat ./tmp/products-1.csv | go run ./cmd/csv-fake-alternate/main.go -seed=1 > ./tmp/products-2.csv

lint: ##@develop Run the go linter
	golangci-lint run -e ./vendor ./...

pkg/pb/products.pb.go: pkg/pb/products.proto
	protoc --go_out=. pkg/pb/products.proto

curl: ##@benchmark show json
	curl -v --silent ${URL} | jq .

warmup: ##@benchmark warmup process
	ab -c 1 -n 1000 -k ${URL}

ab: ##@benchmark show throughput
	ab -c 4 -n 200000 -k ${URL}

results.bin: $(shell find cmd pkg -type f) Makefile
	echo "GET ${URL}" | vegeta attack -workers=4 -max-workers=4 -rate=0 -duration=5s > results.bin

hdrplot: results.bin ##@benchmark show latency histogram
	vegeta report -type hdrplot results.bin

plot: results.bin ##@benchmark show latency ploy
	vegeta plot results.bin > plot.html
	google-chrome plot.html

report: results.bin ##@benchmark show vegeta report
	vegeta report -type text results.bin

break: ##@network Slow down local network
	sudo tc qdisc add dev lo root netem delay 14ms

fix: ##@network Stop slowing down local network
	sudo tc qdisc delete dev lo root netem delay 14ms
