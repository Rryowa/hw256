package ttl

import (
	"sync"
	"time"
)

type TTL[K comparable, V any] struct {
	cache    map[K]cacheItem[V]
	ttl      time.Duration
	mutex    sync.RWMutex
	capacity int
}

type cacheItem[V any] struct {
	value      V
	expiration int64
}

func NewTTLCache[K comparable, V any](capacity int, ttl time.Duration) *TTL[K, V] {
	return &TTL[K, V]{
		cache:    make(map[K]cacheItem[V], capacity),
		ttl:      ttl,
		capacity: capacity,
	}
}

func (c *TTL[K, V]) Get(key K) (value V, ok bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	item, ok := c.cache[key]
	if !ok || item.expiration < time.Now().UnixNano() {
		delete(c.cache, key)
		return zeroValue[V](), false
	}
	return item.value, true
}

func (c *TTL[K, V]) Put(key K, value V) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.evict()

	if len(c.cache) >= c.capacity {
		c.removeOldest()
	}

	c.cache[key] = cacheItem[V]{
		value:      value,
		expiration: time.Now().Add(c.ttl).UnixNano(),
	}

	return nil
}

func (c *TTL[K, V]) Delete(key K) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.cache, key)

	return nil
}

func (c *TTL[K, V]) evict() {
	now := time.Now().UnixNano()
	for key, item := range c.cache {
		if item.expiration < now {
			delete(c.cache, key)
		}
	}
}

func (c *TTL[K, V]) removeOldest() {
	for len(c.cache) > 0 {
		var oldestKey K
		var oldestExpiration = time.Now().UnixNano()

		for key, item := range c.cache {
			if item.expiration < oldestExpiration {
				oldestKey = key
				oldestExpiration = item.expiration
			}
		}
		delete(c.cache, oldestKey)
	}
}

func zeroValue[V any]() V {
	var zero V
	return zero
}