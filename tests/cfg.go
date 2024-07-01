package tests

import (
	"github.com/joho/godotenv"
	"homework/internal/models"
	"log"
	"os"
)

func NewTestConfig() *models.Config {
	//for single stage build using @go test ./tests -tags=integration
	err := godotenv.Load("test.env")
	if err != nil {
		log.Fatalf("err loading: %v", err)
	}
	return &models.Config{
		User:     os.Getenv("TEST_USER"),
		Password: os.Getenv("TEST_PASSWORD"),
		Host:     os.Getenv("TEST_HOST"),
		Port:     os.Getenv("TEST_PORT"),
		DBName:   os.Getenv("TEST_DB"),
	}
}
