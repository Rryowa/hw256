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
		return util.ErrInvalidDate
	}

	if len(id) == 0 {
		return util.ErrOrderIdNotProvided
	}

	if len(userId) == 0 {
		return util.ErrUserIdNotProvided
	}

	if v.storage.Exists(id) {
		return util.ErrOrderExists
	}

	return v.fileService.Write(v.orderService.AcceptOrder(id, userId, dateStr))
}

func (v *orderValidator) ReturnToCourierValidation(id string) error {
	if len(id) == 0 {
		return util.ErrOrderIdNotProvided
	}

	if _, err := strconv.Atoi(id); err != nil {
		return util.ErrOrderIdInvalid
	}

	order := v.storage.Get(id)
	if order.Issued {
		return util.ErrOrderIssued
	}

	if !v.storage.Exists(id) {
		return util.ErrOrderNotFound
	}

	return v.fileService.Write(v.orderService.ReturnOrderToCourier(id))
}

func (v *orderValidator) IssueValidation(ids []string) error {
	if len(ids) == 0 {
		return util.ErrUserIdNotProvided
	}

	var recipientID string
	for i, id := range ids {
		if !v.storage.Exists(id) {
			return util.ErrOrderNotFound
		}
		order := v.storage.Get(id)

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
	}
	return v.fileService.Write(v.orderService.IssueOrders(ids))
}

func (v *orderValidator) ReturnValidation(id, userId string) error {
	if len(id) == 0 {
		return util.ErrOrderIdNotProvided
	}

	if len(userId) == 0 {
		return util.ErrUserIdNotProvided
	}

	if !v.storage.Exists(id) {
		return util.ErrOrderNotFound
	}
	order := v.storage.Get(id)

	if order.UserID != userId {
		return util.ErrOrderDoesNotBelong
	}
	if !order.Issued {
		return util.ErrOrderNotIssued
	}
	if time.Now().After(order.IssuedAt.Add(48 * time.Hour)) {
		return util.ErrReturnPeriodExpired
	}
	return v.fileService.Write(v.orderService.Return(order))
}

func (v *orderValidator) ListReturnsValidation(page, size string) ([]entities.Order, error) {
	p, err := strconv.Atoi(page)
	if err != nil {
		return nil, err
	}
	ps, err := strconv.Atoi(size)
	if err != nil {
		return nil, err
	}
	return v.orderService.ListReturns(p, ps), nil
}

func (v *orderValidator) ListOrdersValidation(userId, limit string) ([]entities.Order, error) {
	l, err := strconv.Atoi(limit)
	if err != nil {
		return nil, err
	}
	return v.orderService.ListOrders(userId, l), nil
}
