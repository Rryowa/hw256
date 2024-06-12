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
	validationService := service.NewOrderValidator(orderStorage)
	orderService := service.NewOrderService(orderStorage)

	commands := cli.NewCLI(validationService, orderService, fileService)
	if err := commands.Run(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Bye!")
}
