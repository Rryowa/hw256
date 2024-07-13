package main

import (
	"context"
	"flag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding/gzip"
	proto "homework/pkg/api/proto/orders/v1/orders/v1"
	"log"
	"time"
)

const (
	id     string = "999"
	userId string = "1"
)

var addr = flag.String("addr", "localhost:50051", "the address to connect to")

func main() {
	flag.Parse()
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := grpc.NewClient(
		*addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.UseCompressor(gzip.Name)),
	)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	client := proto.NewOrderServiceClient(conn)
	acceptOrder(ctx, client)
	//issueOrders(ctx, client)
	//acceptReturn(ctx, client)
	//returnOrderToCourier(ctx, client)
	//listReturns(ctx, client)
	//listOrders(ctx, client)

	log.Println("Client done")
}

func acceptOrder(ctx context.Context, client proto.OrderServiceClient) {
	_, err := client.AcceptOrder(ctx, &proto.AcceptOrderRequest{
		Id:     "5",
		UserId: userId,
		Date:   "2077-07-07",
		Price:  "100",
		Weight: "20",
	})
	if err != nil {
		log.Fatal(err)
	}
}

func issueOrders(ctx context.Context, client proto.OrderServiceClient) {
	_, err := client.IssueOrders(ctx, &proto.IssueOrdersRequest{
		Ids: id,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func acceptReturn(ctx context.Context, client proto.OrderServiceClient) {
	_, err := client.AcceptReturn(ctx, &proto.AcceptReturnRequest{
		Id:     id,
		UserId: userId,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func returnOrderToCourier(ctx context.Context, client proto.OrderServiceClient) {
	_, err := client.ReturnOrderToCourier(ctx, &proto.ReturnOrderToCourierRequest{
		Id: id,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func listReturns(ctx context.Context, client proto.OrderServiceClient) {
	_, err := client.ListReturns(ctx, &proto.ListReturnsRequest{
		Offset: "0",
		Limit:  "10",
	})
	if err != nil {
		log.Fatal(err)
	}
}

func listOrders(ctx context.Context, client proto.OrderServiceClient) {
	_, err := client.ListOrders(ctx, &proto.ListOrdersRequest{
		UserId: userId,
		Offset: "0",
		Limit:  "10",
	})
	if err != nil {
		log.Fatal(err)
	}
}