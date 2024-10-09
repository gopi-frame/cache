package memory

import (
	"github.com/gopi-frame/cache"
	cc "github.com/gopi-frame/contract/cache"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var (
	testCache cc.Cache
)

func TestMain(m *testing.M) {
	c, err := Open(map[string]any{
		"expire": time.Second * 2,
	})
	if err != nil {
		panic(err)
	}
	if err := c.Clear(); err != nil {
		panic(err)
	}
	testCache = c
	m.Run()
}
func TestCache_Set(t *testing.T) {
	t.Run("with duration", func(t *testing.T) {
		if err := testCache.Set("key", "value", time.Second); err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.True(t, testCache.Has("key"))
		time.Sleep(time.Second)
		assert.False(t, testCache.Has("key"))
	})

	t.Run("without duration", func(t *testing.T) {
		if err := testCache.Set("key", "value", 0); err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.True(t, testCache.Has("key"))
		time.Sleep(time.Second)
		assert.True(t, testCache.Has("key"))
		time.Sleep(time.Second)
		assert.False(t, testCache.Has("key"))
	})
}

func TestCache_Get(t *testing.T) {
	t.Run("file not exist", func(t *testing.T) {
		_, err := testCache.Get("key")
		assert.ErrorIs(t, err, cache.ErrCacheNotFound)
	})

	t.Run("cache out of date", func(t *testing.T) {
		if err := testCache.Set("key", "value", time.Second); err != nil {
			assert.FailNow(t, err.Error())
		}
		time.Sleep(time.Second)
		_, err := testCache.Get("key")
		assert.ErrorIs(t, err, cache.ErrCacheNotFound)
	})

	t.Run("cache exists", func(t *testing.T) {
		if err := testCache.Set("key", "value", 0); err != nil {
			assert.FailNow(t, err.Error())
		}
		if value, err := testCache.Get("key"); err != nil {
			assert.FailNow(t, err.Error())
		} else {
			assert.Equal(t, "value", value)
		}
	})
}

func TestCache_Load(t *testing.T) {
	t.Run("cache not exist", func(t *testing.T) {
		if value, err := testCache.Load("key", func() (string, error) {
			return "value", nil
		}, 0); err != nil {
			assert.FailNow(t, err.Error())
		} else {
			assert.Equal(t, "value", value)
		}
	})

	t.Run("cache out of date", func(t *testing.T) {
		if err := testCache.Set("key", "value", time.Second); err != nil {
			assert.FailNow(t, err.Error())
		}
		time.Sleep(time.Second)
		if value, err := testCache.Load("key", func() (string, error) {
			return "value", nil
		}, 0); err != nil {
			assert.FailNow(t, err.Error())
		} else {
			assert.Equal(t, "value", value)
		}
	})

	t.Run("cache exists", func(t *testing.T) {
		if err := testCache.Set("key", "value", 0); err != nil {
			assert.FailNow(t, err.Error())
		}
		if value, err := testCache.Load("key", func() (string, error) {
			return "value1", nil
		}, 0); err != nil {
			assert.FailNow(t, err.Error())
		} else {
			assert.Equal(t, "value", value)
		}
	})
}

func TestCache_Delete(t *testing.T) {
	t.Run("file not exist", func(t *testing.T) {
		if err := testCache.Delete("key"); err != nil {
			assert.FailNow(t, err.Error())
		}
	})

	t.Run("file exists", func(t *testing.T) {
		if err := testCache.Set("key", "value", 0); err != nil {
			assert.FailNow(t, err.Error())
		}
		if err := testCache.Delete("key"); err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.False(t, testCache.Has("key"))
	})
}

func TestCache_Clear(t *testing.T) {
	if err := testCache.Set("key", "value", 0); err != nil {
		assert.FailNow(t, err.Error())
	}
	if err := testCache.Set("key2", "value", 0); err != nil {
		assert.FailNow(t, err.Error())
	}
	if err := testCache.Clear(); err != nil {
		assert.FailNow(t, err.Error())
	}
	assert.False(t, testCache.Has("key"))
	assert.False(t, testCache.Has("key2"))
}
