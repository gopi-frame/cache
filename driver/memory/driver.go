package memory

import (
	"github.com/gopi-frame/cache"
	cc "github.com/gopi-frame/contract/cache"
	"time"
)

// This variable can be replaced through `go build -ldflags=-X github.com/gopi-frame/cache/driver/memory.driverName=custom`
var driverName = "memory"

//goland:noinspection GoBoolExpressions
func init() {
	if driverName != "" {
		cache.Register(driverName, &Driver{})
	}
}

type Driver struct{}

func (d *Driver) Open(config map[string]any) (cc.Cache, error) {
	return New(config["expire"].(time.Duration)), nil
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
