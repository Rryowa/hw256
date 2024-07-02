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

func TestIssue_HappyPath(t *testing.T) {
	ctx := context.Background()
	storageUntil, _ := time.Parse(time.DateOnly, "2077-07-07")
	issuedAt := storageUntil.AddDate(-1, 0, 0)
	orders := []models.Order{
		{
			ID:           "1",
			UserID:       "1",
			StorageUntil: storageUntil,
			Issued:       false,
			Returned:     false,
		},
		{
			ID:           "2",
			UserID:       "1",
			StorageUntil: storageUntil,
			Issued:       false,
			Returned:     false,
		},
	}
	expected := []models.Order{
		{
			ID:           "1",
			UserID:       "1",
			StorageUntil: storageUntil,
			Issued:       true,
			IssuedAt:     issuedAt,
		},
		{
			ID:           "2",
			UserID:       "1",
			StorageUntil: storageUntil,
			Issued:       true,
			IssuedAt:     issuedAt,
		},
	}
	mockRepo := mocks.NewMockStorage(t)
	mockTimer := mocks.NewMockTimer(t)
	orderSvc := NewOrderService(mockRepo, mocks.NewMockPackageService(t), mocks.NewMockHasher(t), mockTimer)

	mockRepo.EXPECT().Get(mock.Anything, "1").Return(orders[0], nil)
	mockRepo.EXPECT().Get(mock.Anything, "2").Return(orders[1], nil)
	mockTimer.EXPECT().TimeNow().Return(issuedAt)
	mockRepo.EXPECT().IssueUpdate(mock.Anything, expected).
		Return([]bool{expected[0].Issued, expected[1].Issued}, nil)

	err := orderSvc.Issue(ctx, "1,2")

	require.NoError(t, err)
}

func TestAccept_ErrOrderNotFound(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	mockRepo := mocks.NewMockStorage(t)
	mockRepo.EXPECT().Get(mock.Anything, "1").Return(models.Order{}, util.ErrOrderNotFound)

	orderSvc := NewOrderService(mockRepo, mocks.NewMockPackageService(t), mocks.NewMockHasher(t), mocks.NewMockTimer(t))

	err := orderSvc.Issue(ctx, "1")

	require.Equal(t, util.ErrOrderNotFound, err)
}

func TestAccept_ErrOrderIssued(t *testing.T) {
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

	err := orderSvc.Issue(ctx, expected.ID)

	require.Equal(t, util.ErrOrderIssued, err)
}

func TestAccept_ErrOrderExpired(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	storageUntil, _ := time.Parse(time.DateOnly, "2002-02-02")
	expected := models.Order{
		ID:           "1",
		UserID:       "2",
		StorageUntil: storageUntil,
		Issued:       false,
	}
	mockRepo := mocks.NewMockStorage(t)
	mockRepo.EXPECT().Get(mock.Anything, expected.ID).Return(expected, nil)

	orderSvc := NewOrderService(mockRepo, mocks.NewMockPackageService(t), mocks.NewMockHasher(t), mocks.NewMockTimer(t))

	err := orderSvc.Issue(ctx, expected.ID)

	require.Equal(t, util.ErrOrderExpired, err)
}

func TestIssue_ErrOrdersUserDiffers(t *testing.T) {
	ctx := context.Background()
	storageUntil, _ := time.Parse(time.DateOnly, "2077-07-07")
	orders := []models.Order{
		{
			ID:           "1",
			UserID:       "1",
			StorageUntil: storageUntil,
			Issued:       false,
			Returned:     false,
		},
		{
			ID:           "2",
			UserID:       "2",
			StorageUntil: storageUntil,
			Issued:       false,
			Returned:     false,
		},
	}
	mockRepo := mocks.NewMockStorage(t)
	mockRepo.EXPECT().Get(mock.Anything, "1").Return(orders[0], nil)
	mockRepo.EXPECT().Get(mock.Anything, "2").Return(orders[1], nil)
	orderSvc := NewOrderService(mockRepo, mocks.NewMockPackageService(t), mocks.NewMockHasher(t), mocks.NewMockTimer(t))

	err := orderSvc.Issue(ctx, "1,2")

	require.Equal(t, util.ErrOrdersUserDiffers, err)
}
