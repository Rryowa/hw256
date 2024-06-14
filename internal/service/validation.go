package service

import (
	"errors"
	"homework-1/internal/entities"
	"homework-1/internal/storage"
	"homework-1/internal/util"
	"strconv"
	"time"
)

type ValidationService interface {
	AcceptValidation(id, userId, dateStr string) error
	ReturnToCourierValidation(id string) error
	IssueValidation(ids []string) error
	ReturnValidation(id, userId string) error
	ListReturnsValidation(page, size string) ([]entities.Order, error)
	ListOrdersValidation(userId, limit string) ([]entities.Order, error)
}

type orderValidator struct {
	repository   *storage.SQLRepository
	orderService OrderService
}

func NewOrderValidator(repository *storage.SQLRepository, orderService OrderService) ValidationService {
	return &orderValidator{
		repository:   repository,
		orderService: orderService,
	}
}

func (v *orderValidator) AcceptValidation(id, userId, dateStr string) error {
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

	if v.orderService.Exists(id) {
		return util.ErrOrderExists
	}

	order := v.orderService.Accept(id, userId, storageUntil)
	return v.repository.Insert(order)
}

func (v *orderValidator) IssueValidation(ids []string) error {
	if len(ids) == 0 {
		return util.ErrUserIdNotProvided
	}

	var orders []entities.Order
	var recipientID string
	for i, id := range ids {
		if !v.orderService.Exists(id) {
			return util.ErrOrderNotFound
		}

		order := v.orderService.Get(id)

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

	modifiedOrders := v.orderService.IssueOrders(orders)

	return v.repository.IssueUpdate(modifiedOrders)
}

func (v *orderValidator) ReturnValidation(id, userId string) error {
	if len(id) == 0 {
		return util.ErrOrderIdNotProvided
	}

	if len(userId) == 0 {
		return util.ErrUserIdNotProvided
	}

	if !v.orderService.Exists(id) {
		return util.ErrOrderNotFound
	}

	order := v.orderService.Get(id)

	if order.UserID != userId {
		return util.ErrOrderDoesNotBelong
	}
	if !order.Issued {
		return util.ErrOrderNotIssued
	}
	if time.Now().After(order.IssuedAt.Add(48 * time.Hour)) {
		return util.ErrReturnPeriodExpired
	}

	orderModified := v.orderService.Return(order)
	return v.repository.Update(orderModified)
}

func (v *orderValidator) ReturnToCourierValidation(id string) error {
	if len(id) == 0 {
		return util.ErrOrderIdNotProvided
	}

	if _, err := strconv.Atoi(id); err != nil {
		return util.ErrOrderIdInvalid
	}

	if !v.orderService.Exists(id) {
		return util.ErrOrderNotFound
	}

	order := v.orderService.Get(id)

	if order.Issued {
		return util.ErrOrderIssued
	}

	//skip checking for a period, to ensure that its working
	//if time.Now().Before(order.StorageUntil) {
	//	return util.ErrOrderNotExpired
	//}

	return v.repository.Delete(id)
}

func (v *orderValidator) ListReturnsValidation(limit, offset string) ([]entities.Order, error) {
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return nil, err
	}
	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		return nil, err
	}

	return v.repository.ListReturns(limitInt, offsetInt)
}

func (v *orderValidator) ListOrdersValidation(userId, limit string) ([]entities.Order, error) {
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return nil, err
	}

	return v.repository.ListOrders(userId, limitInt)
}
