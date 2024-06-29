//go:build integration

package tests

import (
	"context"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"hash/maphash"
	"homework/internal/models"
	storage "homework/internal/storage/db"
	"homework/internal/util"
	"homework/tests/db"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

// TODO: Instead of creating a New database per test create a database at the start of
// the tests and then create a schema per test and set the search_ path to only be of
// that scehema When connecting to the database it's possible to add a connection string
// parameter `?search_path=example` to execute all queries by default in that schema.
// To apply changes make down and make run
type IntTestSuite struct {
	suite.Suite
	createQuery   string
	dropQuery     string
	expectedOrder models.Order
}

func TestIntTestSuite(t *testing.T) {
	suite.Run(t, new(IntTestSuite))
}

func (s *IntTestSuite) SetupTest() {
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

func (s *IntTestSuite) TestInsertOrder() {
	s.T().Run("TestInsertOrder", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		cfg := util.NewTestConfig()
		name := rand.New(rand.NewSource(int64(new(maphash.Hash).Sum64())))
		schemaName := "test" + strconv.FormatInt(name.Int63(), 10)
		tr := db.NewTestRepository(ctx, cfg, schemaName)
		tr.Repo.Pool.Exec(ctx, "CREATE SCHEMA "+schemaName)

		query := storage.AddSchemaPrefix(schemaName, s.createQuery)

		_, err := tr.Repo.Pool.Exec(ctx, query)

		require.NoError(t, err)

		id, err := tr.Repo.Insert(s.expectedOrder, schemaName)
		require.NoError(t, err)
		require.Equal(t, s.expectedOrder.ID, id)

		tr.Repo.Pool.Exec(ctx, "DROP SCHEMA "+schemaName+" CASCADE")
	})
}

func (s *IntTestSuite) TestUpdateOrder() {
	s.T().Run("TestUpd", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		cfg := util.NewTestConfig()
		name := rand.New(rand.NewSource(int64(new(maphash.Hash).Sum64())))
		schemaName := "test" + strconv.FormatInt(name.Int63(), 10)
		tr := db.NewTestRepository(ctx, cfg, schemaName)
		tr.Repo.Pool.Exec(ctx, "CREATE SCHEMA "+schemaName)
		query := storage.AddSchemaPrefix(schemaName, s.createQuery)

		_, err := tr.Repo.Pool.Exec(ctx, query)
		require.NoError(t, err)

		order := s.expectedOrder
		order.Returned = true
		_, err = tr.Repo.Insert(order, schemaName)
		require.NoError(t, err)

		returned, err := tr.Repo.Update(order, schemaName)

		require.NoError(t, err)
		require.Equal(t, order.Returned, returned)

		tr.Repo.Pool.Exec(ctx, "DROP SCHEMA "+schemaName+" CASCADE")
	})
}

func (s *IntTestSuite) TestDelete() {
	s.T().Run("TestDelete", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		cfg := util.NewTestConfig()
		name := rand.New(rand.NewSource(int64(new(maphash.Hash).Sum64())))
		schemaName := "test" + strconv.FormatInt(name.Int63(), 10)
		tr := db.NewTestRepository(ctx, cfg, schemaName)
		tr.Repo.Pool.Exec(ctx, "CREATE SCHEMA "+schemaName)
		query := storage.AddSchemaPrefix(schemaName, s.createQuery)

		//Create table inside a new scheme
		_, err := tr.Repo.Pool.Exec(ctx, query)
		require.NoError(t, err)

		order := s.expectedOrder
		_, err = tr.Repo.Insert(order, schemaName)
		require.NoError(t, err)

		id, err := tr.Repo.Delete(order.ID, schemaName)

		require.NoError(t, err)
		require.Equal(t, order.ID, id)

		tr.Repo.Pool.Exec(ctx, "DROP SCHEMA "+schemaName+" CASCADE")
	})
}

func (s *IntTestSuite) TestGet() {
	s.T().Run("TestGet", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		cfg := util.NewTestConfig()
		name := rand.New(rand.NewSource(int64(new(maphash.Hash).Sum64())))
		schemaName := "test" + strconv.FormatInt(name.Int63(), 10)
		tr := db.NewTestRepository(ctx, cfg, schemaName)
		tr.Repo.Pool.Exec(ctx, "CREATE SCHEMA "+schemaName)
		query := storage.AddSchemaPrefix(schemaName, s.createQuery)

		//Create table inside a new scheme
		_, err := tr.Repo.Pool.Exec(ctx, query)
		require.NoError(t, err)

		order := s.expectedOrder
		_, err = tr.Repo.Insert(order, schemaName)
		require.NoError(t, err)

		order, err = tr.Repo.Get(order.ID, schemaName)

		require.NoError(t, err)
		require.Equal(t, s.expectedOrder.ID, order.ID)

		tr.Repo.Pool.Exec(ctx, "DROP SCHEMA "+schemaName+" CASCADE")
	})
}

func (s *IntTestSuite) TestGetReturns() {
	s.T().Run("TestGetReturns", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		cfg := util.NewTestConfig()
		name := rand.New(rand.NewSource(int64(new(maphash.Hash).Sum64())))
		schemaName := "test" + strconv.FormatInt(name.Int63(), 10)
		tr := db.NewTestRepository(ctx, cfg, schemaName)
		tr.Repo.Pool.Exec(ctx, "CREATE SCHEMA "+schemaName)
		query := storage.AddSchemaPrefix(schemaName, s.createQuery)

		//Create table inside a new scheme
		_, err := tr.Repo.Pool.Exec(ctx, query)
		require.NoError(t, err)

		order := s.expectedOrder
		order.Returned = true
		_, err = tr.Repo.Insert(order, schemaName)
		require.NoError(t, err)

		orders, err := tr.Repo.GetReturns(0, 10, schemaName)

		require.NoError(t, err)
		require.Equal(t, s.expectedOrder.ID, orders[0].ID)

		tr.Repo.Pool.Exec(ctx, "DROP SCHEMA "+schemaName+" CASCADE")
	})
}

func (s *IntTestSuite) TestGetOrders() {
	s.T().Run("TestGetReturns", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		cfg := util.NewTestConfig()
		name := rand.New(rand.NewSource(int64(new(maphash.Hash).Sum64())))
		schemaName := "test" + strconv.FormatInt(name.Int63(), 10)
		tr := db.NewTestRepository(ctx, cfg, schemaName)
		tr.Repo.Pool.Exec(ctx, "CREATE SCHEMA "+schemaName)
		query := storage.AddSchemaPrefix(schemaName, s.createQuery)

		//Create table inside a new scheme
		_, err := tr.Repo.Pool.Exec(ctx, query)
		require.NoError(t, err)

		order := s.expectedOrder
		_, err = tr.Repo.Insert(order, schemaName)
		require.NoError(t, err)

		orders, err := tr.Repo.GetOrders(order.ID, 0, 10, schemaName)

		require.NoError(t, err)
		require.Equal(t, s.expectedOrder.ID, orders[0].ID)
	})
}
