package service

import (
	"fmt"
	"homework-1/internal/entities"
	"homework-1/internal/storage"
	"homework-1/pkg/hash"
	"time"
)

type OrderService interface {
	Exists(id string) bool
	Get(id string) entities.Order
	Delete(id string)
	Return(order entities.Order) entities.Order
	Accept(id, userId string, storageUntil time.Time) entities.Order
	IssueOrders(orders []entities.Order) []entities.Order
}

type orderService struct {
	cache *storage.Cache
}

func NewOrderService() OrderService {
	return &orderService{
		cache: storage.NewCache(),
	}
}

func (os *orderService) Exists(id string) bool {
	return os.cache.Exist(id)
}

func (os *orderService) Get(id string) entities.Order {
	return os.cache.Get(id)
}

func (os *orderService) Delete(id string) {
	os.Delete(id)
}

func (os *orderService) Return(order entities.Order) entities.Order {
	order.Returned = true
	os.cache.Update(order)

	return order
}

func (os *orderService) Accept(id, userId string, storageUntil time.Time) entities.Order {
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

	os.cache.Update(order)

	return order
}

func (os *orderService) IssueOrders(orders []entities.Order) []entities.Order {
	modifiedOrders := make([]entities.Order, 0)
	for _, order := range orders {
		order.Issued = true
		order.IssuedAt = time.Now()
		os.cache.Update(order)
		modifiedOrders = append(modifiedOrders, order)
	}

	return modifiedOrders
}
