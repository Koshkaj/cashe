package cache

import (
	"fmt"
	"sync"
	"time"
)

type Cache struct {
	mu   sync.RWMutex
	data map[string][]byte

	expiryMu sync.RWMutex
	Expiry   map[string]time.Time
}

func New() Cacher {
	return &Cache{
		data:   make(map[string][]byte),
		Expiry: make(map[string]time.Time),
	}
}

func (c *Cache) Set(key, value []byte, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[string(key)] = value

	if ttl > 0 {
		expiryTime := time.Now().Add(ttl * time.Second)
		c.expiryMu.Lock()
		c.Expiry[string(key)] = expiryTime
		c.expiryMu.Unlock()
	}
	return nil
}

func (c *Cache) Get(key []byte) ([]byte, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	keyStr := string(key)
	val, ok := c.data[keyStr]
	if !ok {
		return nil, fmt.Errorf("key `%s` not found", keyStr)
	}
	expiry, ok := c.Expiry[keyStr]
	if ok && time.Now().After(expiry) {
		c.mu.Lock()
		delete(c.data, string(key))
		c.mu.Unlock()

		c.expiryMu.Lock()
		delete(c.Expiry, string(key))
		c.expiryMu.Unlock()

		return nil, nil
	}

	return val, nil
}

func (c *Cache) Delete(key []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.data, string(key))
	delete(c.Expiry, string(key))
	return nil
}

func (c *Cache) Has(key []byte) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	_, ok := c.data[string(key)]
	return ok
}
