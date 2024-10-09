package file

import (
	"os"
	"time"
)

// Config is the configuration for the file cache.
type Config struct {
	// StoragePath is the path to the storage directory.
	// If not set, the default is the "gopi-frame" directory in the user's cache directory.
	StoragePath string `json:"storagePath" yaml:"storagePath" toml:"storagePath" mapstructure:"storagePath"`
	// Prefix is the prefix for the cache files.
	// If not set, the default is "cache".
	Prefix string `json:"prefix" yaml:"prefix" toml:"prefix" mapstructure:"prefix"`
	// Expire is the default expiration time for the cache files.
	// If not set, the default is 72 hours.
	Expire time.Duration `json:"expire" yaml:"expire" toml:"expire" mapstructure:"expire"`
	// DirMode is the mode for the cache directory.
	// If not set, the default is 0755.
	DirMode os.FileMode `json:"dirMode" yaml:"dirMode" toml:"dirMode" mapstructure:"dirMode"`
	// FileMode is the mode for the cache files.
	// If not set, the default is 0644.
	FileMode os.FileMode `json:"fileMode" yaml:"fileMode" toml:"fileMode" mapstructure:"fileMode"`
}
