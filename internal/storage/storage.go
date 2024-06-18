package storage

import "homework-1/internal/models"

// Storage To easily replace postgres with any other db
type Storage interface {
	Insert(order models.Order) error
	Update(order models.Order) error
	IssueUpdate(orders []models.Order) error
	Delete(id string) error
	Get(id string) models.Order
	GetReturns(offset, limit int) ([]models.Order, error)
	GetOrders(userId string, offset, limit int) ([]models.Order, error)
}
