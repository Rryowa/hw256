package service

import (
	"fmt"
	"homework-1/internal/entities"
	"homework-1/internal/storage"
	"homework-1/pkg/hash"
	"log"
	"sort"
	"time"
)

type OrderService interface {
	AcceptOrder(id, userId, dateStr string) map[string]entities.Order
	ReturnOrderToCourier(orderID string) map[string]entities.Order
	IssueOrders(OrderIDs []string) map[string]entities.Order
	AcceptReturn(orderID string) map[string]entities.Order
	ListReturns(page, pageSize int) []entities.Order
	ListOrders(userId string, limit int) []entities.Order
	UpdateCache() error
}

type orderService struct {
	storage     *storage.OrderStorage
	fileService FileService
}

func NewOrderService(storage *storage.OrderStorage, service FileService) OrderService {
	return &orderService{
		storage:     storage,
		fileService: service,
	}
}

func (os *orderService) AcceptOrder(id, userId, dateStr string) map[string]entities.Order {
	storageUntil, _ := time.Parse(time.DateOnly, dateStr)
	order := entities.Order{
		ID:           id,
		UserID:       userId,
		Issued:       false,
		Returned:     false,
		StorageUntil: storageUntil,
	}

	fmt.Println("Calculating hash...")
	// to skip generating in tests
	order.Hash = hash.GenerateHash()

	log.Println("Order accepted.")

	os.storage.Add(order)
	return os.storage.GetOrders()
}

func (os *orderService) ReturnOrderToCourier(orderID string) map[string]entities.Order {
	log.Println("Order returned.")
	return os.storage.DeleteAll(orderID)
}

func (os *orderService) IssueOrders(OrderIDs []string) map[string]entities.Order {
	for _, id := range OrderIDs {
		order := os.storage.Get(id)
		order.Issued = true
		order.IssuedAt = time.Now()

		os.storage.Add(order)

		log.Println("Order issued.")
	}
	return os.storage.GetOrders()
}

func (os *orderService) AcceptReturn(orderID string) map[string]entities.Order {
	order := os.storage.Get(orderID)
	order.Returned = true
	os.storage.Add(order)

	log.Println("Return accepted.")
	return os.storage.Orders
}

func (os *orderService) ListReturns(page, pageSize int) []entities.Order {
	orderIds := os.storage.GetOrderIds()
	ln := len(orderIds)
	start := (page - 1) * pageSize
	if start >= ln {
		return nil
	}

	end := start + pageSize
	if end > ln {
		end = ln
	}

	var returns []entities.Order
	for _, id := range orderIds[start:end] {
		order := os.storage.Get(id)
		if order.Returned {
			returns = append(returns, order)
		}
	}

	// Calculate the start and end indices for slicing
	returnsStart := 0
	if start < len(returns) {
		returnsStart = start
	}
	returnsEnd := end - start
	if returnsEnd > len(returns) {
		returnsEnd = len(returns)
	}

	return returns[returnsStart:returnsEnd]
}

func (os *orderService) ListOrders(userId string, limit int) []entities.Order {
	var userOrders []entities.Order
	orderIds := os.storage.GetOrderIds()

	for _, id := range orderIds {
		order := os.storage.Get(id)
		if order.UserID == userId && !order.Issued {
			userOrders = append(userOrders, order)
			if len(userOrders) == limit {
				break
			}
		}
	}

	// Sort Orders by StorageUntil date in descending order
	sort.Slice(userOrders, func(i, j int) bool {
		return userOrders[i].StorageUntil.Before(userOrders[j].StorageUntil)
	})

	return userOrders
}

func (os *orderService) UpdateCache() error {
	fmt.Println("Updating cache...")
	if os.fileService == nil {
		return fmt.Errorf("fileService is nil")
	}

	f, err := os.fileService.IsEmpty()
	if err != nil {
		return err
	}
	if f {
		return nil
	}
	orders, orderIDs, err := os.fileService.Read()
	if err != nil {
		return err
	}

	os.storage.UpdateAll(orders, orderIDs)

	return nil
}
