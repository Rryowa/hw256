package main

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"homework/internal/api"
	"homework/internal/middleware"
	"homework/internal/service"
	"homework/internal/storage/db"
	"homework/internal/util"
	proto "homework/pkg/api/proto/orders/v1/orders/v1"
	"homework/pkg/hash"
	"log"
	"net"
	"net/http"
	"sync"
)

func main() {
	cfg := util.NewGrpcConfig()
	ctx := context.Background()
	var wg sync.WaitGroup

	repository := db.NewSQLRepository(ctx, util.NewConfig())
	orderService := service.NewOrderService(repository, service.NewPackageService(), &hash.HashGenerator{})
	orderGrpcService := &orders_api.OrderGrpcServer{
		OrderService: orderService,
	}

	logger := service.NewLoggerService(util.NewKafkaConfig(), repository)
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