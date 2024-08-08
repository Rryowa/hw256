package main

import (
	"fmt"
	"homework/internal/models"
	"homework/internal/storage/cache"
	"homework/internal/util"
	"os"
	"time"
)

func main() {
	cfg := util.NewCacheConfig()
	cacher := cache.NewCache(cfg)
	cacher.Put("1", models.Order{ID: "1"})
	cacher.Put("2", models.Order{ID: "1"})
	v, ok := cacher.Get("1")
	if !ok {
		os.Exit(1)
	}
	time.Sleep(5 * time.Second)
	v, ok = cacher.Get("1")
	if !ok {
		fmt.Println("No such key")
		os.Exit(1)
	}
	fmt.Println(v)
}