package storage

import (
	"homework-1/internal/entities"
	"sort"
	"sync"
)

type OrderStorage struct {
	Orders   map[string]entities.Order
	Mu       sync.Mutex
	OrderIDs []string
}

func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		Orders:   make(map[string]entities.Order),
		OrderIDs: []string{},
	}
}

func (ost *OrderStorage) DeleteAll(id string) map[string]entities.Order {
	orderIds := ost.GetOrderIds()
	sort.Strings(orderIds)
	index := sort.SearchStrings(orderIds, id)

	ost.Mu.Lock()
	if orderIds[index] == id {
		if len(orderIds) == 1 {
			ost.OrderIDs = []string{}
		} else {
			ost.OrderIDs = append(orderIds[:index], orderIds[index+1:]...)
		}
		delete(ost.Orders, id)
	}
	ost.Mu.Unlock()

	return ost.GetOrders()
}

func (ost *OrderStorage) UpdateAll(m map[string]entities.Order, ids []string) {
	ost.updateOrders(m)
	ost.updateOrderIds(ids)
}

func (ost *OrderStorage) Get(id string) entities.Order {
	ost.Mu.Lock()
	defer ost.Mu.Unlock()

	return ost.Orders[id]
}

func (ost *OrderStorage) GetOrders() map[string]entities.Order {
	ost.Mu.Lock()
	defer ost.Mu.Unlock()

	//We need to copy a map, to escape data race!
	m := make(map[string]entities.Order)
	for k, v := range ost.Orders {
		m[k] = v
	}
	return m
}

func (ost *OrderStorage) GetOrderIds() []string {
	ost.Mu.Lock()
	defer ost.Mu.Unlock()

	return ost.OrderIDs
}

func (ost *OrderStorage) Add(order entities.Order) {
	ost.Mu.Lock()
	defer ost.Mu.Unlock()
	ost.Orders[order.ID] = order
	ost.OrderIDs = append(ost.OrderIDs, order.ID)
}

func (ost *OrderStorage) updateOrders(m map[string]entities.Order) {
	ost.Mu.Lock()
	defer ost.Mu.Unlock()
	ost.Orders = make(map[string]entities.Order, len(m))
	ost.Orders = m
}

func (ost *OrderStorage) updateOrderIds(ids []string) {
	ost.Mu.Lock()
	defer ost.Mu.Unlock()
	ost.OrderIDs = make([]string, len(ids))
	copy(ost.OrderIDs, ids)
}

func (ost *OrderStorage) Exists(id string) bool {
	ost.Mu.Lock()
	defer ost.Mu.Unlock()
	if _, exists := ost.Orders[id]; exists {
		return true
	}
	return false
}
