package main

import (
	"fmt"
	"homework-1/internal/cli"
	"homework-1/internal/service"
	"homework-1/internal/storage"
	"log"
)

const (
	fileName = "orders.json"
)

func main() {
	fileService := service.NewFileService(fileName)
	orderStorage := storage.NewOrderStorage()
	orderService := service.NewOrderService(orderStorage)
	validationService := service.NewOrderValidator(orderStorage, orderService, fileService)

	commands := cli.NewCLI(validationService)
	if err := commands.Run(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Bye!")
}
