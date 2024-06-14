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

func main() {
	repository := storage.NewSQLRepository(context.Background(), util.NewConfig())
	orderService := service.NewOrderService()
	validationService := service.NewOrderValidator(repository, orderService)

	commands := cli.NewCLI(validationService)
	if err := commands.Run(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Bye!")
}
