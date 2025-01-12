package limiter

import (
	"context"
	"fmt"
	"goThrottle/config"
	"time"

	"github.com/redis/go-redis/v9"
)

type Limiter struct {
	client *redis.Client
	config config.Config
}

func NewLimiter(client *redis.Client, config config.Config) *Limiter {
	return &Limiter{
		client: client,
		config: config,
	}
}

func (l *Limiter) CheckLimit(key string, limit int) (bool, error) {
	ctx := context.Background()

	reqCount, err := l.client.Incr(ctx, key).Result()
	if err != nil {
		return false, err
	}

	duration := time.Duration(l.config.BlockDuration) * time.Second

	if reqCount == 1 {
		l.client.Expire(ctx, key, duration)
	}

	if reqCount > int64(limit) {
		l.client.Set(ctx, fmt.Sprintf("%s:block", key), "blocked", duration)
		return false, nil
	}

	return true, nil
}
