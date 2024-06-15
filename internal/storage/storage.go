package storage

import "homework-1/internal/models"

// Storage To easily replace postgres with any other db
type Storage interface {
	Insert(order models.Order) error
	Update(order models.Order) error
	IssueUpdate(orders []models.Order) error
	Delete(id string) error
	Exists(id string) bool
	Get(id string) models.Order
	GetReturns(limit, offset int) ([]models.Order, error)
	GetOrders(userId string, limit int) ([]models.Order, error)
	AnalyzeQueryPlan(query string, args ...interface{}) error
}
