package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"homework/internal/models"
	"homework/internal/storage"
	"homework/internal/storage/db"
	"log"
	"testing"
)

type TestRepository struct {
	Storage storage.Storage
}

func NewTestRepository(ctx context.Context, cfg *models.Config) TestRepository {
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatal(err, "db connection error")
	}

	log.Println("Connected to db")

	return TestRepository{
		&db.Repository{
			Pool: pool,
			Ctx:  ctx,
		},
	}
}

func (r *TestRepository) SetUp(t *testing.T) {
	t.Helper()
}

func (r *TestRepository) TearDown(t *testing.T) {
	t.Helper()
}
