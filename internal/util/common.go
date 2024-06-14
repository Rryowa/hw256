package util

import (
	"github.com/joho/godotenv"
	"homework-1/internal/models"
	"log"
	"os"
	"strconv"
	"time"
)

func NewConfig() *models.Config {
	err := godotenv.Load("internal/util/.env")
	if err != nil {
		log.Fatalf("err loading: %v", err)
	}

	maxAttempts, err := strconv.Atoi(os.Getenv("MAX_ATTEMPTS"))
	if err != nil {
		log.Fatalf("err converting MAX_ATTEMPTS: %v", err)
	}

	return &models.Config{
		User:        os.Getenv("POSTGRES_USER"),
		Password:    os.Getenv("POSTGRES_PASSWORD"),
		Host:        os.Getenv("DB_HOST"),
		Port:        os.Getenv("DB_PORT"),
		DBName:      os.Getenv("POSTGRES_DB"),
		MaxAttempts: maxAttempts,
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
