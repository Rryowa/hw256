package main

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"homework/internal/metrics"
	"homework/internal/service"
	"homework/internal/storage/db"
	"homework/internal/telemetry"
	"homework/internal/util"
	"homework/internal/view"
	"homework/pkg/cache"
	"homework/pkg/hash"
	"homework/pkg/kafka"
)

func main() {
	ctx := context.Background()
	zapLogger := util.NewZapLogger()

	repository := db.NewSQLRepository(ctx, util.NewDbConfig(), zapLogger)
	packageService := service.NewPackageService()
	loggerService := kafka.NewLoggerService(util.NewKafkaConfig(), repository, zapLogger)
	cacheService := cache.NewCache(util.NewCacheConfig())
	serverMetrics := metrics.NewServerMetrics(prometheus.NewRegistry())
	go metrics.Listen("localhost:9080")
	telemetry.MustSetup(ctx, "cli")

	orderService := service.NewOrderService(repository, packageService,
		&hash.HashGenerator{}, cacheService, serverMetrics)

	commands := view.NewCLI(orderService, loggerService, zapLogger)

	if err := commands.Run(ctx); err != nil {
		log.Fatal(err)
	}
}