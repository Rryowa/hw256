package storage

import "homework/internal/models"

type Storage interface {
	Insert(order models.Order) (string, error)
	Update(order models.Order) (bool, error)
	IssueUpdate(orders []models.Order) error
	Delete(id string) (string, error)
	Get(id string) (models.Order, error)
	GetReturns(offset, limit int) ([]models.Order, error)
	GetOrders(userId string, offset, limit int) ([]models.Order, error)
}
