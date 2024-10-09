package database

import (
	"errors"
	"github.com/gopi-frame/cache"
	"gorm.io/gorm"
	"time"
)

type Cache struct {
	db        *gorm.DB
	prefix    string
	expire    time.Duration
	tableName string
}

func New(config *Config) *Cache {
	if config.DB == nil {
		panic("db is required")
	}
	if config.Prefix == "" {
		config.Prefix = "cache"
	}
	if config.Expire <= 0 {
		config.Expire = time.Hour * 72
	}
	if config.TableName == "" {
		config.TableName = "caches"
	}
	if !config.DB.Migrator().HasTable(config.TableName) {
		if err := config.DB.Table(config.TableName).Migrator().CreateTable(new(CacheModel)); err != nil {
			panic(err)
		}
	}
	c := &Cache{
		db:        config.DB,
		prefix:    config.Prefix,
		expire:    config.Expire,
		tableName: config.TableName,
	}
	go c.gc()
	return c
}

func (c *Cache) buildKey(key string) string {
	return c.prefix + ":" + key
}

func (c *Cache) gc() {
	for {
		time.Sleep(time.Minute)
		if err := c.db.Table(c.tableName).Where("expire < ?", time.Now()).Delete(&CacheModel{}).Error; err != nil {
			panic(err)
		}
	}
}

func (c *Cache) Get(key string) (string, error) {
	var model = new(CacheModel)
	if err := c.db.Table(c.tableName).Where("key = ?", c.buildKey(key)).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", cache.ErrCacheNotFound
		}
		return "", err
	}
	if model.Expire.Before(time.Now()) {
		defer func() {
			c.db.Table(c.tableName).Delete(&model)
		}()
		return "", cache.ErrCacheNotFound
	}
	return model.Value, nil
}

func (c *Cache) Set(key string, value string, expire time.Duration) error {
	if expire <= 0 {
		expire = c.expire
	}
	model := &CacheModel{
		Key:    c.buildKey(key),
		Value:  value,
		Expire: time.Now().Add(expire),
	}
	return c.db.Table(c.tableName).Save(model).Error
}

func (c *Cache) Load(key string, loader func() (string, error), expire time.Duration) (string, error) {
	if v, err := c.Get(key); err == nil {
		return v, nil
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

func (c *Cache) Has(key string) bool {
	var model = new(CacheModel)
	if err := c.db.Table(c.tableName).Where("key = ?", c.buildKey(key)).First(&model).Error; err != nil {
		return false
	}
	if model.Expire.Before(time.Now()) {
		defer func() {
			c.db.Table(c.tableName).Delete(&model)
		}()
		return false
	}
	return model.Expire.After(time.Now())
}

func (c *Cache) Delete(key string) error {
	return c.db.Table(c.tableName).Where("key = ?", c.buildKey(key)).Delete(new(CacheModel)).Error
}

func (c *Cache) Clear() error {
	return c.db.Table(c.tableName).Where("key LIKE ?", c.prefix+":%").Delete(new(CacheModel)).Error
}
