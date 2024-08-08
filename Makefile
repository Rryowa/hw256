include .env
include tests/test.env
DB_STRING="postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable"
TEST_STRING="postgres://$(TEST_USER):$(TEST_PASSWORD)@$(TEST_HOST):$(TEST_PORT)/$(TEST_DB)?sslmode=disable"
MIGRATIONS_DIR="./migrations"

build:
	@mkdir -p bin
	@go build -o bin/cli cmd/main.go
	@chmod +x bin/cli

run:
	@bin/cli

all: build up run

up:
	@goose -dir $(MIGRATIONS_DIR) postgres $(DB_STRING) up

down:
	@goose -dir $(MIGRATIONS_DIR) postgres $(DB_STRING) down

up-d:
	docker-compose up -d

down-d:
	docker-compose down -v

up-test:
	@goose -dir $(MIGRATIONS_DIR) postgres $(TEST_STRING) up

down-test:
	@goose -dir $(MIGRATIONS_DIR) postgres $(TEST_STRING) down

exec-pg:
	@docker exec -it pg psql -U postgres


######################################################################


# Define paths and binaries
LOCAL_BIN:=$(CURDIR)/bin
PROTOC = PATH="$$PATH:$(LOCAL_BIN)" protoc
ORDERS_PROTO_PATH:="api/proto/orders/v1"
PROTOC_GEN_GO := $(LOCAL_BIN)/protoc-gen-go
PROTOC_GEN_GO_GRPC := $(LOCAL_BIN)/protoc-gen-go-grpc
PROTOC_GEN_GRPC_GATEWAY := $(LOCAL_BIN)/protoc-gen-grpc-gateway
PROTOC_GEN_OPENAPIV2 := $(LOCAL_BIN)/protoc-gen-openapiv2
PROTOC_GEN_VALIDATE := $(LOCAL_BIN)/protoc-gen-validate

# Check and install binary dependencies if not present
.PHONY: .bin-deps
.bin-deps: $(PROTOC_GEN_GO) $(PROTOC_GEN_GO_GRPC) $(PROTOC_GEN_GRPC_GATEWAY) $(PROTOC_GEN_OPENAPIV2) $(PROTOC_GEN_VALIDATE)

$(PROTOC_GEN_GO):
	$(info Installing protoc-gen-go...)
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

$(PROTOC_GEN_GO_GRPC):
	$(info Installing protoc-gen-go-grpc...)
	GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

$(PROTOC_GEN_GRPC_GATEWAY):
	$(info Installing protoc-gen-grpc-gateway...)
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest

$(PROTOC_GEN_OPENAPIV2):
	$(info Installing protoc-gen-openapiv2...)
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest

$(PROTOC_GEN_VALIDATE):
	$(info Installing protoc-gen-validate...)
	GOBIN=$(LOCAL_BIN) go install github.com/envoyproxy/protoc-gen-validate@latest

# Вендоринг внешних proto файлов
.vendor-proto: vendor-proto/google vendor-proto/protoc-gen-openapiv2/options vendor-proto/validate

# Устанавливаем proto описания protoc-gen-openapiv2/options
vendor-proto/protoc-gen-openapiv2/options:
	if [ -d "vendor.proto/protoc-gen-openapiv2/options" ]; then \
		rm -rf vendor.proto/protoc-gen-openapiv2/options; \
	fi
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
 		https://github.com/grpc-ecosystem/grpc-gateway vendor.proto/grpc-ecosystem && \
 	cd vendor.proto/grpc-ecosystem && \
	git sparse-checkout set --no-cone protoc-gen-openapiv2/options && \
	git checkout
	mkdir -p vendor.proto/protoc-gen-openapiv2
	mv vendor.proto/grpc-ecosystem/protoc-gen-openapiv2/options vendor.proto/protoc-gen-openapiv2
	rm -rf vendor.proto/grpc-ecosystem


# Устанавливаем proto описания google/protobuf
vendor-proto/google:
	if [ -d "vendor.proto/google" ]; then \
    	rm -rf vendor.proto/google; \
	fi
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
		https://github.com/protocolbuffers/protobuf vendor.proto/protobuf &&\
	cd vendor.proto/protobuf &&\
	git sparse-checkout set --no-cone src/google/protobuf &&\
	git checkout
	mkdir -p vendor.proto/google
	mv vendor.proto/protobuf/src/google/protobuf vendor.proto/google
	rm -rf vendor.proto/protobuf

	git clone -b master --single-branch -n --depth=1 --filter=tree:0 \
		https://github.com/googleapis/googleapis vendor.proto/googleapis && \
	cd vendor.proto/googleapis && \
	git sparse-checkout set --no-cone google/api && \
	git checkout
	mv vendor.proto/googleapis/google/api vendor.proto/google
	rm -rf vendor.proto/googleapis

vendor-proto/validate:
	if [ -d "vendor.proto/validate" ]; then \
    	rm -rf vendor.proto/validate; \
    fi
	git clone -b main --single-branch --depth=2 --filter=tree:0 \
		https://github.com/bufbuild/protoc-gen-validate vendor.proto/tmp && \
		cd vendor.proto/tmp && \
		git sparse-checkout set --no-cone validate &&\
		git checkout
		mkdir -p vendor.proto/validate
		mv vendor.proto/tmp/validate vendor.proto/
		rm -rf vendor.proto/tmp


.PHONY: generate
generate: .bin-deps .vendor-proto
	mkdir -p pkg/${ORDERS_PROTO_PATH}
	protoc -I api/proto \
		-I vendor.proto \
		${ORDERS_PROTO_PATH}/orders.proto \
		--plugin=protoc-gen-go=$(LOCAL_BIN)/protoc-gen-go --go_out=./pkg/${ORDERS_PROTO_PATH} --go_opt=paths=source_relative\
		--plugin=protoc-gen-go-grpc=$(LOCAL_BIN)/protoc-gen-go-grpc --go-grpc_out=./pkg/${ORDERS_PROTO_PATH} --go-grpc_opt=paths=source_relative \
		--plugin=protoc-gen-grpc-gateway=$(LOCAL_BIN)/protoc-gen-grpc-gateway --grpc-gateway_out ./pkg/api/proto/orders/v1  --grpc-gateway_opt  paths=source_relative --grpc-gateway_opt generate_unbound_methods=true \
		--plugin=protoc-gen-openapiv2=$(LOCAL_BIN)/protoc-gen-openapiv2 --openapiv2_out=./pkg/api/proto/orders/v1 \
		--plugin=protoc-gen-validate=$(LOCAL_BIN)/protoc-gen-validate --validate_out="lang=go,paths=source_relative:pkg/api/proto/orders/v1"