package storage

import (
	"context"
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
	InsertEvent(ctx context.Context, request string) (models.Event, error)
}

type StorageTest interface {
	Truncate(ctx context.Context, table string) error
}