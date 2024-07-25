package service

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"homework/internal/models"
	"homework/internal/storage"
	"homework/internal/util"
	"homework/pkg/hash"
	"strconv"
	"strings"
	"time"
)

type OrderService interface {
	Accept(ctx context.Context, dto models.Dto, pkgTypeStr string) error
	Issue(ctx context.Context, idsStr string) ([]models.Order, error)
	Return(ctx context.Context, id, userId string) error
	ReturnToCourier(ctx context.Context, id string) error
	ListReturns(ctx context.Context, offsetStr, limitStr string) ([]models.Order, error)
	ListOrders(ctx context.Context, userId, offsetStr, limitStr string) ([]models.Order, error)
}

type orderService struct {
	repository     storage.Storage
	packageService PackageService
	hashGenerator  hash.Hasher
}

func NewOrderService(r storage.Storage, ps PackageService, hg hash.Hasher) OrderService {
	return &orderService{
		repository:     r,
		packageService: ps,
		hashGenerator:  hg,
	}
}

func (os *orderService) Accept(ctx context.Context, dto models.Dto, pkgTypeStr string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.Accept")
	defer span.Finish()

	_, ok := os.repository.Exists(ctx, dto.ID)
	if ok {
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

	if len(pkgTypeStr) != 0 {
		packageType := models.PackageType(pkgTypeStr)
		if err = os.packageService.ValidatePackage(weight, packageType); err != nil {
			return err
		}
		os.packageService.ApplyPackage(&order, models.PackageType(pkgTypeStr))
	}

	os.calculateHash(&order)

	_, err = os.repository.Insert(ctx, order)
	if err != nil {
		return err
	}

	return nil
}

func (os *orderService) Issue(ctx context.Context, idsStr string) ([]models.Order, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.Issue")
	defer span.Finish()

	ids := strings.Split(idsStr, ",")
	order, ok := os.repository.Exists(ctx, ids[0])
	if !ok {
		return nil, util.ErrOrderNotFound
	}

	recipientID := order.UserID
	var orders []models.Order
	for _, id := range ids {
		order, ok = os.repository.Exists(ctx, id)
		if !ok {
			return nil, util.ErrOrderNotFound
		}
		if order.Issued {
			return nil, util.ErrOrderIssued
		}
		if time.Now().After(order.StorageUntil) {
			return nil, util.ErrOrderExpired
		}

		//Check if users are equal
		if order.UserID != recipientID {
			return nil, util.ErrOrdersUserDiffers
		}

		orders = append(orders, order)
	}
	for i := range orders {
		orders[i].Issued = true
	}

	return os.repository.IssueUpdate(ctx, orders)
}

func (os *orderService) Return(ctx context.Context, id, userId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.Return")
	defer span.Finish()

	order, ok := os.repository.Exists(ctx, id)
	if !ok {
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
	_, err := os.repository.Update(ctx, order)
	if err != nil {
		return err
	}

	return nil
}

func (os *orderService) ReturnToCourier(ctx context.Context, id string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.ReturnToCourier")
	defer span.Finish()

	order, ok := os.repository.Exists(ctx, id)
	if !ok {
		return util.ErrOrderNotFound
	}

	if order.Issued {
		return util.ErrOrderIssued
	}

	if time.Now().Before(order.StorageUntil) {
		return util.ErrOrderNotExpired
	}

	if _, err := os.repository.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}

func (os *orderService) ListReturns(ctx context.Context, offsetStr, limitStr string) ([]models.Order, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.ListReturns")
	defer span.Finish()

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
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.ListOrders")
	defer span.Finish()

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

func (os *orderService) calculateHash(order *models.Order) {
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
	}(order, ticker, done)

	<-done
	fmt.Println()
}