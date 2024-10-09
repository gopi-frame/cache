package cache

import (
	"encoding/json"
	"errors"
	"github.com/gopi-frame/contract/cache"
	"time"
)

var ErrCacheNotFound = errors.New("cache not found")

// Cache is a generic cache wrapper.
type Cache[T any] struct {
	cache.Cache
	encoder func(T) ([]byte, error)
	decoder func([]byte) (T, error)
}

func New[T any](cache cache.Cache, opts ...Option[T]) (*Cache[T], error) {
	c := &Cache[T]{
		Cache: cache,
		encoder: func(value T) ([]byte, error) {
			return json.Marshal(value)
		},
		decoder: func(bs []byte) (T, error) {
			var v T
			err := json.Unmarshal(bs, &v)
			return v, err
		},
	}
	for _, opt := range opts {
		if err := opt.Apply(c); err != nil {
			return nil, err
		}
	}
	return c, nil
}

func (c *Cache[T]) Get(key string) (T, error) {
	v, err := c.Cache.Get(key)
	if err != nil {
		return *new(T), err
	}
	if v, err := c.decoder([]byte(v)); err != nil {
		return *new(T), err
	} else {
		return v, nil
	}
}

func (c *Cache[T]) Set(key string, value T, expire time.Duration) error {
	bs, err := c.encoder(value)
	if err != nil {
		return err
	}
	return c.Cache.Set(key, string(bs), expire)
}

func (c *Cache[T]) Load(key string, loader func() (T, error), expire time.Duration) (T, error) {
	var loaded bool
	var value T
	v, err := c.Cache.Load(key, func() (string, error) {
		v, err := loader()
		if err != nil {
			return "", err
		}
		loaded = true
		value = v
		bs, err := c.encoder(v)
		if err != nil {
			return "", err
		}
		return string(bs), nil
	}, expire)
	if err != nil {
		return *new(T), err
	}
	if loaded {
		return value, nil
	}
	return c.decoder([]byte(v))
}

func (c *Cache[T]) Delete(key string) error {
	return c.Cache.Delete(key)
}

func (c *Cache[T]) Has(key string) bool {
	return c.Cache.Has(key)
}

func (c *Cache[T]) Clear() error {
	return c.Cache.Clear()
}
