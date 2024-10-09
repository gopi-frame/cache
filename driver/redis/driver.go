package redis

import (
	"github.com/go-viper/mapstructure/v2"
	"github.com/gopi-frame/cache"
	cc "github.com/gopi-frame/contract/cache"
	"github.com/gopi-frame/exception"
)

// This variable can be replaced through `go build -ldflags=-X github.com/gopi-frame/cache/driver/redis.driverName=custom`
var driverName = "redis"

//goland:noinspection GoBoolExpressions
func init() {
	if driverName != "" {
		cache.Register(driverName, &Driver{})
	}
}

type Driver struct{}

func (d *Driver) Open(config map[string]any) (cc.Cache, error) {
	var cfg Config
	err := mapstructure.WeakDecode(config, &cfg)
	if err != nil {
		return nil, err
	}
	if cfg.Client == nil {
		return nil, exception.NewArgumentException("client", cfg.Client, "client is required")
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
