package service

import (
	"errors"
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
	SetMaxGoroutinesValidation(ns string) error
}

type orderValidator struct {
	storage      *storage.OrderStorage
	orderService OrderService
	fileService  FileService
}

func NewOrderValidator(storage *storage.OrderStorage, orderService OrderService, fileService FileService) ValidationService {
	return &orderValidator{
		storage:      storage,
		orderService: orderService,
		fileService:  fileService,
	}
}

func (v *orderValidator) AcceptValidation(id, userId, dateStr string) error {
	storageUntil, err := time.Parse(time.DateOnly, dateStr)
	if err != nil {
		return errors.New("error parsing date")
	}
	if storageUntil.Before(time.Now()) {
		return util.InvalidDateError{}
	}

	if len(id) == 0 {
		return util.OrderIdsNotProvidedError{}
	}

	if len(userId) == 0 {
		return util.UserIdIsNotProvided{}
	}

	if v.storage.Exists(id) {
		return util.ExistingOrderError{}
	}

	orders := v.orderService.AcceptOrder(id, userId, dateStr)

	return v.fileService.Write(orders)
}

func (v *orderValidator) ReturnToCourierValidation(id string) error {
	if len(id) == 0 {
		return util.OrderIdsNotProvidedError{}
	}

	if !v.storage.Exists(id) {
		return util.OrderNotFoundError{}
	}

	_, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	return nil
}

func (v *orderValidator) IssueValidation(ids []string) error {
	if len(ids) == 0 {
		return util.UserIdIsNotProvided{}
	}

	var recipientID string
	for i, id := range ids {
		if !v.storage.Exists(id) {
			return util.OrderNotFoundError{}
		}
		order := v.storage.Get(id)

		if time.Now().After(order.StorageUntil) {
			return util.OrderIsExpiredError{}
		}
		if order.Issued {
			return util.OrderIssuedError{}
		}
		if order.Returned {
			return util.OrdersReturnedError{}
		}

		//check if recipients equal
		if i == 0 {
			recipientID = order.UserID
		} else {
			if order.UserID != recipientID {
				return util.OrdersRecipientDiffersError{}
			}
		}
	}
	return nil
}

func (v *orderValidator) ReturnValidation(id, userId string) error {
	if len(id) == 0 {
		return util.OrderIdsNotProvidedError{}
	}

	if len(userId) == 0 {
		return util.UserIdIsNotProvided{}
	}

	if !v.storage.Exists(id) {
		return util.OrderNotFoundError{}
	}
	order := v.storage.Get(id)

	if order.UserID != userId {
		return util.OrderDoesNotBelongError{}
	}
	if !order.Issued {
		return util.OrderHasNotBeenIssuedError{}
	}
	if time.Now().After(order.IssuedAt.Add(48 * time.Hour)) {
		return util.OrderCantBeReturnedError{}
	}
	return nil
}

func (v *orderValidator) SetMaxGoroutinesValidation(ns string) error {
	if len(ns) == 0 {
		return errors.New("number of goroutines is required")
	}
	n, err := strconv.Atoi(ns)
	if err != nil {
		return errors.Join(err, errors.New("invalid argument"))
	}
	if n < 1 {
		return errors.New("number of goroutines must be > 0")
	}
	return nil
}
