package service

import (
	"context"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"homework/internal/models"
	"homework/internal/util"
	"homework/mocks"
	"testing"
	"time"
)

func TestReturn_HappyPath(t *testing.T) {
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
		ID:           "1",
		UserID:       "1",
		StorageUntil: storageUntil,
		Issued:       true,
		IssuedAt:     issuedAt,
		Returned:     true,
	}
	mockRepo := mocks.NewMockStorage(t)
	mockTimer := mocks.NewMockTimer(t)
	orderSvc := NewOrderService(mockRepo, mocks.NewMockPackageService(t), mocks.NewMockHasher(t), mockTimer)
	mockRepo.EXPECT().Get(mock.Anything, "1").Return(order, nil)
	mockRepo.EXPECT().Update(mock.Anything, expected).
		Return(expected.Returned, nil)

	err := orderSvc.Return(ctx, "1", "1")

	require.NoError(t, err)
}

func TestReturn_ErrOrderNotFound(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	mockRepo := mocks.NewMockStorage(t)
	mockRepo.EXPECT().Get(mock.Anything, "1").Return(models.Order{}, util.ErrOrderNotFound)
	orderSvc := NewOrderService(mockRepo, mocks.NewMockPackageService(t), mocks.NewMockHasher(t), mocks.NewMockTimer(t))

	err := orderSvc.Return(ctx, "1", "1")

	require.Equal(t, util.ErrOrderNotFound, err)
}

func TestReturn_ErrOrderDoesNotBelong(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	storageUntil, _ := time.Parse(time.DateOnly, "2077-07-07")
	expected := models.Order{
		ID:           "1",
		UserID:       "2",
		StorageUntil: storageUntil,
		Issued:       true,
	}
	mockRepo := mocks.NewMockStorage(t)
	mockRepo.EXPECT().Get(mock.Anything, expected.ID).Return(expected, nil)
	orderSvc := NewOrderService(mockRepo, mocks.NewMockPackageService(t), mocks.NewMockHasher(t), mocks.NewMockTimer(t))

	err := orderSvc.Return(ctx, expected.ID, "10")

	require.Equal(t, util.ErrOrderDoesNotBelong, err)
}

func TestReturn_ErrOrderNotIssued(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	storageUntil, _ := time.Parse(time.DateOnly, "2077-07-07")
	expected := models.Order{
		ID:           "1",
		UserID:       "2",
		StorageUntil: storageUntil,
		Issued:       false,
	}
	mockRepo := mocks.NewMockStorage(t)
	mockRepo.EXPECT().Get(mock.Anything, expected.ID).Return(expected, nil)
	orderSvc := NewOrderService(mockRepo, mocks.NewMockPackageService(t), mocks.NewMockHasher(t), mocks.NewMockTimer(t))

	err := orderSvc.Return(ctx, expected.ID, expected.UserID)

	require.Equal(t, util.ErrOrderNotIssued, err)
}

func TestReturn_ErrReturnPeriodExpired(t *testing.T) {
	ctx := context.Background()
	storageUntil, _ := time.Parse(time.DateOnly, "2077-07-07")
	issuedAt, _ := time.Parse(time.DateOnly, time.Now().Format(time.DateOnly))
	expected := models.Order{
		ID:           "1",
		UserID:       "2",
		StorageUntil: storageUntil,
		Issued:       true,
		IssuedAt:     issuedAt.AddDate(-1, 0, 0),
	}
	mockRepo := mocks.NewMockStorage(t)
	mockRepo.EXPECT().Get(mock.Anything, expected.ID).Return(expected, nil)
	orderSvc := NewOrderService(mockRepo, mocks.NewMockPackageService(t), mocks.NewMockHasher(t), mocks.NewMockTimer(t))

	err := orderSvc.Return(ctx, "1", "2")

	require.Equal(t, util.ErrReturnPeriodExpired, err)
}

func TestReturnToCourier_HappyPath(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	storageUntil, _ := time.Parse(time.DateOnly, "2002-07-07")
	expected := models.Order{
		ID:           "1",
		UserID:       "2",
		StorageUntil: storageUntil,
		Issued:       false,
	}
	mockRepo := mocks.NewMockStorage(t)
	mockRepo.EXPECT().Get(mock.Anything, expected.ID).Return(expected, nil)
	mockRepo.EXPECT().Delete(mock.Anything, expected.ID).Return(expected.ID, nil)
	orderSvc := NewOrderService(mockRepo, mocks.NewMockPackageService(t), mocks.NewMockHasher(t), mocks.NewMockTimer(t))

	err := orderSvc.ReturnToCourier(ctx, expected.ID)

	require.NoError(t, err)
}

func TestReturnToCourier_ErrOrderNotFound(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	mockRepo := mocks.NewMockStorage(t)
	mockRepo.EXPECT().Get(mock.Anything, "1").Return(models.Order{}, util.ErrOrderNotFound)
	orderSvc := NewOrderService(mockRepo, mocks.NewMockPackageService(t), mocks.NewMockHasher(t), mocks.NewMockTimer(t))

	err := orderSvc.ReturnToCourier(ctx, "1")

	require.Equal(t, util.ErrOrderNotFound, err)
}

func TestReturnToCourier_ErrOrderIssued(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	storageUntil, _ := time.Parse(time.DateOnly, "2077-07-07")
	expected := models.Order{
		ID:           "1",
		UserID:       "2",
		StorageUntil: storageUntil,
		Issued:       true,
	}
	mockRepo := mocks.NewMockStorage(t)
	mockRepo.EXPECT().Get(mock.Anything, expected.ID).Return(expected, nil)
	orderSvc := NewOrderService(mockRepo, mocks.NewMockPackageService(t), mocks.NewMockHasher(t), mocks.NewMockTimer(t))

	err := orderSvc.ReturnToCourier(ctx, expected.ID)

	require.Equal(t, util.ErrOrderIssued, err)
}

func TestReturnToCourier_ErrOrderNotExpired(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	storageUntil, _ := time.Parse(time.DateOnly, "2077-07-07")
	expected := models.Order{
		ID:           "1",
		UserID:       "2",
		StorageUntil: storageUntil,
		Issued:       false,
	}
	mockRepo := mocks.NewMockStorage(t)
	mockRepo.EXPECT().Get(mock.Anything, expected.ID).Return(expected, nil)
	orderSvc := NewOrderService(mockRepo, mocks.NewMockPackageService(t), mocks.NewMockHasher(t), mocks.NewMockTimer(t))

	err := orderSvc.ReturnToCourier(ctx, expected.ID)

	require.Equal(t, util.ErrOrderNotExpired, err)
}
