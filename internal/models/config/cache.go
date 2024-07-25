package config

import "time"

type CacheConfig struct {
	Type   string
	TTL    time.Duration
	Period time.Duration
	Size   int
}