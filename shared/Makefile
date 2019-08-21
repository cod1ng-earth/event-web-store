
include ./etc/help.mk
include ./etc/chaos-network.mk

start: ##@docker Start kafka and zookeeper in the background
	docker-compose up --detach --build
.PHONY: start

stop: ##@docker Stop kafka and zookeeper
	docker-compose stop --timeout 1
.PHONY: stop

clean: ##@docker Stop and remove kafka and zookeeper
	docker-compose down --remove-orphans --volumes --timeout 1
.PHONY: clean

logs: ##@docker Display kafka and zookeepers logs
	docker-compose logs --follow
.PHONY: logs

run: ##@docker Run kafka and zookeeper
	docker-compose up --detach --build
.PHONY: start
