package main

import (
	"fmt"
	"homework-1/internal/cli"
	"homework-1/internal/file"
	"homework-1/internal/storage"
	"os"
)

const (
	ordersName = "orders.json"
)

func main() {
	fileService := file.NewFileService(ordersName)
	orderStorage := storage.NewOrderStorage(fileService)
	commands := cli.NewCLI(orderStorage)
	if err := commands.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Bye!")

}
