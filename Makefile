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
