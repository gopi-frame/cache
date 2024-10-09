package database

import (
	"gorm.io/gorm"
	"time"
)

type Config struct {
	// DB is the database connection.
	DB *gorm.DB `json:"db" yaml:"db" toml:"db" mapstructure:"db"`
	// Prefix is the cache prefix.
	Prefix string `json:"prefix" yaml:"prefix" toml:"prefix" mapstructure:"prefix"`
	// Expire is the default cache expire time, default is 72 hour.
	Expire time.Duration `json:"expire" yaml:"expire" toml:"expire" mapstructure:"expire"`
	// Table is the cache table name.
	TableName string `json:"table_name" yaml:"table_name" toml:"table_name" mapstructure:"table_name"`
}
