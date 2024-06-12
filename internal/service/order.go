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
	Return(order entities.Order) map[string]entities.Order
	ListReturns(page, pageSize int) []entities.Order
	ListOrders(userId string, limit int) []entities.Order
}

type orderService struct {
	storage *storage.OrderStorage
}

func NewOrderService(storage *storage.OrderStorage) OrderService {
	return &orderService{
		storage: storage,
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

	fmt.Print("Calculating hash.")

	ticker := time.NewTicker(time.Second)
	done := make(chan struct{})
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				fmt.Print(" .")
			}
		}
	}()

	go func(order *entities.Order, ticker *time.Ticker, done chan struct{}) {
		order.Hash = hash.GenerateHash()
		ticker.Stop()
		done <- struct{}{}
	}(&order, ticker, done)

	<-done

	//To prettify Ticker
	//time.Sleep(50 * time.Millisecond)
	fmt.Println("\nOrder accepted!")

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

func (os *orderService) Return(order entities.Order) map[string]entities.Order {
	order.Returned = true
	os.storage.Add(order)

	log.Println("Return accepted.")
	return os.storage.GetOrders()
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
