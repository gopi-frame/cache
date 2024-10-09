package cache

import (
	"fmt"
	"github.com/gopi-frame/collection/kv"
	"github.com/gopi-frame/contract/cache"
	"github.com/gopi-frame/exception"
)

var drivers = kv.NewMap[string, cache.Driver]()

// Register registers driver
func Register(driverName string, driver cache.Driver) {
	drivers.Lock()
	defer drivers.Unlock()
	if driver == nil {
		panic(exception.NewEmptyArgumentException("driver"))
	}
	if drivers.ContainsKey(driverName) {
		panic(exception.NewArgumentException("driverName", driverName, fmt.Sprintf("duplicate driver \"%s\"", driverName)))
	}
	drivers.Set(driverName, driver)
}

// Drivers lists registered drivers
func Drivers() []string {
	drivers.RLock()
	defer drivers.RUnlock()
	list := drivers.Keys()
	return list
}

// Open opens cache.
func Open(driverName string, config map[string]any) (cache.Cache, error) {
	drivers.RLock()
	driver, ok := drivers.Get(driverName)
	drivers.RUnlock()
	if !ok {
		return nil, exception.NewArgumentException("driverName", driverName, fmt.Sprintf("unknown driver \"%s\"", driverName))
	}
	return driver.Open(config)
}

// OpenT opens cache with type.
func OpenT[T any](driverName string, config map[string]any, opts ...Option[T]) (*Cache[T], error) {
	drivers.RLock()
	driver, ok := drivers.Get(driverName)
	drivers.RUnlock()
	if !ok {
		return nil, exception.NewArgumentException("driverName", driverName, fmt.Sprintf("unknown driver \"%s\"", driverName))
	}
	c, err := driver.Open(config)
	if err != nil {
		return nil, err
	}
	return New[T](c, opts...)
}
