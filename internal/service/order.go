package service

import (
	"context"
	"errors"
	"fmt"
	"homework/internal/models"
	"homework/internal/storage"
	"homework/internal/util"
	"homework/pkg/hash"
	"homework/pkg/timer"
	"strconv"
	"strings"
	"time"
)

type OrderService interface {
	Accept(ctx context.Context, dto models.Dto, pkgTypeStr string) error
	Issue(ctx context.Context, idsStr string) error
	Return(ctx context.Context, id, userId string) error
	ReturnToCourier(ctx context.Context, id string) error
	ListReturns(ctx context.Context, offsetStr, limitStr string) ([]models.Order, error)
	ListOrders(ctx context.Context, userId, offsetStr, limitStr string) ([]models.Order, error)
	PrintList(orders []models.Order)
}

type orderService struct {
	repository     storage.Storage
	packageService PackageService
	hashGenerator  hash.Hasher
	timeGenerator  timer.Timer
}

func NewOrderService(repository storage.Storage, packageService PackageService, hashGenerator hash.Hasher, timeGenerator timer.Timer) OrderService {
	return &orderService{
		repository:     repository,
		packageService: packageService,
		hashGenerator:  hashGenerator,
		timeGenerator:  timeGenerator,
	}
}

func (os *orderService) Accept(ctx context.Context, dto models.Dto, pkgTypeStr string) error {

	_, err := os.repository.Get(ctx, dto.ID)
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

	order := models.Order{
		ID:           dto.ID,
		UserID:       dto.UserID,
		StorageUntil: storageUntil,
		OrderPrice:   orderPrice,
		Weight:       weight,
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
		order.Hash = os.hashGenerator.GenerateHash()
		ticker.Stop()
		done <- struct{}{}
	}(&order, ticker, done)

	<-done
	fmt.Println()

	if len(pkgTypeStr) != 0 {
		packageType := models.PackageType(pkgTypeStr)
		if err = os.packageService.ValidatePackage(weight, packageType); err != nil {
			return err
		}
		os.packageService.ApplyPackage(&order, models.PackageType(pkgTypeStr))
	}

	_, err = os.repository.Insert(ctx, order)
	return err
}

func (os *orderService) Issue(ctx context.Context, idsStr string) error {
	ids := strings.Split(idsStr, ",")
	order, err := os.repository.Get(ctx, ids[0])
	if err != nil {
		return util.ErrOrderNotFound
	}
	recipientID := order.UserID

	var orders []models.Order
	for _, id := range ids {
		order, err := os.repository.Get(ctx, id)
		if err != nil {
			return util.ErrOrderNotFound
		}
		if order.Issued {
			return util.ErrOrderIssued
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
		orders[i].IssuedAt = os.timeGenerator.TimeNow()
	}

	_, err = os.repository.IssueUpdate(ctx, orders)
	return err
}

func (os *orderService) Return(ctx context.Context, id, userId string) error {
	order, err := os.repository.Get(ctx, id)
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

	_, err = os.repository.Update(ctx, order)
	return err
}

func (os *orderService) ReturnToCourier(ctx context.Context, id string) error {
	order, err := os.repository.Get(ctx, id)
	if err != nil {
		return util.ErrOrderNotFound
	}

	if order.Issued {
		return util.ErrOrderIssued
	}

	if time.Now().Before(order.StorageUntil) {
		return util.ErrOrderNotExpired
	}

	_, err = os.repository.Delete(ctx, id)
	return err
}

func (os *orderService) ListReturns(ctx context.Context, offsetStr, limitStr string) ([]models.Order, error) {
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		return []models.Order{}, util.ErrOffsetInvalid
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return []models.Order{}, util.ErrLimitInvalid
	}

	return os.repository.GetReturns(ctx, offset, limit)
}

func (os *orderService) ListOrders(ctx context.Context, userId, offsetStr, limitStr string) ([]models.Order, error) {
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		return []models.Order{}, util.ErrOffsetInvalid
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return []models.Order{}, util.ErrLimitInvalid
	}

	return os.repository.GetOrders(ctx, userId, offset, limit)
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
