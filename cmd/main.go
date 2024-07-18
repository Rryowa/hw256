package main

import (
	"context"
	log "github.com/sirupsen/logrus"
	"homework/internal/service"
	"homework/internal/storage/db"
	"homework/internal/util"
	"homework/internal/view"
	"homework/pkg/cache"
	"homework/pkg/hash"
	"homework/pkg/kafka"
)

// TODO: prometheus collects data - grafana displays
func main() {
	ctx := context.Background()
	zapLogger := util.NewZapLogger()
	repository := db.NewSQLRepository(ctx, util.NewDbConfig(), zapLogger)
	cacheService := cache.NewCache(util.NewCacheConfig())
	orderService := service.NewOrderService(repository, cacheService, service.NewPackageService(), &hash.HashGenerator{})
	loggerService := kafka.NewLoggerService(util.NewKafkaConfig(), repository, zapLogger)
	commands := view.NewCLI(orderService, loggerService, zapLogger)
	if err := commands.Run(ctx); err != nil {
		log.Fatal(err)
	}
}