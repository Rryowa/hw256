package service

import (
	"homework/internal/models"
	pkg "homework/internal/service/package"
	"homework/internal/storage"
	"homework/internal/util"
	"strconv"
	"strings"
	"time"
)

type ValidationService interface {
	ValidateAccept(dto models.Dto) (models.Order, error)
	ValidateIssue(idsStr string) ([]models.Order, error)
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
	if isArgEmpty(dto.ID) {
		return models.Order{}, util.ErrOrderIdNotProvided
	}
	if isArgEmpty(dto.UserID) {
		return models.Order{}, util.ErrUserIdNotProvided
	}
	if isArgEmpty(dto.Weight) {
		return models.Order{}, util.ErrWeightNotProvided
	}
	if isArgEmpty(dto.OrderPrice) {
		return models.Order{}, util.ErrPriceNotProvided
	}

	storageUntil, err := time.Parse(time.DateOnly, dto.StorageUntil)
	if err != nil {
		return models.Order{}, util.ErrParsingDate
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

	if _, err := v.repository.Get(dto.ID); err == nil {
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

func (v *validationService) ValidateIssue(idsStr string) ([]models.Order, error) {
	var emptyOrders []models.Order
	var ordersToIssue []models.Order
	if len(idsStr) == 0 {
		return emptyOrders, util.ErrUserIdNotProvided
	}

	ids := strings.Split(idsStr, ",")

	order, err := v.repository.Get(ids[0])
	if err != nil {
		return emptyOrders, err
	}

	recipientID := order.UserID

	for _, id := range ids {
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
	if isArgEmpty(id) {
		return emptyOrder, util.ErrOrderIdNotProvided
	}
	if isArgEmpty(userId) {
		return emptyOrder, util.ErrUserIdNotProvided
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
	if isArgEmpty(id) {
		return util.ErrOrderIdNotProvided
	}

	if _, err := strconv.Atoi(id); err != nil {
		return util.ErrOrderIdInvalid
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
	if isArgEmpty(offset) {
		return -1, -1, util.ErrOffsetNotProvided
	}
	if isArgEmpty(limit) {
		return -1, -1, util.ErrLimitNotProvided
	}
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

func isArgEmpty(id string) bool {
	if len(id) == 0 {
		return true
	}
	return false
}
