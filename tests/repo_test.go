//go:build intgr_test

package tests

import (
	"context"
	"github.com/stretchr/testify/require"
	"homework/internal/models"
	"homework/internal/util"
	"homework/tests/db"
	"strconv"
	"testing"
	"time"
)

func TestInsertOrder(t *testing.T) {
	var (
		ctx = context.Background()
	)
	repo := db.NewTestRepository(ctx, util.NewTestConfig())
	repo.SetUp(t)
	defer repo.TearDown(t)
	dto := models.Dto{
		ID:           "1",
		UserID:       "1",
		StorageUntil: "2077-01-01",
		OrderPrice:   "999.99",
		Weight:       "10",
		PackageType:  "box",
	}
	storageUntilDate, _ := time.Parse(time.DateOnly, dto.StorageUntil)
	orderPriceFloat, _ := strconv.ParseFloat(dto.OrderPrice, 64)
	weightFloat, _ := strconv.ParseFloat(dto.Weight, 64)
	expectedOrder := models.Order{
		ID:           dto.ID,
		UserID:       dto.UserID,
		StorageUntil: storageUntilDate,
		OrderPrice:   models.Price(orderPriceFloat),
		Weight:       models.Weight(weightFloat),
	}
	id, err := repo.Storage.Insert(expectedOrder)
	require.NoError(t, err)
	require.Equal(t, expectedOrder.ID, id)
}
