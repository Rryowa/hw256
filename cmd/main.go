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
	orderStorage := storage.NewOrderStorage()
	validationService := service.NewOrderValidator(orderStorage)
	fileService := service.NewFileService(fileName)
	orderService := service.NewOrderService(orderStorage, fileService)

	commands := cli.NewCLI(validationService, orderService, fileService)
	if err := commands.Run(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Bye!")
}
