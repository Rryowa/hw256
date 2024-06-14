package storage

import (
	"homework-1/internal/models"
	"sync"
)

type Cache struct {
	cache map[string]models.Order
	mu    sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{
		cache: make(map[string]models.Order),
	}
}

func (c *Cache) Get(id string) models.Order {
	c.mu.RLock()
	defer c.mu.RUnlock()
	order := c.cache[id]

	return order
}

func (c *Cache) Exist(id string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if _, exist := c.cache[id]; exist {
		return true
	} else {
		return false
	}
}

func (c *Cache) Update(order models.Order) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	c.cache[order.ID] = order
}

func (c *Cache) Delete(id string) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	delete(c.cache, id)
}
