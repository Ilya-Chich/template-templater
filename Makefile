-include .env.dev
export
COMPOSE := docker-compose --env-file .env.dev
GOLANGCI := go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.5

init: \
	docker-clean \
	up

init-prod: \
	docker-clean \
	lint \
	up


up: \
	stop \
	docker-up

stop:
	$(COMPOSE) stop

docker-clean:
	$(COMPOSE) down --volumes --remove-orphans --rmi local

docker-up:
	$(COMPOSE) up --build --remove-orphans

build:
	$(COMPOSE) build

.PHONY: lint test
lint:
	$(GOLANGCI) run

test:
	go test -race ./...