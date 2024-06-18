package main

import (
	"context"
	"fmt"
	"homework-1/internal/cli"
	"homework-1/internal/service"
	"homework-1/internal/storage/db"
	"homework-1/internal/util"
	"log"
)

func main() {
	repository := db.NewSQLRepository(context.Background(), util.NewConfig())
	orderService := service.NewOrderService(repository)

	commands := cli.NewCLI(orderService)
	if err := commands.Run(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Bye!")
}
