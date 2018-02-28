package utils

import (
	"sync"
	"time"
)

type cacheItem struct {
	data    []byte
	expTime time.Time
}

type Cache struct {
	guard    sync.RWMutex
	cache    map[string]*cacheItem
	duration time.Duration
	ticker   *time.Ticker
}

func NewCache(duration time.Duration) *Cache {
	cache := &Cache{
		cache:    make(map[string]*cacheItem),
		duration: duration,
		ticker:   time.NewTicker(time.Second),
	}

	go func() {
		for range cache.ticker.C {
			cache.removeExpired()
		}
	}()

	return cache
}

func (c *Cache) AddItem(key string, data []byte) {
	c.guard.Lock()
	defer c.guard.Unlock()

	c.cache[key] = &cacheItem{
		data:    data,
		expTime: time.Now().Add(c.duration),
	}
}

func (c *Cache) GetItem(key string) ([]byte, bool) {
	c.guard.RLock()
	defer c.guard.RUnlock()

	i, ok := c.cache[key]
	if !ok {
		return nil, false
	}

	return i.data, true
}

func (c *Cache) removeExpired() {
	c.guard.Lock()
	defer c.guard.Unlock()

	t := time.Now()
	for k, v := range c.cache {
		if t.After(v.expTime) {
			delete(c.cache, k)
		}
	}
}
