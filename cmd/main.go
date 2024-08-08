package main

import (
	"context"
	"homework/internal/metrics"
	"homework/internal/service"
	"homework/internal/storage/db"
	"homework/internal/util"
	"homework/internal/view"
	"homework/pkg/hash"
	"homework/pkg/kafka"
)

func main() {
	ctx := context.Background()
	zapLogger := util.NewZapLogger()

	repository := db.NewSQLRepository(ctx, util.NewDbConfig(), zapLogger)
	packageService := service.NewPackageService()
	loggerService := kafka.NewLoggerService(util.NewKafkaConfig(), repository, zapLogger)
	go metrics.Listen(ctx, util.NewMetricsConfig(), zapLogger)
	orderService := service.NewOrderService(repository, packageService, &hash.HashGenerator{})

	commands := view.NewCLI(orderService, loggerService, zapLogger)

	if err := commands.Run(ctx); err != nil {
		zapLogger.Fatalln(err)
	}
}