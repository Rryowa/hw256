//go:build integration

package tests

import (
	"context"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"hash/maphash"
	"homework/internal/models"
	"homework/tests/db"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

type IntTestSuite struct {
	suite.Suite

	tr            db.TestRepository
	expectedOrder models.Order
}

func TestIntTestSuite(t *testing.T) {
	suite.Run(t, new(IntTestSuite))
}

func (s *IntTestSuite) SetupSuite() {
	cfg := NewTestConfig()
	ctx := context.Background()

	s.tr = db.NewTestRepository(ctx, cfg)

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
	s.expectedOrder = models.Order{
		ID:           dto.ID,
		UserID:       dto.UserID,
		StorageUntil: storageUntilDate,
		Issued:       false,
		OrderPrice:   models.Price(orderPriceFloat),
		Weight:       models.Weight(weightFloat),
		Returned:     false,
	}
}

func (s *IntTestSuite) TearDownSuite() {
	ctx := context.Background()

	err := s.tr.Repo.Truncate(ctx, "orders")

	require.NoError(s.T(), err)
}

func (s *IntTestSuite) setupOrder() (context.Context, models.Order) {
	ctx := context.Background()
	r := rand.New(rand.NewSource(int64(new(maphash.Hash).Sum64())))
	order := s.expectedOrder
	order.ID = strconv.FormatInt(r.Int63(), 10)
	return ctx, order
}

func (s *IntTestSuite) TestInsert() {
	s.T().Run("TestInsert", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		r := rand.New(rand.NewSource(int64(new(maphash.Hash).Sum64())))
		order := s.expectedOrder
		order.ID = strconv.FormatInt(r.Int63(), 10)

		id, err := s.tr.Repo.Insert(ctx, order)

		require.NoError(s.T(), err)
		require.Equal(s.T(), order.ID, id)
	})
}

func (s *IntTestSuite) TestUpdateOrder() {
	s.T().Run("TestUpd", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		r := rand.New(rand.NewSource(int64(new(maphash.Hash).Sum64())))
		order := s.expectedOrder
		order.ID = strconv.FormatInt(r.Int63(), 10)
		_, err := s.tr.Repo.Insert(ctx, order)
		order.Returned = true

		returned, err := s.tr.Repo.Update(ctx, order)

		require.NoError(t, err)
		require.Equal(t, order.Returned, returned)
	})
}

func (s *IntTestSuite) TestDelete() {
	s.T().Run("TestDelete", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		r := rand.New(rand.NewSource(int64(new(maphash.Hash).Sum64())))
		order := s.expectedOrder
		order.ID = strconv.FormatInt(r.Int63(), 10)
		_, err := s.tr.Repo.Insert(ctx, order)

		id, err := s.tr.Repo.Delete(ctx, order.ID)

		require.NoError(t, err)
		require.Equal(t, order.ID, id)
	})
}

func (s *IntTestSuite) TestGet() {
	s.T().Run("TestGet", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		r := rand.New(rand.NewSource(int64(new(maphash.Hash).Sum64())))
		order := s.expectedOrder
		order.ID = strconv.FormatInt(r.Int63(), 10)
		_, err := s.tr.Repo.Insert(ctx, order)

		order2, err := s.tr.Repo.Get(ctx, order.ID)

		require.NoError(t, err)
		require.Equal(t, order.ID, order2.ID)
	})
}

func (s *IntTestSuite) TestGetOrders() {
	s.T().Run("TestGetOrders", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		r := rand.New(rand.NewSource(int64(new(maphash.Hash).Sum64())))
		order := s.expectedOrder
		order.ID = strconv.FormatInt(r.Int63(), 10)
		_, err := s.tr.Repo.Insert(ctx, order)

		orders, err := s.tr.Repo.GetOrders(ctx, order.UserID, 0, 10)

		require.NoError(t, err)
		require.Equal(t, order.Issued, orders[0].Issued)
	})
}

func (s *IntTestSuite) TestGetReturns() {
	s.T().Run("TestGetReturns", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		r := rand.New(rand.NewSource(int64(new(maphash.Hash).Sum64())))
		order := s.expectedOrder
		order.ID = strconv.FormatInt(r.Int63(), 10)
		order.Returned = true
		_, err := s.tr.Repo.Insert(ctx, order)

		orders, err := s.tr.Repo.GetReturns(ctx, 0, 10)

		require.NoError(t, err)
		require.Equal(t, order.Returned, orders[0].Returned)
	})
}