package cache

import (
	"sync"

	"github.com/gopi-frame/collection/kv"
	"github.com/gopi-frame/contract/cache"
)

// CacheManager is a cache manager.
type CacheManager struct {
	once sync.Once
	cache.Cache

	defaultStore string
	stores       *kv.Map[string, cache.Cache]
}

func NewCacheManager() *CacheManager {
	return &CacheManager{
		stores: kv.NewMap[string, cache.Cache](),
	}
}

// SetDefaultStore sets the default cache store name.
func (c *CacheManager) SetDefaultStore(name string) {
	c.defaultStore = name
}

// AddStore adds a cache store to the manager.
func (c *CacheManager) AddStore(name string, store cache.Cache) {
	c.stores.Lock()
	defer c.stores.Unlock()
	c.stores.Set(name, store)
}

// HasStore checks if the cache store exists.
func (c *CacheManager) HasStore(name string) bool {
	c.stores.RLock()
	if c.stores.ContainsKey(name) {
		c.stores.RUnlock()
		return true
	}
	c.stores.RUnlock()
	return false
}

// TryStore gets the cache store.
// It will return an error if the cache store is not configured or error occurs.
func (c *CacheManager) TryStore(name string) (cache.Cache, error) {
	c.stores.RLock()
	if store, ok := c.stores.Get(name); ok {
		c.stores.RUnlock()
		return store, nil
	}
	c.stores.RUnlock()
	return nil, NewStoreNotConfiguredException(name)
}

// GetStore gets the cache store.
// It will panic if the cache store is not configured or error occurs.
func (c *CacheManager) GetStore(name string) cache.Cache {
	if store, err := c.TryStore(name); err != nil {
		panic(err)
	} else {
		return store
	}
}
