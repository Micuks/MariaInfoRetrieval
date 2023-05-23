package query_process

import (
	. "MariaInfoRetrieval/maria_types"
	"sync"
)

// Simple in-memory cache
type Cache struct {
	Mu    sync.Mutex
	Cache map[string][]SearchResult
}

// Create a new cache
func NewCache() *Cache {
	return &Cache{
		Cache: make(map[string][]SearchResult),
	}
}

// Retrieve an item from cache
func (c *Cache) Get(key string) ([]SearchResult, bool) {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	val, found := c.Cache[key]
	return val, found
}

func (c *Cache) Set(key string, val []SearchResult) {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	c.Cache[key] = val
}
