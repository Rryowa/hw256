package config

import "time"

type CacheConfig struct {
	Type string
	TTL  time.Duration
	Size int
}