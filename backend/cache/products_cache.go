package cache

import (
	"sync"
	"time"
)

type CacheEntry struct {
	Data      []byte
	ExpiresAt time.Time
}

type ProductsCache struct {
	mu      sync.RWMutex
	entries map[string]*CacheEntry
	ttl     time.Duration
}

func NewProductsCache(ttl time.Duration) *ProductsCache {
	return &ProductsCache{
		entries: make(map[string]*CacheEntry),
		ttl:     ttl,
	}
}

func (c *ProductsCache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.entries[key]
	if !exists || entry == nil {
		return nil, false
	}

	if time.Now().After(entry.ExpiresAt) {
		// Entry expired, but don't delete here (let it be cleaned up on next Set)
		return nil, false
	}

	return entry.Data, true
}

func (c *ProductsCache) Set(key string, data []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Clean up expired entries
	c.cleanupExpired()

	c.entries[key] = &CacheEntry{
		Data:      data,
		ExpiresAt: time.Now().Add(c.ttl),
	}
}

func (c *ProductsCache) cleanupExpired() {
	now := time.Now()
	for k, entry := range c.entries {
		if entry == nil || now.After(entry.ExpiresAt) {
			delete(c.entries, k)
		}
	}
}

func (c *ProductsCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries = make(map[string]*CacheEntry)
}
