package cache

import (
	cache "github.com/kiricle/in-memory-cache"
)

type Cache struct {
	Storage *cache.Cache
}

func NewCache() *Cache {
	storage := cache.New()

	return &Cache{Storage: storage}
}
