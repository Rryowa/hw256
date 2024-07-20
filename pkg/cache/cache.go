package cache

import (
	"fmt"
	"homework/internal/models"
	"homework/internal/models/config"
	"homework/pkg/cache/arc"
	"homework/pkg/cache/ttl"
)

// CacheService is a wrapper for [string, order] cache
type CacheService interface {
	Get(key string) (models.Order, bool)
	Put(key string, val models.Order) error
	Delete(key string) error
}

func NewCache(cfg *config.CacheConfig) CacheService {
	var cacheService CacheService
	switch cfg.Type {
	case "ARC":
		cacheService = arc.NewArcCache[string, models.Order](cfg.Size)
	case "TTL":
		cacheService = ttl.NewTTLCache[string, models.Order](cfg.Size, cfg.TTL)
	default:
		fmt.Println("Unknown CACHE_TYPE, defaulting to ARC")
		cacheService = arc.NewArcCache[string, models.Order](cfg.Size)
	}

	return cacheService
}

//func (CacheService) Get(key string) (models.Order, bool) {
//	return c.cache.Get(key)
//}
//
//func (c *Cache) Put(key string, val models.Order) error {
//	return c.cache.Put(key, val)
//}
//
//func (c *Cache) Delete(key string) error {
//	return c.cache.Delete(key)
//}