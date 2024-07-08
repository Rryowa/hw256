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
	docker-compose up -d db db_test zookeeper kafka1 kafka2 kafka3

down-d:
	docker-compose down -v

up-test:
	@goose -dir $(MIGRATIONS_DIR) postgres $(TEST_STRING) up

down-test:
	@goose -dir $(MIGRATIONS_DIR) postgres $(TEST_STRING) down

exec-pg:
	@docker exec -it pg psql -U postgres

#	kafka-topics.sh --create --topic orders --partitions 3 --replication-factor 1 --bootstrap-server localhost:9091
#	kafka-topics.sh --list --bootstrap-server localhost:9091
#	kafka-topics.sh --delete --topic orders --bootstrap-server localhost:9091
#Number of partitions:
#If you have a small cluster < 6 brokers, create 3x of brokers you have, else 2x,
#cause if you have more brokers over time, you will have enough partitions to cover that.
#If you need 20 consumers at peak time, you need at least 20 partitions in your topic.
#Replication factor:
#It should be at least 2 and a maximum of 4. The recommended number is 3 as it provides
#the right balance between performance and fault tolerance.