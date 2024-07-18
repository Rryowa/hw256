package cache

import (
	log "github.com/sirupsen/logrus"
	"homework/internal/models"
	"homework/internal/models/config"
	"homework/pkg/cache/arc"
)

// CacheService is a wrapper for [string, order] cache
type CacheService interface {
	Get(key string) (models.Order, bool)
	Put(key string, val models.Order) error
	Delete(key string) error
	Len() int
}

type Cache struct {
	arc *arc.ARC[string, models.Order]
}

func NewCache(cfg *config.CacheConfig) CacheService {
	log.Debugln("Cache SIZE:", cfg.Size)
	return &Cache{
		arc: arc.NewArcCache[string, models.Order](cfg.Size),
	}
}

func (c *Cache) Get(key string) (models.Order, bool) {
	return c.arc.Get(key)
}

func (c *Cache) Put(key string, val models.Order) error {
	return c.arc.Put(key, val)
}

func (c *Cache) Delete(key string) error {
	return c.arc.Delete(key)
}

func (c *Cache) Len() int {
	return c.arc.Len()
}