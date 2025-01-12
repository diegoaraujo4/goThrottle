package limiter

import (
	"context"
	"errors"
	"goThrottle/config"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
)

func TestNewLimiter_Error(t *testing.T) {
	client, _ := redismock.NewClientMock()
	testConfig := config.Config{
		IPLimit:       -1,
		TokenLimit:    -1,
		BlockDuration: -1,
	}
	_, err := NewLimiter(client, testConfig)
	assert.ErrorContains(t, err, "invalid configuration values")
}

func TestLimiter_CheckLimit(t *testing.T) {
	client, mock := redismock.NewClientMock()
	testConfig := config.Config{
		IPLimit:       5,
		TokenLimit:    10,
		BlockDuration: 1,
	}
	limiter, err := NewLimiter(client, testConfig)
	assert.NoError(t, err)

	duration := time.Duration(limiter.config.BlockDuration) * time.Second

	t.Run("within IP limit", func(t *testing.T) {
		ctx := context.Background()
		key := "test_key_within_limit"
		for i := 0; i < limiter.config.IPLimit; i++ {
			mock.ExpectIncr(key).SetVal(int64(i + 1))
			if i == 0 {
				mock.ExpectExpire(key, duration).SetVal(true)
			}

			allowed, err := limiter.CheckLimit(ctx, key, IPLimit)
			assert.NoError(t, err)
			assert.True(t, allowed)
		}
	})

	t.Run("within Token limit", func(t *testing.T) {
		ctx := context.Background()
		key := "test_key_within_limit"
		for i := 0; i < limiter.config.TokenLimit; i++ {
			mock.ExpectIncr(key).SetVal(int64(i + 1))
			if i == 0 {
				mock.ExpectExpire(key, duration).SetVal(true)
			}

			allowed, err := limiter.CheckLimit(ctx, key, TokenLimit)
			assert.NoError(t, err)
			assert.True(t, allowed)
		}
	})

	t.Run("within invalid limit", func(t *testing.T) {
		ctx := context.Background()
		key := "test_key_within_limit"
		allowed, err := limiter.CheckLimit(ctx, key, "someLimit	")
		assert.False(t, allowed)
		assert.ErrorContains(t, err, "unknown limit type: someLimit")
	})

	t.Run("exceed limit", func(t *testing.T) {
		ctx := context.Background()
		key := "test_key_exceed_limit"
		for i := 0; i <= limiter.config.IPLimit; i++ {
			mock.ExpectIncr(key).SetVal(int64(i + 1))
			if i == 0 {
				mock.ExpectExpire(key, duration).SetVal(true)
			}

			allowed, err := limiter.CheckLimit(ctx, key, IPLimit)

			if i < limiter.config.IPLimit {
				assert.NoError(t, err)
				assert.True(t, allowed)
			}

			if i == limiter.config.IPLimit {
				assert.NoError(t, err)
				assert.False(t, allowed)
			}
		}
	})

	t.Run("block duration", func(t *testing.T) {
		ctx := context.Background()
		key := "test_key_block_duration"
		requestCount := 0
		for requestCount < limiter.config.IPLimit {
			mock.ExpectIncr(key).SetVal(int64(requestCount + 1))
			if requestCount == 0 {
				mock.ExpectExpire(key, duration).SetVal(true)
			}

			allowed, err := limiter.CheckLimit(ctx, key, IPLimit)
			assert.NoError(t, err)
			assert.True(t, allowed)

			requestCount++
		}

		// Invalid blocked request
		mock.ExpectIncr(key).SetVal(int64(requestCount + 1))
		allowed, err := limiter.CheckLimit(ctx, key, IPLimit)
		assert.NoError(t, err)
		assert.False(t, allowed)

		// Wait for block duration
		time.Sleep(duration)

		mock.ExpectIncr(key).SetVal(int64(1))
		mock.ExpectExpire(key, duration).SetVal(true)
		allowed, err = limiter.CheckLimit(ctx, key, IPLimit)
		assert.NoError(t, err)
		assert.True(t, allowed)
	})

	t.Run("within invalid limit", func(t *testing.T) {
		ctx := context.Background()
		key := "test_key_within_limit"
		allowed, err := limiter.CheckLimit(ctx, key, "someLimit	")
		assert.False(t, allowed)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "unknown limit type: someLimit")
	})

	t.Run("Error on Increasing limit", func(t *testing.T) {
		ctx := context.Background()
		key := "test_key_within_limit_error"

		mock.ExpectIncr(key).SetErr(errors.New("error increasing limit"))
		allowed, err := limiter.CheckLimit(ctx, key, IPLimit)
		assert.False(t, allowed)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "error increasing limit")
	})

	t.Run("Error on adding token expiration", func(t *testing.T) {
		ctx := context.Background()
		key := "test_key_within_expiration_error"

		mock.ExpectIncr(key).SetVal(int64(1))
		mock.ExpectExpire(key, duration).SetErr(errors.New("error setting expiration"))
		allowed, err := limiter.CheckLimit(ctx, key, IPLimit)
		assert.False(t, allowed)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "error setting expiration")
	})

}
