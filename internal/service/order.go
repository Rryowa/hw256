package service

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"homework/internal/metrics"
	"homework/internal/models"
	"homework/internal/storage"
	"homework/internal/util"
	"homework/pkg/cache"
	"homework/pkg/hash"
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
}

type orderService struct {
	repository     storage.Storage
	packageService PackageService
	hashGenerator  hash.Hasher
	cache          cache.CacheService
	serverMetrics  metrics.Metrics
}

func NewOrderService(r storage.Storage, ps PackageService, hg hash.Hasher, c cache.CacheService, sm metrics.Metrics) OrderService {
	return &orderService{
		repository:     r,
		packageService: ps,
		hashGenerator:  hg,
		cache:          c,
		serverMetrics:  sm,
	}
}

func (os *orderService) Accept(ctx context.Context, dto models.Dto, pkgTypeStr string) error {
	os.serverMetrics.IncrementMethodCallCounter("Accept")
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.Accept")
	defer span.Finish()

	_, ok := os.cache.Get(dto.ID)
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

	order, err = os.repository.Insert(ctx, order)
	if err != nil {
		return err
	}

	return os.putInCache(order)
}

func (os *orderService) Issue(ctx context.Context, idsStr string) error {
	os.serverMetrics.IncrementMethodCallCounter("Issue")
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.Issue")
	defer span.Finish()

	start := time.Now()
	ids := strings.Split(idsStr, ",")

	order, ok := os.cache.Get(ids[0])
	if !ok {
		return util.ErrOrderNotFound
	}

	recipientID := order.UserID
	var orders []models.Order
	for _, id := range ids {
		order, ok = os.cache.Get(id)
		if !ok {
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
	}

	issuedOrders, err := os.repository.IssueUpdate(ctx, orders)
	if err != nil {
		return err
	}

	for _, order := range issuedOrders {
		if err := os.putInCache(order); err != nil {
			return err
		}
		os.serverMetrics.IncrementIssuedCounter()
	}
	os.serverMetrics.ObserveRequestDuration("200", time.Since(start))

	return nil
}

func (os *orderService) Return(ctx context.Context, id, userId string) error {
	os.serverMetrics.IncrementMethodCallCounter("Return")
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.Return")
	defer span.Finish()

	order, ok := os.cache.Get(id)
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
	order, err := os.repository.Update(ctx, order)
	if err != nil {
		return err
	}

	if err := os.putInCache(order); err != nil {
		return err
	}

	return nil
}

func (os *orderService) ReturnToCourier(ctx context.Context, id string) error {
	os.serverMetrics.IncrementMethodCallCounter("ReturnToCourier")
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.ReturnToCourier")
	defer span.Finish()

	order, ok := os.cache.Get(id)
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

	if err := os.cache.Delete(id); err != nil {
		return err
	}

	return nil
}

func (os *orderService) ListReturns(ctx context.Context, offsetStr, limitStr string) ([]models.Order, error) {
	os.serverMetrics.IncrementMethodCallCounter("ListReturns")
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
	os.serverMetrics.IncrementMethodCallCounter("ListOrders")
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

func (os *orderService) putInCache(order models.Order) error {
	err := os.cache.Put(order.ID, order)
	if err != nil {
		return err
	}
	return nil
}