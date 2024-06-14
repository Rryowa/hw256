package service

import (
	"errors"
	"homework-1/internal/entities"
	"homework-1/internal/storage"
	"homework-1/internal/util"
	"log"
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
	repository *storage.SQLRepository
}

func NewOrderValidator(repository *storage.SQLRepository) ValidationService {
	return &orderValidator{
		repository: repository,
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

	exists, err := v.repository.Exists(id)
	if err != nil {
		return err
	}
	if exists {
		return util.ErrOrderExists
	}

	order := Accept(id, userId, storageUntil)
	return v.repository.Insert(order)
}

//TODO: WHY ISSUE OT WORKING PROPERLY IF REPEAT SAME issue -ids=2,3
//TODO IT WONT DISPLAY ERROR!

func (v *orderValidator) IssueValidation(ids []string) error {
	if len(ids) == 0 {
		return util.ErrUserIdNotProvided
	}

	var orders []entities.Order
	var recipientID string
	for i, id := range ids {
		exists, err := v.repository.Exists(id)
		if err != nil {
			return err
		}
		if !exists {
			return util.ErrOrderNotFound
		}

		order, err := v.repository.Get(id)
		if err != nil {
			return err
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

	for _, order := range IssueOrders(orders) {
		if err := v.repository.Update(order); err != nil {
			return err
		}
		log.Println("Order issued.")
	}
	return nil
}

func (v *orderValidator) ReturnValidation(id, userId string) error {
	if len(id) == 0 {
		return util.ErrOrderIdNotProvided
	}

	if len(userId) == 0 {
		return util.ErrUserIdNotProvided
	}

	exists, err := v.repository.Exists(id)
	if err != nil {
		return err
	}
	if !exists {
		return util.ErrOrderNotFound
	}
	order, err := v.repository.Get(id)
	if err != nil {
		return err
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

	return v.repository.Update(Return(order))
}

func (v *orderValidator) ReturnToCourierValidation(id string) error {
	if len(id) == 0 {
		return util.ErrOrderIdNotProvided
	}

	if _, err := strconv.Atoi(id); err != nil {
		return util.ErrOrderIdInvalid
	}

	exists, err := v.repository.Exists(id)
	if err != nil {
		return err
	}
	if !exists {
		return util.ErrOrderNotFound
	}

	order, err := v.repository.Get(id)
	if err != nil {
		return err
	}

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
