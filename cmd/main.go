package main

import (
	"context"
	"fmt"
	"homework/internal/service"
	"homework/internal/storage/db"
	"homework/internal/util"
	"homework/internal/view"
	"log"
)

func main() {
	repository := db.NewSQLRepository(context.Background(), util.NewConfig())

	packageService := service.NewPackageService()
	orderService := service.NewOrderService(repository, packageService)

	commands := view.NewCLI(orderService)
	if err := commands.Run(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Bye!")
}
