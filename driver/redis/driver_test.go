package redis

import (
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestOpenT(t *testing.T) {
	t.Run("without client", func(t *testing.T) {
		_, err := OpenT[any](map[string]any{
			"expire": time.Second * 2,
		})
		assert.Error(t, err)
	})

	t.Run("with client", func(t *testing.T) {
		_, err := OpenT[any](map[string]any{
			"expire": time.Second * 2,
			"client": redis.NewClient(&redis.Options{
				Addr:     "localhost:6379",
				Password: "", // no password set
				DB:       0,  // use default DB
			}),
		})
		assert.Nil(t, err)
	})
}
