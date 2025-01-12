package limiter

import (
	"goThrottle/config"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
)

func TestLimiter_CheckLimit(t *testing.T) {
	client, mock := redismock.NewClientMock()
	config := config.Config{
		IPLimit:       5,
		TokenLimit:    10,
		BlockDuration: 1,
	}
	limiter := NewLimiter(client, config)

	duration := time.Duration(config.BlockDuration) * time.Second

	t.Run("within limit", func(t *testing.T) {
		key := "test_key_within_limit"
		for i := 0; i < config.IPLimit; i++ {
			mock.ExpectIncr(key).SetVal(int64(i + 1))
			if i == 0 {
				mock.ExpectExpire(key, duration).SetVal(true)
			}

			allowed, err := limiter.CheckLimit(key, config.IPLimit)
			assert.NoError(t, err)
			assert.True(t, allowed)
		}
	})

	t.Run("exceed limit", func(t *testing.T) {
		key := "test_key_exceed_limit"
		for i := 0; i <= config.IPLimit; i++ {
			mock.ExpectIncr(key).SetVal(int64(i + 1))
			if i == 0 {
				mock.ExpectExpire(key, duration).SetVal(true)
			}

			allowed, err := limiter.CheckLimit(key, config.IPLimit)

			if i < config.IPLimit {
				assert.NoError(t, err)
				assert.True(t, allowed)
			}

			if i == config.IPLimit {
				assert.NoError(t, err)
				assert.False(t, allowed)
			}
		}
	})

	t.Run("block duration", func(t *testing.T) {
		key := "test_key_block_duration"
		requestCount := 0
		for requestCount < config.IPLimit {
			mock.ExpectIncr(key).SetVal(int64(requestCount + 1))
			if requestCount == 0 {
				mock.ExpectExpire(key, duration).SetVal(true)
			}

			allowed, err := limiter.CheckLimit(key, config.IPLimit)
			assert.NoError(t, err)
			assert.True(t, allowed)

			requestCount++
		}

		// Invalid blocked request
		mock.ExpectIncr(key).SetVal(int64(requestCount + 1))
		allowed, err := limiter.CheckLimit(key, config.IPLimit)
		assert.NoError(t, err)
		assert.False(t, allowed)

		// Wait for block duration
		time.Sleep(duration)

		mock.ExpectIncr(key).SetVal(int64(1))
		mock.ExpectExpire(key, duration).SetVal(true)
		allowed, err = limiter.CheckLimit(key, config.IPLimit)
		assert.NoError(t, err)
		assert.True(t, allowed)
	})
}
