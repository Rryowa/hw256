package storage

import "homework/internal/models"

type Storage interface {
	Insert(order models.Order, schemaName string) (string, error)
	Update(order models.Order, schemaName string) (bool, error)
	IssueUpdate(orders []models.Order, schemaName string) error
	Delete(id string, schemaName string) (string, error)
	Get(id string, schemaName string) (models.Order, error)
	GetReturns(offset, limit int, schemaName string) ([]models.Order, error)
	GetOrders(userId string, offset, limit int, schemaName string) ([]models.Order, error)
}
