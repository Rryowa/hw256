package main

import (
	"homework/pkg/cache/arc"
	"log"
)

func main() {
	cache := arc.New(2)
	cache.Put("Hello1", "World1")
	cache.Put("Hello2", "World2")
	val, _ := cache.Get("Hello1")
	log.Println(val)

	cache.Put("Hello3", "World3")
	cache.Put("Hello4", "World4")

	//Value is not lost
	val, _ = cache.Get("Hello1")
	log.Println(val)
}