package orders_api

import (
	"context"
	"errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"homework/internal/models"
	"homework/internal/service"
	"homework/internal/util"
	"homework/mocks"
	proto "homework/pkg/api/proto/orders/v1/orders/v1"
	"log"
	"net"
	"strconv"
	"testing"
	"time"
)

const buffer = 1024 * 1024

func server(t *testing.T) (proto.OrderServiceClient, *mocks.MockStorage, *mocks.MockPackageService, *mocks.MockHasher, func()) {
	lis := bufconn.Listen(buffer)

	repository := mocks.NewMockStorage(t)
	pkgService := mocks.NewMockPackageService(t)
	hashService := mocks.NewMockHasher(t)
	orderService := service.NewOrderService(repository, pkgService, hashService)
	orderGrpcService := &OrderGrpcServer{
		OrderService: orderService,
	}
	baseServer := grpc.NewServer()
	proto.RegisterOrderServiceServer(baseServer, orderGrpcService)
	go func() {
		if err := baseServer.Serve(lis); err != nil {
			log.Printf("error serving server: %v", err)
		}
	}()

	conn, err := grpc.DialContext(
		context.Background(),
		"bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("error connecting to server: %v", err)
	}

	closer := func() {
		err := lis.Close()
		if err != nil {
			log.Printf("error closing listener: %v", err)
		}
		baseServer.Stop()
		conn.Close()
	}

	client := proto.NewOrderServiceClient(conn)

	return client, repository, pkgService, hashService, closer
}

func TestAcceptOrder(t *testing.T) {
	ctx := context.Background()
	client, repository, packageService, hashService, closer := server(t)
	defer closer()
	dto := models.Dto{
		ID:           "1",
		UserID:       "2",
		StorageUntil: "2077-07-07",
		OrderPrice:   "100",
		Weight:       "20",
		PackageType:  "box",
	}
	storageUntil, _ := time.Parse(time.DateOnly, dto.StorageUntil)
	orderPriceFloat, _ := strconv.ParseFloat(dto.OrderPrice, 64)
	weightFloat, _ := strconv.ParseFloat(dto.Weight, 64)
	expected := models.Order{
		ID:           dto.ID,
		UserID:       dto.UserID,
		StorageUntil: storageUntil,
		OrderPrice:   models.Price(orderPriceFloat),
		Weight:       models.Weight(weightFloat),
	}
	repository.EXPECT().Get(mock.Anything, expected.ID).Return(models.Order{}, util.ErrOrderNotFound)
	packageService.EXPECT().ValidatePackage(expected.Weight, models.PackageType(dto.PackageType)).Return(nil)
	packageService.EXPECT().ApplyPackage(&expected, models.PackageType(dto.PackageType))
	hashService.EXPECT().GenerateHash().Return(expected.Hash)
	repository.EXPECT().Insert(mock.Anything, expected).Return(expected.ID, nil)
	_, err := client.AcceptOrder(ctx, &proto.AcceptOrderRequest{
		Id:     dto.ID,
		UserId: dto.UserID,
		Date:   dto.StorageUntil,
		Price:  dto.OrderPrice,
		Weight: dto.Weight,
	})

	require.NoError(t, err)
}

func TestAcceptOrder_InvalidRequest(t *testing.T) {
	ctx := context.Background()
	client, _, _, _, closer := server(t)
	defer closer()

	_, err := client.AcceptOrder(ctx, &proto.AcceptOrderRequest{
		Id:     "",
		UserId: "2",
		Date:   "2077-07-07",
		Price:  "100",
		Weight: "20",
	})

	require.EqualError(t, err, errors.New("rpc error: code = InvalidArgument desc = invalid AcceptOrderRequest.Id: value does not match regex pattern \"^[0-9]+$\"").Error())
}
func TestAcceptOrder_InvalidRequest_Non_numeric_argument(t *testing.T) {
	ctx := context.Background()
	client, _, _, _, closer := server(t)
	defer closer()

	_, err := client.AcceptOrder(ctx, &proto.AcceptOrderRequest{
		Id:     "abc",
		UserId: "2",
		Date:   "2077-07-07",
		Price:  "100",
		Weight: "20",
	})

	require.EqualError(t, err, errors.New("rpc error: code = InvalidArgument desc = invalid AcceptOrderRequest.Id: value does not match regex pattern \"^[0-9]+$\"").Error())
}

func TestIssueOrders(t *testing.T) {
	ctx := context.Background()
	client, repository, _, _, closer := server(t)
	defer closer()
	dto := models.Dto{
		ID:           "1",
		UserID:       "2",
		StorageUntil: "2077-07-07",
	}
	storageUntil, _ := time.Parse(time.DateOnly, dto.StorageUntil)
	order := models.Order{
		ID:           dto.ID,
		UserID:       dto.UserID,
		StorageUntil: storageUntil,
		Issued:       false,
		Returned:     false,
	}
	expected := models.Order{
		ID:           dto.ID,
		UserID:       dto.UserID,
		StorageUntil: storageUntil,
		Issued:       true,
		Returned:     false,
	}

	repository.EXPECT().Get(mock.Anything, "1").Return(order, nil)
	repository.EXPECT().IssueUpdate(mock.Anything, []models.Order{expected}).
		Return([]bool{expected.Issued}, nil)

	_, err := client.IssueOrders(ctx, &proto.IssueOrdersRequest{
		Ids: "1",
	})

	require.NoError(t, err)
}

func TestIssueOrders_InvalidRequest(t *testing.T) {
	ctx := context.Background()
	client, _, _, _, closer := server(t)
	defer closer()

	_, err := client.IssueOrders(ctx, &proto.IssueOrdersRequest{
		Ids: "",
	})

	require.EqualError(t, err, errors.New("rpc error: code = InvalidArgument desc = invalid IssueOrdersRequest.Ids: value length must be at least 1 runes").Error())
}

func TestReturn(t *testing.T) {
	client, repository, _, _, closer := server(t)
	defer closer()
	ctx := context.Background()
	storageUntil, _ := time.Parse(time.DateOnly, "2077-07-07")
	issuedAt, _ := time.Parse(time.DateOnly, time.Now().Format(time.DateOnly))
	order := models.Order{
		ID:           "1",
		UserID:       "1",
		StorageUntil: storageUntil,
		Issued:       true,
		IssuedAt:     issuedAt,
		Returned:     false,
	}
	expected := models.Order{
		ID:           order.ID,
		UserID:       order.UserID,
		StorageUntil: storageUntil,
		Issued:       true,
		IssuedAt:     issuedAt,
		Returned:     true,
	}

	repository.EXPECT().Get(mock.Anything, "1").Return(order, nil)
	repository.EXPECT().Update(mock.Anything, expected).Return(expected.Returned, nil)

	_, err := client.AcceptReturn(ctx, &proto.AcceptReturnRequest{
		Id:     order.ID,
		UserId: order.UserID,
	})

	require.NoError(t, err)
}

func TestReturn_InvalidRequest(t *testing.T) {
	client, _, _, _, closer := server(t)
	defer closer()
	ctx := context.Background()

	_, err := client.AcceptReturn(ctx, &proto.AcceptReturnRequest{
		Id:     "",
		UserId: "1",
	})

	require.EqualError(t, err, errors.New("rpc error: code = InvalidArgument desc = invalid AcceptReturnRequest.Id: value does not match regex pattern \"^[0-9]+$\"").Error())
}

func TestReturnToCourier(t *testing.T) {
	client, repository, _, _, closer := server(t)
	defer closer()
	ctx := context.Background()
	storageUntil, _ := time.Parse(time.DateOnly, "2002-07-07")
	expected := models.Order{
		ID:           "1",
		UserID:       "2",
		StorageUntil: storageUntil,
		Issued:       false,
	}
	repository.EXPECT().Get(mock.Anything, expected.ID).Return(expected, nil)
	repository.EXPECT().Delete(mock.Anything, expected.ID).Return(expected.ID, nil)

	_, err := client.ReturnOrderToCourier(ctx, &proto.ReturnOrderToCourierRequest{
		Id: expected.ID,
	})

	require.NoError(t, err)
}

func TestReturnToCourier_InvalidRequest(t *testing.T) {
	client, _, _, _, closer := server(t)
	defer closer()
	ctx := context.Background()

	_, err := client.ReturnOrderToCourier(ctx, &proto.ReturnOrderToCourierRequest{
		Id: "",
	})

	require.EqualError(t, err, errors.New("rpc error: code = InvalidArgument desc = invalid ReturnOrderToCourierRequest.Id: value does not match regex pattern \"^[0-9]+$\"").Error())
}