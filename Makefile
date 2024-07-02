include .env
include tests/test.env
BINARY_NAME=cli
BIN_DIR=bin
CMD_DIR=cmd
EXPLAIN_DIR=explain
DB_STRING="postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable"
TEST_STRING="postgres://$(TEST_USER):$(TEST_PASSWORD)@$(TEST_HOST):$(TEST_PORT)/$(TEST_DB)?sslmode=disable"
MIGRATIONS_DIR="./migrations"

.PHONY: build
build:
	@echo "Building the CLI application..."
	@mkdir -p $(BIN_DIR)
	@go build -o $(BIN_DIR)/$(BINARY_NAME) $(CMD_DIR)/main.go
	@chmod +x $(BIN_DIR)/$(BINARY_NAME)
	@echo "Build completed. Binary is located at $(BIN_DIR)/$(BINARY_NAME)"

.PHONY: run
run:
	@echo "Running the CLI application..."
	@$(BIN_DIR)/$(BINARY_NAME)

.PHONY: all
all: build up run

.PHONY: test-int
test-int:
	@go test ./tests -tags=integration

.PHONY: up
up:
	@goose -dir $(MIGRATIONS_DIR) postgres $(DB_STRING) up

.PHONY: down
down:
	@goose -dir $(MIGRATIONS_DIR) postgres $(DB_STRING) down

.PHONY: up-test
up-test:
	@goose -dir $(MIGRATIONS_DIR) postgres $(TEST_STRING) up

.PHONY: down-test
down-test:
	@goose -dir $(MIGRATIONS_DIR) postgres $(TEST_STRING) down

.PHONY: compose-up
compose-up:
	@docker compose up -d
#	@docker exec -it go bash

.PHONY: compose-rm
compose-rm:
	@docker compose rm -fvs

.PHONY: exec-pg
exec-pg:
	@docker exec -it pg psql -U postgres