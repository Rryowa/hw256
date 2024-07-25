package storage

import (
	"context"
	"homework/internal/models"
)

type Storage interface {
	Insert(ctx context.Context, order models.Order) (models.Order, error)
	Update(ctx context.Context, order models.Order) (models.Order, error)
	IssueUpdate(ctx context.Context, orders []models.Order) ([]models.Order, error)
	Delete(ctx context.Context, id string) (string, error)
	GetReturns(ctx context.Context, offset, limit int) ([]models.Order, error)
	GetOrders(ctx context.Context, userId string, offset, limit int) ([]models.Order, error)
	Exists(ctx context.Context, id string) (models.Order, bool)
	Event
}

type Event interface {
	InsertEvent(ctx context.Context, request string) (models.Event, error)
	UpdateEvent(ctx context.Context, event models.Event) (models.Event, error)
}