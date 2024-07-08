package util

import (
	"github.com/joho/godotenv"
	"homework/internal/models"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func NewConfig() *models.Config {
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

	useKafka, err := strconv.ParseBool(os.Getenv("USE_KAFKA"))
	if err != nil {
		log.Fatalf("Error parsing USE_KAFKA: %v\n", err)
	}

	return &models.Config{
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		DBName:   os.Getenv("POSTGRES_DB"),
		Attempts: attempts,
		Timeout:  timeout,
		UseKafka: useKafka,
		Brokers:  strings.Split(os.Getenv("BROKERS"), ","),
		Topic:    os.Getenv("TOPIC"),
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