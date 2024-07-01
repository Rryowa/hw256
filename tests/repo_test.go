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

// IntTestSuite Parallel integration tests using postgres schemas
type IntTestSuite struct {
	suite.Suite

	createQuery   string
	expectedOrder models.Order
	cfg           *models.Config
}

func TestIntTestSuite(t *testing.T) {
	suite.Run(t, new(IntTestSuite))
}

func (s *IntTestSuite) SetupSuite() {
	s.cfg = NewTestConfig()
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
		OrderPrice:   models.Price(orderPriceFloat),
		Weight:       models.Weight(weightFloat),
		Returned:     false,
	}

	s.createQuery = `CREATE TABLE orders (
		id VARCHAR(255) PRIMARY KEY,
		user_id VARCHAR(255) NOT NULL,
		storage_until TIMESTAMPTZ NOT NULL,
		issued BOOLEAN NOT NULL,
		issued_at TIMESTAMPTZ,
		returned BOOLEAN NOT NULL,
		order_price FLOAT NOT NULL,
		weight FLOAT NOT NULL,
		package_type VARCHAR(255) NOT NULL,
		package_price FLOAT NOT NULL,
		hash VARCHAR(255) NOT NULL
	);
	CREATE INDEX id_asc ON orders (id ASC);
	CREATE INDEX user_id_storage_asc ON orders (user_id, storage_until ASC);`
}

func (s *IntTestSuite) createRepoAndSchema(ctx context.Context) db.TestRepository {
	name := rand.New(rand.NewSource(int64(new(maphash.Hash).Sum64())))
	schemaName := "test" + strconv.FormatInt(name.Int63(), 10)

	tr := db.NewTestRepository(ctx, s.cfg, schemaName)

	_, err := tr.Repo.Pool.Exec(ctx, "CREATE SCHEMA "+schemaName)
	require.NoError(s.T(), err)

	//query := storage.AddSchemaPrefix(schemaName, s.createQuery)
	_, err = tr.Repo.Pool.Exec(ctx, s.createQuery)
	require.NoError(s.T(), err)

	return tr
}

func (s *IntTestSuite) TestInsertOrder() {
	s.T().Run("TestInsertOrder", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		tr := s.createRepoAndSchema(ctx)
		defer tr.DropSchema(ctx, t)

		id, err := tr.Repo.Insert(ctx, s.expectedOrder)

		require.NoError(t, err)
		require.Equal(t, s.expectedOrder.ID, id)
	})
}

func (s *IntTestSuite) TestUpdateOrder() {
	s.T().Run("TestUpd", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		tr := s.createRepoAndSchema(ctx)
		defer tr.DropSchema(ctx, t)
		order := s.expectedOrder
		order.Returned = true
		_, err := tr.Repo.Insert(ctx, order)
		require.NoError(t, err)

		returned, err := tr.Repo.Update(ctx, order)

		require.NoError(t, err)
		require.Equal(t, order.Returned, returned)
	})
}

func (s *IntTestSuite) TestDelete() {
	s.T().Run("TestDelete", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		tr := s.createRepoAndSchema(ctx)
		defer tr.DropSchema(ctx, t)
		order := s.expectedOrder
		_, err := tr.Repo.Insert(ctx, order)
		require.NoError(t, err)

		id, err := tr.Repo.Delete(ctx, order.ID)

		require.NoError(t, err)
		require.Equal(t, order.ID, id)
	})
}

func (s *IntTestSuite) TestGet() {
	s.T().Run("TestGet", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		tr := s.createRepoAndSchema(ctx)
		defer tr.DropSchema(ctx, t)
		order := s.expectedOrder
		_, err := tr.Repo.Insert(ctx, order)
		require.NoError(t, err)

		order, err = tr.Repo.Get(ctx, order.ID)

		require.NoError(t, err)
		require.Equal(t, s.expectedOrder.ID, order.ID)
	})
}

func (s *IntTestSuite) TestGetReturns() {
	s.T().Run("TestGetReturns", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		tr := s.createRepoAndSchema(ctx)
		defer tr.DropSchema(ctx, t)
		order := s.expectedOrder
		order.Returned = true
		_, err := tr.Repo.Insert(ctx, order)
		require.NoError(t, err)

		orders, err := tr.Repo.GetReturns(ctx, 0, 10)

		require.NoError(t, err)
		require.Equal(t, s.expectedOrder.ID, orders[0].ID)
	})
}

func (s *IntTestSuite) TestGetOrders() {
	s.T().Run("TestGetOrders", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		tr := s.createRepoAndSchema(ctx)
		defer tr.DropSchema(ctx, t)
		order := s.expectedOrder
		_, err := tr.Repo.Insert(ctx, order)
		require.NoError(t, err)

		orders, err := tr.Repo.GetOrders(ctx, order.ID, 0, 10)

		require.NoError(t, err)
		require.Equal(t, s.expectedOrder.ID, orders[0].ID)
	})
}
