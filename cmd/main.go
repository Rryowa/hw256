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
	if err := repository.ApplyMigrations("up"); err != nil {
		log.Fatalf("Migration failed: %v\n", err)
	}
	orderService := service.NewOrderService()
	validationService := service.NewOrderValidator(repository, orderService)

	commands := cli.NewCLI(validationService)
	if err := commands.Run(); err != nil {
		log.Fatal(err)
	}

	if err := repository.ApplyMigrations("down"); err != nil {
		log.Fatalf("Migration failed: %v\n", err)
	}
	fmt.Println("Bye!")
}
