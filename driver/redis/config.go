package redis

import (
	"github.com/gopi-frame/contract/redis"
	"time"
)

// Config is the redis cache config.
type Config struct {
	// Client is the redis client.
	Client redis.Client `json:"client" yaml:"client" toml:"client" mapstructure:"client"`
	// Prefix is the cache prefix.
	Prefix string `json:"prefix" yaml:"prefix" toml:"prefix" mapstructure:"prefix"`
	// Expire is the default cache expire time, default is 72 hour.
	Expire time.Duration `json:"expire" yaml:"expire" toml:"expire" mapstructure:"expire"`
}
