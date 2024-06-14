package service

import (
	"fmt"
	"homework-1/internal/entities"
	"homework-1/pkg/hash"
	"time"
)

func Accept(id, userId string, storageUntil time.Time) entities.Order {
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

	return order
}

func IssueOrders(orders []entities.Order) []entities.Order {
	for _, order := range orders {
		order.Issued = true
		order.IssuedAt = time.Now()
	}

	return orders
}

func Return(order entities.Order) entities.Order {
	order.Returned = true

	return order
}
