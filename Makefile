.DEFAULT_GOAL := help
help: ### Display this help text.
	@awk -F ':.*###\ ' '$$0 ~ FS {printf "    %-8s%s\n", $$1 ": ", $$2}' $(MAKEFILE_LIST) | sort
.PHONY: help

setup: clean ### Stop and remove all containers, networks and volumes. (Re)build all images. Import products into Kafka.
	docker-compose build
	docker-compose up --detach kafka
	docker-compose run --rm backend bash -euo pipefail -c 'echo "Waiting for Kafka... "; while ! nc -z kafka 9093; do sleep 1; done; echo "Kafka is ready!"'
	docker-compose run --rm backend bash -euo pipefail -c 'make products-1.csv && make import && make products-1-stock.csv && make stock'
.PHONY: setup

start: ### Start all containers, networks and volumes in the background.
	docker-compose up --detach --renew-anon-volumes
.PHONY: start

stop: ### Stop all containers, networks, and volumes.
	docker-compose down --remove-orphans
.PHONY: stop

clean: ### Stop and remove all containers, images, networks, and volumes.
	docker-compose down --remove-orphans --rmi all --volumes
.PHONY: clean

logs: ### Display container logs.
	docker-compose logs --follow
.PHONY: logs
