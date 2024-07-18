package main

import (
	"context"
	"homework/internal/service"
	"homework/internal/storage/db"
	"homework/internal/util"
	"homework/internal/view"
	"homework/pkg/cache"
	"homework/pkg/hash"
	"log"
)

// TODO: logrus
// TODO: prometheus collects data - grafana displays
func main() {
	ctx := context.Background()
	repository := db.NewSQLRepository(ctx, util.NewDbConfig())
	cacheService := cache.NewCache(util.NewCacheConfig())
	orderService := service.NewOrderService(repository, cacheService, service.NewPackageService(), &hash.HashGenerator{})
	loggerService := service.NewLoggerService(util.NewKafkaConfig(), repository)

	commands := view.NewCLI(orderService, loggerService)
	if err := commands.Run(ctx); err != nil {
		log.Fatal(err)
	}
}