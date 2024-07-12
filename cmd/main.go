package main

import (
	"context"
	"homework/internal/service"
	"homework/internal/storage/db"
	"homework/internal/util"
	"homework/internal/view"
	"homework/pkg/hash"
	"log"
)

func main() {
	ctx := context.Background()
	repository := db.NewSQLRepository(ctx, util.NewConfig())
	orderService := service.NewOrderService(repository, service.NewPackageService(), &hash.HashGenerator{})
	kafkaService := service.NewKafkaService(util.NewKafkaConfig(), repository)

	commands := view.NewCLI(orderService, kafkaService)
	if err := commands.Run(ctx); err != nil {
		log.Fatal(err)
	}
}