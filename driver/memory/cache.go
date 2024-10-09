package memory

import (
	"github.com/gopi-frame/cache"
	"sync"
	"time"
)

type Cache struct {
	expire time.Duration
	data   map[string]struct {
		value  string
		expire time.Time
	}
	mu *sync.RWMutex
}

func New(expire time.Duration) *Cache {
	if expire <= 0 {
		expire = time.Hour * 72
	}
	return &Cache{
		expire: expire,
		data: make(map[string]struct {
			value  string
			expire time.Time
		}),
		mu: &sync.RWMutex{},
	}
}

func (c *Cache) Get(key string) (string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if v, ok := c.data[key]; ok {
		if v.expire.After(time.Now()) {
			return v.value, nil
		}
		delete(c.data, key)
	}
	return "", cache.ErrCacheNotFound
}

func (c *Cache) Set(key string, value string, expire time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if expire <= 0 {
		expire = c.expire
	}
	c.data[key] = struct {
		value  string
		expire time.Time
	}{
		value:  value,
		expire: time.Now().Add(expire),
	}
	return nil
}

func (c *Cache) Load(key string, loader func() (string, error), expire time.Duration) (string, error) {
	if c.Has(key) {
		return c.Get(key)
	}
	value, err := loader()
	if err != nil {
		return "", err
	}
	if err := c.Set(key, value, expire); err != nil {
		return "", err
	}
	return value, nil
}

func (c *Cache) Has(key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if v, ok := c.data[key]; ok {
		if v.expire.After(time.Now()) {
			return true
		}
		delete(c.data, key)
	}
	return false
}

func (c *Cache) Delete(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
	return nil
}

func (c *Cache) Clear() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data = make(map[string]struct {
		value  string
		expire time.Time
	})
	return nil
}
