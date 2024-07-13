package main

import (
	"context"
	"flag"
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

const (
	host     = "localhost"
	grpcPort = 50051
	httpPort = ":44444"
)

var (
	grpcServerEndpoint = flag.String("grpc-server-endpoint", host+httpPort, "gRPC server endpoint")
)

func main() {
	flag.Parse()
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

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
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
	err = proto.RegisterOrderServiceHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint, opts)
	if err != nil {
		log.Fatalf("failed to RegisterOrderServiceHandlerFromEndpoint: %v", err)
		return
	}

	go func() {
		gwServer := &http.Server{
			Addr:    httpPort,
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