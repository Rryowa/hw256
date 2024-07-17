package util

import (
	"github.com/joho/godotenv"
	"homework/internal/models/config"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func NewConfig() *config.DbConfig {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatalf("err loading: %v", err)
	}

	attempts, err := strconv.Atoi(os.Getenv("ATTEMPTS"))
	if err != nil {
		log.Fatalf("err converting ATTEMPTS: %v\n", err)
	}

	timeout, err := time.ParseDuration(os.Getenv("TIMEOUT"))
	if err != nil {
		log.Fatalf("Error parsing TIMEOUT: %v\n", err)
	}

	return &config.DbConfig{
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		DBName:   os.Getenv("POSTGRES_DB"),
		Attempts: attempts,
		Timeout:  timeout,
	}
}

func NewKafkaConfig() *config.KafkaConfig {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatalf("err loading: %v", err)
	}
	kafkaUse, err := strconv.ParseBool(os.Getenv("KAFKA_USE"))
	if err != nil {
		log.Fatalf("Error parsing KAFKA_USE: %v\n", err)
	}

	return &config.KafkaConfig{
		KafkaUse:     kafkaUse,
		KafkaBrokers: strings.Split(os.Getenv("KAFKA_BROKERS"), ","),
		KafkaTopics:  strings.Split(os.Getenv("KAFKA_TOPICS"), ","),
		KafkaGroupID: os.Getenv("KAFKA_GROUP_ID"),
	}
}

func NewGrpcConfig() *config.GrpcConfig {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatalf("err loading: %v", err)
	}
	return &config.GrpcConfig{
		Host:     os.Getenv("GRPC_HOST"),
		HttpPort: os.Getenv("GRPC_HTTP_PORT"),
		GrpcPort: os.Getenv("GRPC_PORT"),
	}
}

func DoWithTries(fn func() error, attempts int, delay time.Duration) (err error) {
	for attempts > 0 {
		if err = fn(); err != nil {
			time.Sleep(delay)
			attempts--
			continue
		}
		return nil
	}
	return
}