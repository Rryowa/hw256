# Variables
BINARY_NAME=cli
BIN_DIR=bin
CMD_DIR=cmd

all: build

build:
	@echo "Building the CLI application..."
	@mkdir -p $(BIN_DIR)
	@go build -o $(BIN_DIR)/$(BINARY_NAME) $(CMD_DIR)/main.go
	@chmod +x $(BIN_DIR)/$(BINARY_NAME)
	@echo "Build completed. Binary is located at $(BIN_DIR)/$(BINARY_NAME)"

run: build
	@echo "Running the CLI application..."
	@$(BIN_DIR)/$(BINARY_NAME)

.PHONY: all build run