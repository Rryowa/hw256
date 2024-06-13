package main

import (
	"context"
	"fmt"
	"homework-1/internal/cli"
	"homework-1/internal/service"
	"homework-1/internal/storage"
	"homework-1/internal/util"
	"log"
)

const (
	fileName = "orders.json"
)

func main() {
	orderStorage := storage.NewOrderStorage()
	repository := storage.NewSQLRepository(context.Background(), util.NewConfig())
	orderService := service.NewOrderService(orderStorage)
	validationService := service.NewOrderValidator(orderStorage, orderService, repository)

	commands := cli.NewCLI(validationService)
	if err := commands.Run(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Bye!")
}
