package service

import (
	"context"
	"homework/internal/metrics"
	"homework/internal/models"
	"time"
)

const statusOk = "ok"

type MetricsService struct {
	OrderService
}

func NewMetricsService(os OrderService) *MetricsService {
	return &MetricsService{
		OrderService: os,
	}
}

func (mos *MetricsService) Accept(ctx context.Context, dto models.Dto, pkgTypeStr string) error {
	const op = "Accept"
	startTime := time.Now()
	metrics.IncrementMethodCallCounter(op)
	defer func() {
		duration := time.Since(startTime)
		metrics.ObserveRequestDuration(statusOk, duration)
	}()

	return mos.OrderService.Accept(ctx, dto, pkgTypeStr)
}

func (mos *MetricsService) Issue(ctx context.Context, idsStr string) ([]models.Order, error) {
	const op = "Issue"
	startTime := time.Now()
	metrics.IncrementMethodCallCounter(op)
	defer func() {
		duration := time.Since(startTime)
		metrics.ObserveRequestDuration(statusOk, duration)
	}()

	issuedOrders, err := mos.OrderService.Issue(ctx, idsStr)
	if err != nil {
		return nil, err
	}
	for range issuedOrders {
		metrics.IncrementIssuedOrdersCounter()
	}

	return issuedOrders, nil
}

func (mos *MetricsService) Return(ctx context.Context, id, userId string) error {
	const op = "Return"
	startTime := time.Now()
	metrics.IncrementMethodCallCounter(op)
	defer func() {
		duration := time.Since(startTime)
		metrics.ObserveRequestDuration(statusOk, duration)
	}()

	return mos.OrderService.Return(ctx, id, userId)
}

func (mos *MetricsService) ReturnToCourier(ctx context.Context, id string) error {
	const op = "ReturnToCourier"
	startTime := time.Now()
	metrics.IncrementMethodCallCounter(op)
	defer func() {
		duration := time.Since(startTime)
		metrics.ObserveRequestDuration(statusOk, duration)
	}()

	return mos.OrderService.ReturnToCourier(ctx, id)
}

func (mos *MetricsService) ListReturns(ctx context.Context, offsetStr, limitStr string) ([]models.Order, error) {
	const op = "ListReturns"
	startTime := time.Now()
	metrics.IncrementMethodCallCounter(op)
	defer func() {
		duration := time.Since(startTime)
		metrics.ObserveRequestDuration(statusOk, duration)
	}()

	return mos.OrderService.ListReturns(ctx, offsetStr, limitStr)
}

func (mos *MetricsService) ListOrders(ctx context.Context, userId, offsetStr, limitStr string) ([]models.Order, error) {
	const op = "ListOrders"
	startTime := time.Now()
	metrics.IncrementMethodCallCounter(op)
	defer func() {
		duration := time.Since(startTime)
		metrics.ObserveRequestDuration(statusOk, duration)
	}()

	return mos.OrderService.ListOrders(ctx, userId, offsetStr, limitStr)
}