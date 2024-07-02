package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"homework/internal/models"
	"homework/internal/storage"
	"homework/internal/storage/db"
	"log"
)

type TestRepository struct {
	Repo storage.Storage
}

func NewTestRepository(ctx context.Context, cfg *models.Config) TestRepository {
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatal(err, "db connection error")
	}
	log.Println("Connected to db")
	return TestRepository{
		Repo: &db.Repository{
			Pool: pool,
		},
	}
}
