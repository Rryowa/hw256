package service

import (
	"fmt"
	"homework-1/internal/models"
	"homework-1/pkg/hash"
	"time"
)

type OrderService interface {
	Return(order models.Order) models.Order
	Accept(id, userId string, storageUntil time.Time) models.Order
	IssueOrders(orders []models.Order) []models.Order
}

func Return(order models.Order) models.Order {
	order.Returned = true
	return order
}

func Accept(id, userId string, storageUntil time.Time) models.Order {
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

	return order
}

func IssueOrders(orders []models.Order) []models.Order {
	modifiedOrders := make([]models.Order, 0)
	for _, order := range orders {
		order.Issued = true
		order.IssuedAt = time.Now()
		modifiedOrders = append(modifiedOrders, order)
	}

	return modifiedOrders
}
