package database

import (
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"testing"
	"time"
)

func TestOpenT(t *testing.T) {
	t.Run("without db", func(t *testing.T) {
		_, err := OpenT[any](map[string]any{
			"expire": time.Second * 2,
		})
		assert.Error(t, err)
	})

	t.Run("with db", func(t *testing.T) {
		db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		_, err = OpenT[any](map[string]any{
			"expire": time.Second * 2,
			"db":     db,
		})
		assert.Nil(t, err)
	})
}
