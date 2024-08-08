package orders_api

import (
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"homework/internal/models"
	"homework/internal/service"
	proto "homework/pkg/api/proto/orders/v1/orders/v1"
)

type OrderGrpcServer struct {
	OrderService service.OrderService
	proto.UnimplementedOrderServiceServer
}

func (ogs *OrderGrpcServer) AcceptOrder(ctx context.Context, req *proto.AcceptOrderRequest) (*proto.OrderResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	dto, err := requestToDto(req)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = ogs.OrderService.Accept(ctx, dto, "box")
	if err != nil {
		return nil, err
	}

	return &proto.OrderResponse{Message: "success"}, nil
}

func (ogs *OrderGrpcServer) IssueOrders(ctx context.Context, req *proto.IssueOrdersRequest) (*proto.OrderResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	err := ogs.OrderService.Issue(ctx, req.Ids)
	if err != nil {
		return nil, err
	}

	return &proto.OrderResponse{Message: "success"}, nil
}

func (ogs *OrderGrpcServer) AcceptReturn(ctx context.Context, req *proto.AcceptReturnRequest) (*proto.OrderResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	err := ogs.OrderService.Return(ctx, req.Id, req.UserId)
	if err != nil {
		return nil, err
	}

	return &proto.OrderResponse{Message: "success"}, nil
}

func (ogs *OrderGrpcServer) ReturnOrderToCourier(ctx context.Context, req *proto.ReturnOrderToCourierRequest) (*proto.OrderResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	err := ogs.OrderService.ReturnToCourier(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &proto.OrderResponse{Message: "success"}, nil
}

func (ogs *OrderGrpcServer) ListReturns(ctx context.Context, req *proto.ListReturnsRequest) (*proto.ListReturnsResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	orders, err := ogs.OrderService.ListReturns(ctx, req.Offset, req.Limit)
	if err != nil {
		return nil, err
	}

	var protoOrders []*proto.Order
	for _, order := range orders {
		protoOrders = append(protoOrders, convertModelOrderToProtoOrder(order))
	}

	return &proto.ListReturnsResponse{Orders: protoOrders}, nil
}

func (ogs *OrderGrpcServer) ListOrders(ctx context.Context, req *proto.ListOrdersRequest) (*proto.ListOrdersResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	orders, err := ogs.OrderService.ListOrders(ctx, req.UserId, req.Offset, req.Limit)
	if err != nil {
		return nil, err
	}

	var protoOrders []*proto.Order
	for _, order := range orders {
		protoOrders = append(protoOrders, convertModelOrderToProtoOrder(order))
	}

	return &proto.ListOrdersResponse{Orders: protoOrders}, nil
}

func requestToDto(req *proto.AcceptOrderRequest) (models.Dto, error) {
	if len(req.Id) == 0 || len(req.UserId) == 0 || len(req.Date) == 0 ||
		len(req.Price) == 0 || len(req.Weight) == 0 {
		return models.Dto{}, errors.New("invalid request")
	}

	return models.Dto{
		ID:           req.Id,
		UserID:       req.UserId,
		StorageUntil: req.Date,
		OrderPrice:   req.Price,
		Weight:       req.Weight,
	}, nil
}

func convertModelOrderToProtoOrder(order models.Order) *proto.Order {
	return &proto.Order{
		Id:           order.ID,
		UserId:       order.UserID,
		StorageUntil: timestamppb.New(order.StorageUntil),
		Issued:       order.Issued,
		IssuedAt:     timestamppb.New(order.IssuedAt),
		Returned:     order.Returned,
		OrderPrice:   float64(order.OrderPrice),
		Weight:       float64(order.Weight),
		PackageType:  string(order.PackageType),
		PackagePrice: float64(order.PackagePrice),
		Hash:         order.Hash,
	}
}