LOCAL_BIN := $(shell pwd)/bin

appName = bank
compose = docker-compose -f docker-compose-debug.yml -p $(appName)

DB_BANK = postgresql://postgres:postgres@localhost:5432/bank?sslmode=disable
DB_BANK_MIGRATION_DSN = postgresql://postgres:postgres@localhost:5432/bank?search_path=main

install-deps:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
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


# Для тестов
COVERAGE_DIR := .coverage
CURRENT_COVERAGE := $(COVERAGE_DIR)/current.out
PREVIOUS_COVERAGE := $(COVERAGE_DIR)/previous.out
SUMMARY_FILE := $(COVERAGE_DIR)/summary.txt
HTML_REPORT := $(COVERAGE_DIR)/coverage.html

# Убедиться, что директория для покрытия существует
ensure-coverage-dir:
	@mkdir -p $(COVERAGE_DIR)

# Запустить тесты и снять покрытие
test: ensure-coverage-dir
	@if [ -f $(CURRENT_COVERAGE) ]; then \
		mv $(CURRENT_COVERAGE) $(PREVIOUS_COVERAGE); \
	fi
	$(eval @_scope := $(or $(addprefix './',$(filter-out $@,$(MAKECMDGOALS))), './...'))
	$(info $(M) running tests for $(@_scope))
	go test -v -race ./internal/core/bank/service/... -coverprofile=$(CURRENT_COVERAGE)
	$(info $(M) Tests completed. Coverage saved to $(CURRENT_COVERAGE).)

# Сравнить текущее покрытие с предыдущим
compare: test
	@if [ -f $(PREVIOUS_COVERAGE) ]; then \
		PREVIOUS_PERCENT=$$(go tool cover -func=$(PREVIOUS_COVERAGE) | grep "total:" | awk '{print $$3}' | sed 's/%//'); \
		CURRENT_PERCENT=$$(go tool cover -func=$(CURRENT_COVERAGE) | grep "total:" | awk '{print $$3}' | sed 's/%//'); \
		echo "Previous coverage: $$PREVIOUS_PERCENT%"; \
		echo "Current coverage: $$CURRENT_PERCENT%"; \
		if [ $$(echo "$$CURRENT_PERCENT < $$PREVIOUS_PERCENT" | bc -l) -eq 1 ]; then \
			echo "Warning: Coverage decreased! ($$PREVIOUS_PERCENT% -> $$CURRENT_PERCENT%)"; \
			echo "Coverage decreased! ($$PREVIOUS_PERCENT% -> $$CURRENT_PERCENT%)" > $(SUMMARY_FILE); \
		else \
			echo "Coverage maintained or improved. ($$PREVIOUS_PERCENT% -> $$CURRENT_PERCENT%)"; \
			echo "Coverage maintained or improved. ($$PREVIOUS_PERCENT% -> $$CURRENT_PERCENT%)" > $(SUMMARY_FILE); \
		fi; \
	else \
		echo "No previous coverage to compare."; \
		CURRENT_PERCENT=$$(go tool cover -func=$(CURRENT_COVERAGE) | grep "total:" | awk '{print $$3}' | sed 's/%//'); \
		echo "Current coverage: $$CURRENT_PERCENT%" > $(SUMMARY_FILE); \
	fi
	@cat $(SUMMARY_FILE)

# Сгенерировать HTML отчет
coverage: compare
	@go tool cover -html=$(CURRENT_COVERAGE) -o $(HTML_REPORT)
	@echo "HTML coverage report generated: $(HTML_REPORT)"