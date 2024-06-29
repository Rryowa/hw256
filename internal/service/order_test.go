package service

import (
	"github.com/stretchr/testify/require"
	"homework/internal/models"
	"homework/mocks"
	"testing"
)

func (uts *UnitTestSuite) Test_Exists() {
	uts.T().Run("Test_Exists", func(t *testing.T) {
		t.Parallel()
		repo := mocks.NewMockStorage(t)
		orderService := NewOrderService("public", repo, uts.packageService)
		repo.EXPECT().Get("1", "public").Return(models.Order{}, nil)

		_, exists := orderService.Exists("1")
		require.Equal(t, true, exists)
	})
}
