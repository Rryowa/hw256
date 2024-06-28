include .env
BINARY_NAME=cli
BIN_DIR=bin
CMD_DIR=cmd
EXPLAIN_DIR=explain
DB_STRING="postgres://$(POSTGRES_USER):${POSTGRES_PASSWORD}@${POSTGRES_HOST}/${POSTGRES_DB}?sslmode=disable"
TEST_STRING="postgres://$(TEST_USER):${TEST_PASSWORD}@${TEST_HOST}/${TEST_DB}?sslmode=disable"
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

.PHONY: test
test:
	@echo "Testing..."
#	Single stage build
#	@go test ./tests -tags=integration

#	Multi-stage build
#	@./tests/repo_test wont work if testconfig asks for ../.env
#	use cd tests && ./repo_test instead

#	if testconfig will ask for ./.env
	@./bin/repo_test
#	will work

.PHONY: up
up:
	@goose -dir $(MIGRATIONS_DIR) postgres $(DB_STRING) up
	@goose -dir $(MIGRATIONS_DIR) postgres $(TEST_STRING) up

.PHONY: down
down:
	@goose -dir $(MIGRATIONS_DIR) postgres $(DB_STRING) down
	@goose -dir $(MIGRATIONS_DIR) postgres $(TEST_STRING) down

.PHONY: compose-up
compose-up:
	@docker compose up app db db_test -d
	@docker exec -it go bash

.PHONY: compose-rm
compose-rm:
	@docker compose rm app db db_test -fvs