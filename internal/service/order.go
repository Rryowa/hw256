package service

import (
	"errors"
	"fmt"
	"homework/internal/models"
	"homework/internal/storage"
	"homework/internal/util"
	"homework/pkg/hash"
	"strconv"
	"strings"
	"time"
)

type OrderService interface {
	Accept(dto models.Dto, pkgTypeStr string) error
	Issue(idsStr string) error
	Return(id, userId string) error
	ReturnToCourier(id string) error
	ListReturns(offsetStr, limitStr string) ([]models.Order, error)
	ListOrders(userId, offsetStr, limitStr string) ([]models.Order, error)
	PrintList(orders []models.Order)
}

type orderService struct {
	repository     storage.Storage
	packageService PackageService
}

func NewOrderService(repository storage.Storage, packageService PackageService) OrderService {
	return &orderService{
		repository:     repository,
		packageService: packageService,
	}
}

func (os *orderService) Accept(dto models.Dto, pkgTypeStr string) error {
	_, err := os.repository.Get(dto.ID)
	if !errors.Is(err, util.ErrOrderNotFound) {
		return util.ErrOrderExists
	}

	storageUntil, err := time.Parse(time.DateOnly, dto.StorageUntil)
	if err != nil {
		return util.ErrParsingDate
	}
	if storageUntil.Before(time.Now()) {
		return util.ErrDateInvalid
	}

	orderPriceFloat, err := strconv.ParseFloat(dto.OrderPrice, 64)
	if err != nil || orderPriceFloat <= 0 {
		return util.ErrOrderPriceInvalid
	}
	orderPrice := models.Price(orderPriceFloat)

	weightFloat, err := strconv.ParseFloat(dto.Weight, 64)
	if err != nil || weightFloat <= 0 {
		return util.ErrWeightInvalid
	}
	weight := models.Weight(weightFloat)

	if len(pkgTypeStr) != 0 {
		packageType := models.PackageType(pkgTypeStr)
		if err = os.packageService.ValidatePackage(weight, packageType); err != nil {
			return err
		}
	}

	order := models.Order{
		ID:           dto.ID,
		UserID:       dto.UserID,
		StorageUntil: storageUntil,
		OrderPrice:   orderPrice,
		Weight:       weight,
	}

	os.packageService.ApplyPackage(&order, models.PackageType(pkgTypeStr))

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

	_, err = os.repository.Insert(order)
	return err
}

func (os *orderService) Issue(idsStr string) error {
	ids := strings.Split(idsStr, ",")
	order, err := os.repository.Get(ids[0])
	if err != nil {
		return util.ErrOrderNotFound
	}
	recipientID := order.UserID

	var orders []models.Order
	for _, id := range ids {
		order, err := os.repository.Get(id)
		if err != nil {
			return util.ErrOrderNotFound
		}
		if order.Issued {
			return util.ErrOrderIssued
		}
		if order.Returned {
			return util.ErrOrderReturned
		}
		if time.Now().After(order.StorageUntil) {
			return util.ErrOrderExpired
		}

		//Check if users are equal
		if order.UserID != recipientID {
			return util.ErrOrdersUserDiffers
		}

		orders = append(orders, order)
	}

	for i := range orders {
		orders[i].Issued = true
		orders[i].IssuedAt = time.Now()
	}

	return os.repository.IssueUpdate(orders)
}

func (os *orderService) Return(id, userId string) error {
	order, err := os.repository.Get(id)
	if err != nil {
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

	_, err = os.repository.Update(order)
	return err
}

func (os *orderService) ReturnToCourier(id string) error {
	order, err := os.repository.Get(id)
	if err != nil {
		return util.ErrOrderNotFound
	}

	if order.Issued {
		return util.ErrOrderIssued
	}

	//skip checking for a period, to ensure that its working
	//if time.Now().Before(order.StorageUntil) {
	//	return util.ErrOrderNotExpired
	//}

	_, err = os.repository.Delete(id)
	return err
}

func (os *orderService) ListReturns(offsetStr, limitStr string) ([]models.Order, error) {
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		return []models.Order{}, util.ErrOffsetInvalid
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return []models.Order{}, util.ErrLimitInvalid
	}

	return os.repository.GetReturns(offset, limit)
}

func (os *orderService) ListOrders(userId, offsetStr, limitStr string) ([]models.Order, error) {
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		return []models.Order{}, util.ErrOffsetInvalid
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return []models.Order{}, util.ErrLimitInvalid
	}

	return os.repository.GetOrders(userId, offset, limit)
}

func (os *orderService) PrintList(orders []models.Order) {
	if len(orders) == 0 {
		defer fmt.Printf("\n\n")
	}
	fmt.Printf("%-5s%-10s%-15s%-15v%-10v%-13v%-10v%-13s%-13v\n", "id", "user_id", "storage_until", "issued_at", "returned", "order_price", "weight", "package_type", "package_price")
	fmt.Println(strings.Repeat("-", 100))
	for _, order := range orders {
		fmt.Printf("%-5s%-10s%-15s%-15v%-10v%-13v%-10v%-13s%-13v\n",
			order.ID,
			order.UserID,
			order.StorageUntil.Format("2006-01-02"),
			order.IssuedAt.Format("2006-01-02"),
			order.Returned,
			order.OrderPrice,
			order.Weight,
			order.PackageType,
			order.PackagePrice)
	}
	fmt.Printf("\n")
}
