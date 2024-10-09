package database

import (
	"github.com/go-viper/mapstructure/v2"
	"github.com/gopi-frame/cache"
	cc "github.com/gopi-frame/contract/cache"
	"github.com/gopi-frame/exception"
)

var driverName = "database"

//goland:noinspection GoBoolExpressions
func init() {
	if driverName != "" {
		cache.Register(driverName, &Driver{})
	}
}

type Driver struct{}

func (d *Driver) Open(config map[string]any) (cc.Cache, error) {
	var cfg Config
	err := mapstructure.Decode(config, &cfg)
	if err != nil {
		return nil, err
	}
	if cfg.DB == nil {
		return nil, exception.NewArgumentException("db", cfg.DB, "db is required")
	}
	c := New(&cfg)
	return c, nil
}

func Open(config map[string]any) (cc.Cache, error) {
	return (new(Driver)).Open(config)
}

func OpenT[T any](config map[string]any, opts ...cache.Option[T]) (*cache.Cache[T], error) {
	c, err := Open(config)
	if err != nil {
		return nil, err
	}
	return cache.New[T](c, opts...)
}
