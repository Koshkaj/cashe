package cache

import (
	"fmt"
	"sync"
	"time"
)

type Cache struct {
	mu   sync.RWMutex
	data map[string][]byte
}

func New() Cacher {
	return &Cache{
		data: make(map[string][]byte),
	}
}

func (c *Cache) Set(key, value []byte, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[string(key)] = value

	// if ttl > 0 {
	// 	go func() {
	// 		// Move to the one goroutine with locking (channels)
	// 		<-time.After(ttl)
	// 		delete(c.data, string(key))
	// 	}()
	// }
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

	return val, nil
}

func (c *Cache) Delete(key []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.data, string(key))
	return nil
}

func (c *Cache) Has(key []byte) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	_, ok := c.data[string(key)]
	return ok
}
