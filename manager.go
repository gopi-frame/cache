package cache

import (
	"github.com/gopi-frame/collection/kv"
	"github.com/gopi-frame/contract/cache"
	"sync"
)

// CacheManager is a cache manager.
type CacheManager struct {
	once sync.Once
	cache.Cache

	defaultStore string
	stores       *kv.Map[string, cache.Cache]
	deferStores  *kv.Map[string, func() (cache.Cache, error)]
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

// Use sets the default cache.
// It will only set the default cache once.
func (c *CacheManager) Use(cache cache.Cache) *CacheManager {
	c.once.Do(func() {
		if c.Cache == nil {
			c.Cache = cache
		}
	})
	return c
}

// AddStore adds a cache store to the manager.
func (c *CacheManager) AddStore(name string, store cache.Cache) {
	c.stores.Lock()
	defer c.stores.Unlock()
	c.stores.Set(name, store)
}

// AddDeferStore adds a defer cache store to the manager.
func (c *CacheManager) AddDeferStore(name string, config map[string]any) {
	c.stores.Lock()
	defer c.stores.Unlock()
	c.deferStores.Set(name, func() (cache.Cache, error) {
		driverName := config["driver"].(string)
		return Open(driverName, config)
	})
}

// HasStore checks if the cache store exists.
func (c *CacheManager) HasStore(name string) bool {
	c.stores.RLock()
	if c.stores.ContainsKey(name) {
		c.stores.RUnlock()
		return true
	}
	c.stores.RUnlock()
	c.deferStores.RLock()
	if c.deferStores.ContainsKey(name) {
		c.deferStores.RUnlock()
		return true
	}
	c.deferStores.RUnlock()
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
	c.deferStores.RLock()
	if store, ok := c.deferStores.Get(name); ok {
		c.deferStores.RUnlock()
		if store, err := store(); err != nil {
			return nil, err
		} else {
			c.stores.Lock()
			defer c.stores.Unlock()
			c.stores.Set(name, store)
			return store, nil
		}
	}
	c.deferStores.RUnlock()
	return nil, NewStoreNotConfiguredException(name)
}

// Store gets the cache store.
// It will panic if the cache store is not configured or error occurs.
func (c *CacheManager) Store(name string) cache.Cache {
	if store, err := c.TryStore(name); err != nil {
		panic(err)
	} else {
		return store
	}
}

// StoreOrDefault gets the cache store.
// It will return the default cache store if the cache store is not configured or error occurs.
func (c *CacheManager) StoreOrDefault(name string) cache.Cache {
	store, err := c.TryStore(name)
	if err != nil {
		return c.Store(c.defaultStore)
	}
	return store
}
