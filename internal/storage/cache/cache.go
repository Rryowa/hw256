package cache

import (
	"fmt"
	"homework/internal/models"
	"homework/internal/models/config"
	"homework/internal/storage/cache/arc"
	"homework/internal/storage/cache/ttl"
)

// Cacher is a wrapper for [string, order] cache
type Cacher interface {
	Get(key string) (models.Order, bool)
	Put(key string, val models.Order) error
	Delete(key string) error
}

func NewCache(cfg *config.CacheConfig) Cacher {
	var cacheService Cacher
	switch cfg.Type {
	case "ARC":
		cacheService = arc.NewArcCache[string, models.Order](cfg.Size)
	case "TTL":
		cacheService = ttl.NewTTLCache[string, models.Order](cfg.Size, cfg.TTL, cfg.Period)
	default:
		fmt.Println("Unknown CACHE_TYPE, defaulting to ARC")
		cacheService = arc.NewArcCache[string, models.Order](cfg.Size)
	}

	return cacheService
}