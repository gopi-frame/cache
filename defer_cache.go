package cache

import (
	"sync"
	"time"

	"github.com/gopi-frame/contract/cache"
)

type DeferCache struct {
	cache.Cache

	once   sync.Once
	driver string
	config map[string]any
}

func NewDeferCache(driver string, config map[string]any) *DeferCache {
	return &DeferCache{
		driver: driver,
		config: config,
	}
}

func (c *DeferCache) deferInit() {
	c.once.Do(func() {
		var err error
		if c.Cache, err = Open(c.driver, c.config); err != nil {
			panic(err)
		}
	})
}

func (c *DeferCache) Get(key string) (string, error) {
	c.deferInit()
	return c.Cache.Get(key)
}

func (c *DeferCache) Set(key string, value string, expire time.Duration) error {
	c.deferInit()
	return c.Cache.Set(key, value, expire)
}

func (c *DeferCache) Load(key string, loader func() (value string, err error), expire time.Duration) (string, error) {
	c.deferInit()
	return c.Cache.Load(key, loader, expire)
}

func (c *DeferCache) Delete(key string) error {
	c.deferInit()
	return c.Cache.Delete(key)
}

func (c *DeferCache) Has(key string) bool {
	c.deferInit()
	return c.Cache.Has(key)
}

func (c *DeferCache) Clear() error {
	c.deferInit()
	return c.Cache.Clear()
}
