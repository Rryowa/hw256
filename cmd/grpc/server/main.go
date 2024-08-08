package main

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"homework/internal/api"
	"homework/internal/metrics"
	"homework/internal/middleware"
	"homework/internal/service"
	"homework/internal/storage/cache"
	"homework/internal/storage/db"
	"homework/internal/telemetry"
	"homework/internal/util"
	proto "homework/pkg/api/proto/orders/v1/orders/v1"
	"homework/pkg/hash"
	"homework/pkg/kafka"
	"net"
	"net/http"
	"sync"
)

func main() {
	cfg := util.NewGrpcConfig()
	ctx := context.Background()
	var wg sync.WaitGroup
	zapLogger := util.NewZapLogger()

	repository := db.NewSQLRepository(ctx, util.NewDbConfig(), zapLogger)
	cacheService := cache.NewCache(util.NewCacheConfig())
	packageService := service.NewPackageService()
	serverMetrics := metrics.NewServerMetrics(prometheus.NewRegistry())
	go metrics.Listen("localhost:9080")
	telemetry.MustSetup(ctx, "cli")

	orderService := service.NewOrderService(repository, packageService,
		&hash.HashGenerator{}, cacheService, serverMetrics)

	orderGrpcService := &orders_api.OrderGrpcServer{
		OrderService: orderService,
	}

	logger := kafka.NewLoggerService(util.NewKafkaConfig(), repository, zapLogger)
	loggerClose := logger.Start(ctx, &wg)
	defer loggerClose()
	go logger.DisplayKafkaEvents()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GrpcPort))
	if err != nil {
		log.Fatalf("failed to listen grpc: %v", err)
	}
	log.Println("GRPC listening on", lis.Addr())

	grpcServer := grpc.NewServer(
		grpc.RPCDecompressor(grpc.NewGZIPDecompressor()),
		grpc.ChainUnaryInterceptor(
			middleware.AddLoggerToContext(logger),
			middleware.Logging(),
		),
	)
	proto.RegisterOrderServiceServer(grpcServer, orderGrpcService)

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err = proto.RegisterOrderServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf("%s:%s", cfg.Host, cfg.HttpPort), opts)
	if err != nil {
		log.Fatalf("failed to RegisterOrderServiceHandlerFromEndpoint: %v", err)
		return
	}

	go func() {
		gwServer := &http.Server{
			Addr:    fmt.Sprintf(":%s", cfg.HttpPort),
			Handler: middleware.WithHTTPLoggingMiddleware(mux),
		}

		err := gwServer.ListenAndServe()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	wg.Wait()
	log.Println("Done")
}