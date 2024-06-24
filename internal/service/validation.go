package service

import (
	"errors"
	"homework/internal/models"
	pkg "homework/internal/service/package"
	"homework/internal/storage"
	"homework/internal/util"
	"strconv"
	"time"
)

type ValidationService interface {
	ValidateAccept(dto models.Dto) (models.Order, error)
	ValidateIssue(ids []string) ([]models.Order, error)
	ValidateAcceptReturn(id, userId string) (models.Order, error)
	ValidateReturnToCourier(id string) error
	ValidateList(offset, limit string) (int, int, error)
}

type validationService struct {
	repository     storage.Storage
	packageService pkg.PackageService
}

func NewValidationService(repository storage.Storage, packageService pkg.PackageService) ValidationService {
	return &validationService{
		repository:     repository,
		packageService: packageService,
	}
}

func (v *validationService) ValidateAccept(dto models.Dto) (models.Order, error) {
	if len(dto.ID) == 0 {
		return models.Order{}, util.ErrOrderIdNotProvided
	}
	if len(dto.UserID) == 0 {
		return models.Order{}, util.ErrUserIdNotProvided
	}
	if len(dto.Weight) == 0 {
		return models.Order{}, util.ErrWeightNotProvided
	}
	if len(dto.OrderPrice) == 0 {
		return models.Order{}, util.ErrPriceNotProvided
	}

	storageUntil, err := time.Parse(time.DateOnly, dto.StorageUntil)
	if err != nil {
		return models.Order{}, errors.New("error parsing date")
	} else if storageUntil.Before(time.Now()) {
		return models.Order{}, util.ErrDateInvalid
	}

	orderPriceFloat, err := strconv.ParseFloat(dto.OrderPrice, 64)
	if err != nil || orderPriceFloat <= 0 {
		return models.Order{}, util.ErrOrderPriceInvalid
	}

	weightFloat, err := strconv.ParseFloat(dto.Weight, 64)
	if err != nil || weightFloat <= 0 {
		return models.Order{}, util.ErrWeightInvalid
	}

	if v.repository.Exists(dto.ID) {
		return models.Order{}, util.ErrOrderExists
	}

	orderPrice := models.Price(orderPriceFloat)
	weight := models.Weight(weightFloat)
	packageType := models.PackageType(dto.PackageType)

	if err = v.packageService.ValidatePackage(weight, packageType); err != nil {
		return models.Order{}, err
	}

	order := models.Order{
		ID:           dto.ID,
		UserID:       dto.UserID,
		StorageUntil: storageUntil,
		OrderPrice:   orderPrice,
		Weight:       weight,
	}

	return order, nil
}

func (v *validationService) ValidateIssue(ids []string) ([]models.Order, error) {
	var emptyOrders []models.Order
	var ordersToIssue []models.Order

	if len(ids) == 0 {
		return emptyOrders, util.ErrUserIdNotProvided
	}

	if !v.repository.Exists(ids[0]) {
		return emptyOrders, util.ErrOrderNotFound
	}
	order, err := v.repository.Get(ids[0])
	if err != nil {
		return emptyOrders, err
	}

	recipientID := order.UserID

	for _, id := range ids {
		if !v.repository.Exists(id) {
			return emptyOrders, util.ErrOrderNotFound
		}
		order, err = v.repository.Get(id)
		if err != nil {
			return emptyOrders, err
		}
		if order.Issued {
			return emptyOrders, util.ErrOrderIssued
		}
		if order.Returned {
			return emptyOrders, util.ErrOrderReturned
		}
		if time.Now().After(order.StorageUntil) {
			return emptyOrders, util.ErrOrderExpired
		}

		//Check if users are equal
		if order.UserID != recipientID {
			return emptyOrders, util.ErrOrdersUserDiffers
		}

		ordersToIssue = append(ordersToIssue, order)
	}

	return ordersToIssue, nil
}

func (v *validationService) ValidateAcceptReturn(id, userId string) (models.Order, error) {
	emptyOrder := models.Order{}
	if len(id) == 0 {
		return emptyOrder, util.ErrOrderIdNotProvided
	}

	if len(userId) == 0 {
		return emptyOrder, util.ErrUserIdNotProvided
	}

	if !v.repository.Exists(id) {
		return emptyOrder, util.ErrOrderNotFound
	}
	order, err := v.repository.Get(id)
	if err != nil {
		return emptyOrder, err
	}

	if order.UserID != userId {
		return emptyOrder, util.ErrOrderDoesNotBelong
	}
	if !order.Issued {
		return emptyOrder, util.ErrOrderNotIssued
	}
	if time.Now().After(order.IssuedAt.Add(48 * time.Hour)) {
		return emptyOrder, util.ErrReturnPeriodExpired
	}

	return order, nil
}

func (v *validationService) ValidateReturnToCourier(id string) error {
	if len(id) == 0 {
		return util.ErrOrderIdNotProvided
	}

	if _, err := strconv.Atoi(id); err != nil {
		return util.ErrOrderIdInvalid
	}

	if !v.repository.Exists(id) {
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

	return nil
}

func (v *validationService) ValidateList(offset, limit string) (int, int, error) {
	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		return -1, -1, err
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return -1, -1, err
	}

	return offsetInt, limitInt, nil
}
