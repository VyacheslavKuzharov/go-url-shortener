#!make
include .env
export $(shell sed 's/=.*//' .env)
export DATE=$(shell date +%Y-%m-%d)
#SHELL:=/bin/bash
EXEC = docker-compose exec
RUN = docker-compose run --rm
START = docker-compose up -d
STOP = docker-compose stop
LOGS = docker-compose logs

.PHONY: help

help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

status: ### Containers statuses
	@echo "\n\033[01;33m Containers statuses \033[0m"
	@docker-compose ps
.PHONY: status

up: ### Up docker-compose
	@echo "\n\033[0;33m Spinning up docker environment... \033[0m"
	@$(START)
	@$(MAKE) --no-print-directory status
.PHONY: up

go-run-pg: ### Run Go code with pg storage
	go run cmd/shortener/main.go -d '$(PG_URL)?sslmode=disable'
.PHONY: go-run-pg

go-run: ### Run Go code with in memory storage
	go run cmd/shortener/main.go
.PHONY: go-run

stop: ### Stop docker-compose
	@echo "\n\033[0;33m Halting containers... \033[0m"
	@docker-compose stop
	@$(MAKE) --no-print-directory status
.PHONY: stop

down: ### Destroy containers
	@echo "\n\033[0;33m Destroy containers... \033[0m"
	@docker-compose down
	@$(MAKE) --no-print-directory status
.PHONY: down

MIGRATION_NAME := $(or $(MIGRATION_NAME),migration_name)
migrate-create:  ### create new migration. With specific name: MIGRATION_NAME="some_name"
	migrate create -ext sql -dir migrations $(MIGRATION_NAME)
.PHONY: db-migrate-create

migrate-up: ### migration up
	migrate -path migrations -database '$(PG_URL)?sslmode=disable' up
.PHONY: db-migrate-up

migrate-down: ### migration down
	migrate -path migrations -database '$(PG_URL)?sslmode=disable' down
.PHONY: db-migrate-down