package storage

import (
	"context"
	"github.com/jackc/pgx/v5"
	"homework/internal/models"
)

type Storage interface {
	Insert(ctx context.Context, order models.Order) (string, error)
	Update(ctx context.Context, order models.Order) (bool, error)
	IssueUpdate(ctx context.Context, orders []models.Order) ([]bool, error)
	Delete(ctx context.Context, id string) (string, error)
	Get(ctx context.Context, id string) (models.Order, error)
	GetReturns(ctx context.Context, offset, limit int) ([]models.Order, error)
	GetOrders(ctx context.Context, userId string, offset, limit int) ([]models.Order, error)
	StorageTest
	Event
}

type Event interface {
	CreateEvent(ctx context.Context, request string) error
	GetEvent(ctx context.Context) (models.Event, pgx.Tx, error)
	ProcessEvent(ctx context.Context, e models.Event, tx pgx.Tx) (models.Event, error)
}

type StorageTest interface {
	Truncate(ctx context.Context, table string) error
}