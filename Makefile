.DEFAULT_GOAL := help
help: ### Display this help text.
	@awk -F ':.*###\ ' '$$0 ~ FS {printf "    %-14s%s\n", $$1 ": ", $$2}' $(MAKEFILE_LIST) | sort
.PHONY: help

setup: ### Stop and remove all containers, networks and volumes. (Re)build all images. Import products into Kafka.
	docker-compose build
	touch frontend/dist/main.js
	chmod o+w frontend/dist/main.js

	$(MAKE) start
##	docker-compose exec backend make import
##	docker-compose exec backend make stock
.PHONY: setup

cli-frontend: ### Access backend cli
	docker-compose exec frontend bash

cli-backend: ### Access backend cli
	docker-compose exec backend bash

start: ### Start all containers, networks and volumes in the background.
	docker-compose up --detach
.PHONY: start

stop: ### Stop all containers, networks, and volumes.
	docker-compose stop
.PHONY: stop

clean: ### Stop and remove all containers, images, networks, and volumes.
	docker-compose down --remove-orphans --volumes --timeout 1
.PHONY: clean

logs: ### Display container logs.
	docker-compose logs --follow
.PHONY: logs


protoc: backend/pkg/catalog/catalog.pb.go backend/pkg/checkout/checkout.pb.go frontend/src/Catalog.elm frontend/src/Checkout.elm # backend/pkg/checkout/catalog.pb.go

backend/pkg/catalog/catalog.pb.go: backend/pkg/catalog/catalog.proto
	cd backend/pkg/catalog/ && protoc --go_out=. catalog.proto

backend/pkg/checkout/checkout.pb.go: backend/pkg/checkout/checkout.proto backend/pkg/catalog/catalog.proto
	cd backend/pkg/checkout/ && protoc --go_out=. --proto_path=../../../backend/pkg/checkout/ --proto_path=../../../backend/pkg/catalog/ checkout.proto

#backend/pkg/checkout/catalog.pb.go: backend/pkg/catalog/catalog.proto
#	cd backend/pkg/checkout/ && protoc --go_out=. --proto_path=../../../backend/pkg/checkout/ --proto_path=../../../backend/pkg/catalog/ catalog.proto

frontend/src/Catalog.elm: backend/pkg/catalog/catalog.proto
	protoc --elm_out=frontend/src --proto_path=backend/pkg/catalog/ catalog.proto

frontend/src/Checkout.elm: backend/pkg/checkout/checkout.proto
	protoc --elm_out=frontend/src --proto_path=backend/pkg/checkout/ --proto_path=backend/pkg/catalog/ checkout.proto
