package service

import (
	"context"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"homework/internal/models"
	"homework/internal/util"
	"homework/mocks"
	"strconv"
	"testing"
	"time"
)

func TestAccept_HappyPath(t *testing.T) {
	ctx := context.Background()
	dto := models.Dto{
		ID:           "1",
		UserID:       "test",
		StorageUntil: "2077-07-07",
		OrderPrice:   "100",
		Weight:       "20",
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
		Hash:         "hash",
	}
	mockRepo := mocks.NewMockStorage(t)
	mockRepo.EXPECT().Get(mock.Anything, dto.ID).Return(models.Order{}, util.ErrOrderNotFound)
	mockPackageSrvc := mocks.NewMockPackageService(t)
	mockPackageSrvc.EXPECT().ValidatePackage(expected.Weight, models.PackageType("box")).Return(nil)
	mockPackageSrvc.EXPECT().ApplyPackage(&expected, models.PackageType("box"))
	//I mocked Hasher because creating an alias for GenerateHash() is not concurrency safe.
	mockHasher := mocks.NewMockHasher(t)
	mockHasher.EXPECT().GenerateHash().Return(expected.Hash)
	orderSvc := NewOrderService(mockRepo, mockPackageSrvc, mockHasher)

	mockRepo.EXPECT().Insert(mock.Anything, expected).Return(expected.ID, nil)

	err := orderSvc.Accept(ctx, dto, "box")

	require.NoError(t, err)
}

func TestAccept_ErrOrderExists(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	dto := models.Dto{
		ID:           "1",
		UserID:       "2",
		StorageUntil: "2077-07-07",
		OrderPrice:   "100",
		Weight:       "5",
	}
	mockRepo := mocks.NewMockStorage(t)
	mockRepo.EXPECT().Get(mock.Anything, dto.ID).Return(models.Order{}, nil)
	orderSvc := NewOrderService(mockRepo, mocks.NewMockPackageService(t), mocks.NewMockHasher(t))

	err := orderSvc.Accept(ctx, dto, "")

	require.Equal(t, util.ErrOrderExists, err)
}

func TestAccept_ErrDateInvalid(t *testing.T) {
	t.Run("InvalidDateFormat", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		dto := models.Dto{
			ID:           "1",
			UserID:       "2",
			StorageUntil: "Some date",
			OrderPrice:   "100",
			Weight:       "5",
		}
		mockRepo := mocks.NewMockStorage(t)
		mockRepo.EXPECT().Get(mock.Anything, dto.ID).Return(models.Order{}, util.ErrOrderNotFound)
		orderSvc := NewOrderService(mockRepo, mocks.NewMockPackageService(t), mocks.NewMockHasher(t))

		err := orderSvc.Accept(ctx, dto, "")

		require.Equal(t, util.ErrParsingDate, err)
	})
	t.Run("DateBeforeNow", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		dto := models.Dto{
			ID:           "1",
			UserID:       "2",
			StorageUntil: time.Now().AddDate(-1, 0, 0).Format(time.DateOnly),
			OrderPrice:   "100",
			Weight:       "5",
		}
		mockRepo := mocks.NewMockStorage(t)
		mockRepo.EXPECT().Get(mock.Anything, dto.ID).Return(models.Order{}, util.ErrOrderNotFound)
		orderSvc := NewOrderService(mockRepo, mocks.NewMockPackageService(t), mocks.NewMockHasher(t))

		err := orderSvc.Accept(ctx, dto, "")

		require.Equal(t, util.ErrDateInvalid, err)
	})
}

func TestAccept_ErrOrderPriceInvalid(t *testing.T) {
	t.Run("InvalidPrice", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		dto := models.Dto{
			ID:           "1",
			UserID:       "2",
			StorageUntil: "2077-07-07",
			OrderPrice:   "some price",
			Weight:       "5",
		}
		mockRepo := mocks.NewMockStorage(t)
		mockRepo.EXPECT().Get(mock.Anything, dto.ID).Return(models.Order{}, util.ErrOrderNotFound)
		orderSvc := NewOrderService(mockRepo, mocks.NewMockPackageService(t), mocks.NewMockHasher(t))

		err := orderSvc.Accept(ctx, dto, "")
		require.Equal(t, util.ErrOrderPriceInvalid, err)
	})
	t.Run("NegativePrice", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		dto := models.Dto{
			ID:           "1",
			UserID:       "2",
			StorageUntil: "2077-07-07",
			OrderPrice:   "-100",
			Weight:       "5",
		}
		mockRepo := mocks.NewMockStorage(t)
		mockRepo.EXPECT().Get(mock.Anything, dto.ID).Return(models.Order{}, util.ErrOrderNotFound)
		orderSvc := NewOrderService(mockRepo, mocks.NewMockPackageService(t), mocks.NewMockHasher(t))

		err := orderSvc.Accept(ctx, dto, "")
		require.Equal(t, util.ErrOrderPriceInvalid, err)
	})
}

func TestAccept_ErrWeightInvalid(t *testing.T) {
	t.Run("InvalidWeight", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		dto := models.Dto{
			ID:           "1",
			UserID:       "2",
			StorageUntil: "2077-07-07",
			OrderPrice:   "100",
			Weight:       "some weight",
		}
		mockRepo := mocks.NewMockStorage(t)
		mockRepo.EXPECT().Get(mock.Anything, dto.ID).Return(models.Order{}, util.ErrOrderNotFound)
		orderSvc := NewOrderService(mockRepo, mocks.NewMockPackageService(t), mocks.NewMockHasher(t))

		err := orderSvc.Accept(ctx, dto, "")
		require.Equal(t, util.ErrWeightInvalid, err)
	})
	t.Run("NegativeWeight", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		dto := models.Dto{
			ID:           "1",
			UserID:       "2",
			StorageUntil: "2077-07-07",
			OrderPrice:   "100",
			Weight:       "-5",
		}
		mockRepo := mocks.NewMockStorage(t)
		mockRepo.EXPECT().Get(mock.Anything, dto.ID).Return(models.Order{}, util.ErrOrderNotFound)
		orderSvc := NewOrderService(mockRepo, mocks.NewMockPackageService(t), mocks.NewMockHasher(t))

		err := orderSvc.Accept(ctx, dto, "")
		require.Equal(t, util.ErrWeightInvalid, err)
	})
}