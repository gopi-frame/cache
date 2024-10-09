package file

import (
	"github.com/go-viper/mapstructure/v2"
	"github.com/gopi-frame/cache"
)
import cc "github.com/gopi-frame/contract/cache"

// This variable can be replaced through `go build -ldflags=-X github.com/gopi-frame/cache/driver/file.driverName=custom`
var driverName = "file"

//goland:noinspection GoBoolExpressions
func init() {
	if driverName != "" {
		cache.Register(driverName, new(Driver))
	}
}

type Driver struct{}

func (d *Driver) Open(config map[string]any) (cc.Cache, error) {
	var cfg Config
	err := mapstructure.Decode(config, &cfg)
	if err != nil {
		return nil, err
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
