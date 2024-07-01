package service

import (
	"github.com/stretchr/testify/require"
	"homework/internal/models"
	"homework/mocks"
	"testing"
)

func Test_Exists(t *testing.T) {
	t.Run("Test_Exists", func(t *testing.T) {
		t.Parallel()
		repo := mocks.NewMockStorage(t)
		packageService := mocks.NewMockPackageService(t)
		orderService := NewOrderService(repo, packageService)
		repo.EXPECT().Get("1").Return(models.Order{}, nil)

		_, exists := orderService.Exists("1")
		require.Equal(t, true, exists)
	})
}
