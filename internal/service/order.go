package service

import (
	"fmt"
	"homework-1/internal/models"
	"homework-1/internal/storage"
	"homework-1/pkg/hash"
	"time"
)

type OrderService interface {
	Exists(id string) bool
	Get(id string) models.Order
	Delete(id string)
	Return(order models.Order) models.Order
	Accept(id, userId string, storageUntil time.Time) models.Order
	IssueOrders(orders []models.Order) []models.Order
}

type orderService struct {
	cache *storage.Cache
}

func NewOrderService() OrderService {
	return &orderService{
		cache: storage.NewCache(),
	}
}

func (os *orderService) Exists(id string) bool { return os.cache.Exist(id) }

func (os *orderService) Get(id string) models.Order { return os.cache.Get(id) }

func (os *orderService) Delete(id string) { os.Delete(id) }

func (os *orderService) Return(order models.Order) models.Order {
	order.Returned = true
	os.cache.Update(order)

	return order
}

func (os *orderService) Accept(id, userId string, storageUntil time.Time) models.Order {
	order := models.Order{
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

	go func(order *models.Order, ticker *time.Ticker, done chan struct{}) {
		order.Hash = hash.GenerateHash()
		ticker.Stop()
		done <- struct{}{}
	}(&order, ticker, done)

	<-done

	os.cache.Update(order)

	return order
}

func (os *orderService) IssueOrders(orders []models.Order) []models.Order {
	modifiedOrders := make([]models.Order, 0)
	for _, order := range orders {
		order.Issued = true
		order.IssuedAt = time.Now()
		os.cache.Update(order)
		modifiedOrders = append(modifiedOrders, order)
	}

	return modifiedOrders
}
