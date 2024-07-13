package middleware

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"homework/internal/service"
	"log"
)

type contextKey string

const loggerKey contextKey = "logger"

func Logging() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		logger, ok := ctx.Value(loggerKey).(service.LoggerService)
		if !ok {
			log.Printf("[interceptor.Logging] no logger found in context")
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			header := md.Get("x-my-header")
			log.Printf("[interceptor.Logging] header: %v", header)
		}

		raw, _ := protojson.Marshal((req).(proto.Message))
		log.Printf("[interceptor.Logging] start: %v, %v", info.FullMethod, string(raw))

		// Create a logging event at the start
		event, err := logger.CreateEvent(ctx, info.FullMethod+" - start")
		if err != nil {
			log.Printf("[interceptor.Logging] error creating start event: %v", err)
		}

		resp, err = handler(ctx, req)
		if err != nil {
			log.Printf("[interceptor.Logging] error: %v", err.Error())
			// Create a logging event on error
			if _, err := logger.CreateEvent(ctx, info.FullMethod+" - error: "+err.Error()); err != nil {
				log.Printf("[interceptor.Logging] error creating error event: %v", err)
			}
			return nil, err
		}

		// Create a logging event at the end
		if err := logger.ProcessEvent(ctx, event); err != nil {
			log.Printf("[interceptor.Logging] error processing event: %v", err)
		}

		log.Println("[interceptor.Logging] end")
		return resp, nil
	}
}

// AddLoggerToContext provides access to kafka service
func AddLoggerToContext(logger service.LoggerService) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		ctx = context.WithValue(ctx, loggerKey, logger)
		return handler(ctx, req)
	}
}