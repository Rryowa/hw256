package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"homework/internal/models"
	"homework/internal/storage/db"
	"log"
	"testing"
)

type TestRepository struct {
	Repo   *db.Repository
	schema string
}

func NewTestRepository(ctx context.Context, cfg *models.Config, schema string) TestRepository {
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?search_path=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, schema)
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatal(err, "db connection error")
	}
	log.Println("Connected to db")

	return TestRepository{
		Repo: &db.Repository{
			Pool: pool,
		},
		schema: schema,
	}
}

func (tr *TestRepository) DropSchema(ctx context.Context, t *testing.T) {
	_, err := tr.Repo.Pool.Exec(ctx, "DROP SCHEMA "+tr.schema+" CASCADE")
	require.NoError(t, err)
}
