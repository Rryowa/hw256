# Variables
BINARY_NAME=cli
BIN_DIR=bin
CMD_DIR=cmd
EXPLAIN_DIR=explain

build:
	@echo "Building the CLI application..."
	@mkdir -p $(BIN_DIR)
	@go build -o $(BIN_DIR)/$(BINARY_NAME) $(CMD_DIR)/main.go
	@chmod +x $(BIN_DIR)/$(BINARY_NAME)
	@echo "Build completed. Binary is located at $(BIN_DIR)/$(BINARY_NAME)"

run: build
	@echo "Running the CLI application..."
	@$(BIN_DIR)/$(BINARY_NAME)

explain:
	@mkdir -p $(BIN_DIR)
	@go build -o $(BIN_DIR)/$(EXPLAIN_DIR) $(EXPLAIN_DIR)/$(EXPLAIN_DIR).go
	@chmod +x $(BIN_DIR)/$(EXPLAIN_DIR)
	@echo "Analysis..."
	@$(BIN_DIR)/$(EXPLAIN_DIR)

.PHONY: build run explain