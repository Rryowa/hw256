package service

import (
	"errors"
	"fmt"
	"homework-1/internal/models"
	"homework-1/internal/storage"
	"homework-1/internal/util"
	"homework-1/pkg/hash"
	"strconv"
	"time"
)

type OrderService interface {
	Accept(id, userId, dateStr string) error
	ReturnToCourier(id string) error
	Issue(ids []string) error
	Return(id, userId string) error
	ListReturns(page, size string) ([]models.Order, error)
	ListOrders(userId, limit string) ([]models.Order, error)
}

type orderService struct {
	repository storage.Storage
}

func NewOrderService(repository storage.Storage) OrderService {
	return &orderService{
		repository: repository,
	}
}

func (os *orderService) Accept(id, userId, dateStr string) error {
	storageUntil, err := time.Parse(time.DateOnly, dateStr)
	if err != nil {
		return errors.New("error parsing date")
	}
	if storageUntil.Before(time.Now()) {
		return util.ErrInvalidDate
	}

	if len(id) == 0 {
		return util.ErrOrderIdNotProvided
	}

	if len(userId) == 0 {
		return util.ErrUserIdNotProvided
	}

	emptyOrder := models.Order{}
	order := os.repository.Get(id)
	if order != emptyOrder {
		return util.ErrOrderExists
	}

	newOrder := Create(id, userId, storageUntil)
	return os.repository.Insert(newOrder)
}

func (os *orderService) Issue(ids []string) error {
	if len(ids) == 0 {
		return util.ErrUserIdNotProvided
	}

	var orders []models.Order
	var recipientID string
	for i, id := range ids {
		emptyOrder := models.Order{}
		order := os.repository.Get(id)
		if order == emptyOrder {
			return util.ErrOrderNotFound
		}

		if time.Now().After(order.StorageUntil) {
			return util.ErrOrderExpired
		}
		if order.Issued {
			return util.ErrOrderIssued
		}
		if order.Returned {
			return util.ErrOrderReturned
		}

		//Check if users are equal
		if i == 0 {
			recipientID = order.UserID
		} else {
			if order.UserID != recipientID {
				return util.ErrOrdersUserDiffers
			}
		}
		orders = append(orders, order)
	}

	modifiedOrders := IssueOrders(orders)
	return os.repository.IssueUpdate(modifiedOrders)
}

func (os *orderService) Return(id, userId string) error {
	if len(id) == 0 {
		return util.ErrOrderIdNotProvided
	}

	if len(userId) == 0 {
		return util.ErrUserIdNotProvided
	}

	emptyOrder := models.Order{}
	order := os.repository.Get(id)
	if order == emptyOrder {
		return util.ErrOrderNotFound
	}

	if order.UserID != userId {
		return util.ErrOrderDoesNotBelong
	}
	if !order.Issued {
		return util.ErrOrderNotIssued
	}
	if time.Now().After(order.IssuedAt.Add(48 * time.Hour)) {
		return util.ErrReturnPeriodExpired
	}

	order.Returned = true
	return os.repository.Update(order)
}

func (os *orderService) ReturnToCourier(id string) error {
	if len(id) == 0 {
		return util.ErrOrderIdNotProvided
	}

	if _, err := strconv.Atoi(id); err != nil {
		return util.ErrOrderIdInvalid
	}

	emptyOrder := models.Order{}
	order := os.repository.Get(id)
	if order == emptyOrder {
		return util.ErrOrderNotFound
	}

	if order.Issued {
		return util.ErrOrderIssued
	}

	//skip checking for a period, to ensure that its working
	//if time.Now().Before(order.StorageUntil) {
	//	return util.ErrOrderNotExpired
	//}

	return os.repository.Delete(id)
}

func (os *orderService) ListReturns(limit, offset string) ([]models.Order, error) {
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return nil, err
	}
	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		return nil, err
	}

	return os.repository.GetReturns(limitInt, offsetInt)
}

func (os *orderService) ListOrders(userId, limit string) ([]models.Order, error) {
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return nil, err
	}

	return os.repository.GetOrders(userId, limitInt)
}

func Create(id, userId string, storageUntil time.Time) models.Order {
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
