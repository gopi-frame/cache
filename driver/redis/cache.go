package redis

import (
	"context"
	"github.com/gopi-frame/cache"
	"github.com/gopi-frame/contract/redis"
	"time"
)

type Cache struct {
	client redis.Client
	prefix string
	expire time.Duration
}

// New creates a new cache.
func New(config *Config) *Cache {
	if config.Client == nil {
		panic("client is required")
	}
	if config.Prefix == "" {
		config.Prefix = "cache"
	}
	if config.Expire <= 0 {
		config.Expire = time.Hour * 72
	}
	return &Cache{
		client: config.Client,
		prefix: config.Prefix,
		expire: config.Expire,
	}
}

func (c *Cache) buildKey(key string) string {
	return c.prefix + ":" + key
}

func (c *Cache) Get(key string) (string, error) {
	if !c.Has(key) {
		return "", cache.ErrCacheNotFound
	}
	return c.client.Get(context.Background(), c.buildKey(key)).Result()
}

func (c *Cache) Set(key string, value string, expire time.Duration) error {
	if expire <= 0 {
		expire = c.expire
	}
	return c.client.Set(context.Background(), c.buildKey(key), value, expire).Err()
}

func (c *Cache) Load(key string, loader func() (value string, err error), expire time.Duration) (string, error) {
	key = c.buildKey(key)
	if c.Has(key) {
		return c.Get(key)
	}
	v, err := loader()
	if err != nil {
		return "", err
	}
	if err := c.Set(key, v, expire); err != nil {
		return "", err
	}
	return v, nil
}

func (c *Cache) Delete(key string) error {
	return c.client.Del(context.Background(), c.buildKey(key)).Err()
}

func (c *Cache) Has(key string) bool {
	return c.client.Exists(context.Background(), c.buildKey(key)).Val() > 0
}

func (c *Cache) Clear() error {
	iter := c.client.Scan(context.Background(), 0, c.prefix+":*", 0).Iterator()
	for iter.Next(context.Background()) {
		if err := c.client.Del(context.Background(), iter.Val()).Err(); err != nil {
			return err
		}
	}
	return iter.Err()
}
