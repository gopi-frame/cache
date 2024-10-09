package file

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/gopi-frame/cache"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type Cache struct {
	mu          *sync.Mutex
	storagePath string
	prefix      string
	expire      time.Duration
	dirMode     os.FileMode
	fileMode    os.FileMode
}

func New(config *Config) *Cache {
	if config.StoragePath == "" {
		userCacheDir, err := os.UserCacheDir()
		if err != nil {
			return nil
		}
		config.StoragePath = filepath.Join(userCacheDir, "gopi-frame")
	}
	if config.DirMode == 0 {
		config.DirMode = 0755
	}
	if config.FileMode == 0 {
		config.FileMode = 0644
	}
	if config.Prefix == "" {
		config.Prefix = "cache"
	}
	if config.Expire <= 0 {
		config.Expire = time.Hour * 72
	}
	c := &Cache{
		mu:          &sync.Mutex{},
		storagePath: config.StoragePath,
		prefix:      config.Prefix,
		expire:      config.Expire,
		dirMode:     config.DirMode,
		fileMode:    config.FileMode,
	}
	go c.gc()
	return c
}

func (c *Cache) buildKey(key string) string {
	h := md5.New()
	h.Write([]byte(key))
	return c.prefix + hex.EncodeToString(h.Sum(nil))
}

func (c *Cache) buildPath(key string) string {
	return filepath.Join(c.storagePath, c.buildKey(key)+".bin")
}

func (c *Cache) gc() {
	for {
		files, err := os.ReadDir(c.storagePath)
		if err != nil {
			return
		}
		for _, file := range files {
			if strings.HasPrefix(file.Name(), c.prefix) && strings.HasSuffix(file.Name(), ".bin") {
				path := filepath.Join(c.storagePath, file.Name())
				if s, err := os.Stat(path); err != nil {
					return
				} else if s.ModTime().Before(time.Now()) {
					_ = os.Remove(path)
				}
			}
		}
		time.Sleep(time.Second)
	}
}

func (c *Cache) Get(key string) (string, error) {
	path := c.buildPath(key)
	if s, err := os.Stat(path); err != nil {
		return "", cache.ErrCacheNotFound
	} else if s.ModTime().Before(time.Now()) {
		_ = os.Remove(path)
		return "", cache.ErrCacheNotFound
	}
	content, err := os.ReadFile(path)
	if err != nil {
		_ = os.Remove(path)
		return "", err
	}
	return string(content), nil
}

func (c *Cache) Set(key string, value string, expire time.Duration) error {
	path := c.buildPath(key)
	if err := os.MkdirAll(filepath.Dir(path), c.dirMode); err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, c.fileMode)
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
	}()
	if _, err := f.WriteString(value); err != nil {
		return err
	}
	if expire <= 0 {
		expire = c.expire
	}
	return os.Chtimes(path, time.Now(), time.Now().Add(expire))
}

func (c *Cache) Load(key string, loader func() (string, error), expire time.Duration) (string, error) {
	if c.Has(key) {
		return c.Get(key)
	}
	value, err := loader()
	if err != nil {
		return "", err
	}
	if err := c.Set(key, value, expire); err != nil {
		return "", err
	}
	return value, nil
}

func (c *Cache) Delete(key string) error {
	if err := os.Remove(c.buildPath(key)); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return nil
}

func (c *Cache) Has(key string) bool {
	path := c.buildPath(key)
	if s, err := os.Stat(path); err != nil {
		return false
	} else if s.ModTime().Before(time.Now()) {
		_ = os.Remove(path)
		return false
	}
	return true
}

func (c *Cache) Clear() error {
	entries, err := os.ReadDir(c.storagePath)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !strings.HasPrefix(entry.Name(), c.prefix) {
			continue
		}
		err := os.Remove(filepath.Join(c.storagePath, entry.Name()))
		if err != nil {
			return err
		}
	}
	return nil
}
