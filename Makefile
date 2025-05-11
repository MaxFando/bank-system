LOCAL_BIN := $(shell pwd)/bin

appName = bank
compose = docker-compose -f docker-compose-debug.yml -p $(appName)

DB_BANK = postgresql://postgres:postgres@localhost:5432/bank?sslmode=disable
DB_BANK_MIGRATION_DSN = postgresql://postgres:postgres@localhost:5432/bank?search_path=main

install-deps:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0
	GOBIN=$(LOCAL_BIN) go install go.uber.org/mock/mockgen@v0.5.0
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.21.1

db-create-migration:
	@echo "Enter migration name:"
	@read name; \
	GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(DB_BANK_MIGRATION_DSN) $(LOCAL_BIN)/goose -dir migrations create $$name sql

db-migrate:
	GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(DB_BANK_MIGRATION_DSN) $(LOCAL_BIN)/goose -dir migrations up

db-rollback:
	GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(DB_BANK_MIGRATION_DSN) $(LOCAL_BIN)/goose -dir migrations down

.PHONY: fixtures
fixtures:
	go run cmd/fixtures/main.go -dsn=$(DB_BANK)

up: down build
	@echo "Starting app..."
	$(compose) up -d
	@echo "Docker images built and started!"
	make db-migrate

build:
	@echo "Building images"
	$(compose) build
	@echo "Docker images built!"

down:
	@echo "Stopping docker compose..."
	$(compose) down
	@echo "Done!"

down-v:
	@echo "Stopping docker compose..."
	$(compose) down -v
	@echo "Done!"


generate:
	@echo "Generating code..."
	go generate ./...
	@echo "Code generated!"

lint:
	$(LOCAL_BIN)/golangci-lint run -c .golangci.yaml

lint-fix:
	$(LOCAL_BIN)/golangci-lint run -c .golangci.yaml --fix ./...