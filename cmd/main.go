package main

import (
	"context"
	"fmt"
	"homework/internal/service"
	"homework/internal/storage/db"
	"homework/internal/util"
	"homework/internal/view"
	"homework/pkg/hash"
	"homework/pkg/timer"
	"log"
)

func main() {
	repository := db.NewSQLRepository(context.Background(), util.NewConfig())

	packageService := service.NewPackageService()
	hashGenerator := &hash.HashGenerator{}
	timeGenerator := &timer.TimeGenerator{}
	orderService := service.NewOrderService(repository, packageService, hashGenerator, timeGenerator)

	commands := view.NewCLI(orderService)
	if err := commands.Run(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Bye!")
}
