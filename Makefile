.DEFAULT_GOAL := help
help: ### Display this help text.
	@awk -F ':.*###\ ' '$$0 ~ FS {printf "    %-8s%s\n", $$1 ": ", $$2}' $(MAKEFILE_LIST) | sort
.PHONY: help

start: ### Build and start all containers, networks, and volumes in the background.
	docker-compose up --build --detach
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
