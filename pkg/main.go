package main

import (
	"homework/internal/models"
	"homework/pkg/cache/arc"
	"log"
)

func main() {
	cache := arc.NewArcCache[string, models.Order](32768)
	order := models.Order{
		ID:     "1",
		UserID: "2",
		Issued: true,
	}

	cache.Put("Hello1", order)
	val, _ := cache.Get("Hello1")
	log.Println(val)
	cache.Put("Hello2", order)
	log.Println(val)
	cache.Delete("Hello1")

	val, ok := cache.Get("Hello1")
	if ok {
		log.Println(val)
	} else {
		log.Println(ok)
	}
}