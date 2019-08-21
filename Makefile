
include ./shared/etc/help.mk
include ./shared/etc/chaos-network.mk

setup: ##@docker Start all conponents and import products
	$(MAKE) start
	cd backend && make import
.PHONY: setup

start: ## Start all components
	cd shared && make start
	cd backend && make start
	cd frontend && make start
.PHONY: start

stop: ## Stop all components
	cd frontend && make stop
	cd backend && make stop
	cd shared && make stop
.PHONY: stop

clean: ## Stop and remove all components
	cd frontend && make clean
	cd backend && make clean
	cd shared && make clean
.PHONY: clean
