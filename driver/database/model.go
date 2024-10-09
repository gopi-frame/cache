package database

import "time"

type CacheModel struct {
	Key    string    `gorm:"column:key;type:varchar(255);not null;primaryKey;uniqueIndex" json:"key" yaml:"key" toml:"key" mapstructure:"key"`
	Value  string    `gorm:"column:value;type:text;not null" json:"value" yaml:"value" toml:"value" mapstructure:"value"`
	Expire time.Time `gorm:"column:expire;type:timestamp;not null" json:"expire" yaml:"expire" toml:"expire" mapstructure:"expire"`
}
