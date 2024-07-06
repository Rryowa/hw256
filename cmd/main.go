package main

import (
	"context"
	"fmt"
	"homework/internal/service"
	"homework/internal/storage/db"
	"homework/internal/util"
	"homework/internal/view"
	"homework/pkg/hash"
	"log"
)

func main() {
	ctx := context.Background()
	cfg := util.NewConfig()
	repository := db.NewSQLRepository(ctx, cfg)
	packageService := service.NewPackageService()
	hashGenerator := &hash.HashGenerator{}
	orderService := service.NewOrderService(repository, packageService, hashGenerator)
	newOutbox := service.NewOutbox(ctx, cfg)

	commands := view.NewCLI(orderService, newOutbox)
	if err := commands.Run(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Bye!")
}