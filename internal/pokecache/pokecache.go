package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	cache    map[string]cacheEntry
	mux      *sync.Mutex
	interval time.Duration
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) Cache {
	cache := Cache{
		cache: map[string]cacheEntry{},
		mux:   &sync.Mutex{},
	}
	go cache.reapLoop()
	return cache
}

func (c *Cache) Add(key string, val []byte) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.cache[key] = cacheEntry{
		createdAt: time.Now().UTC(),
		val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mux.Lock()
	defer c.mux.Unlock()
	val, ok := c.cache[key]
	return val.val, ok
}

func (c *Cache) reapLoop() {
	for key, entry := range c.cache {
		c.mux.Lock()
		defer c.mux.Unlock()
		if entry.createdAt.Before(time.Now().Add(-1 * c.interval)) {
			delete(c.cache, key)
		}
	}
}
