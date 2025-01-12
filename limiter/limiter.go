package limiter

import (
	"context"
	"fmt"
	"goThrottle/config"
	"log"
	"time"
)

const (
	IPLimit    = "ipLimit"
	TokenLimit = "tokenLimit"
)

// Limiter struct will implement the rate limiting logic
type Limiter struct {
	client RedisClientInterface
	config LimiterConfig
}

// LimiterConfig contains the configuration values for the rate limiter
type LimiterConfig struct {
	IPLimit       int
	TokenLimit    int
	BlockDuration int
}

func NewLimiter(client RedisClientInterface, config config.Config) (*Limiter, error) {
	if config.IPLimit < 0 || config.TokenLimit < 0 || config.BlockDuration < 0 {
		return nil, fmt.Errorf("invalid configuration values")
	}

	return &Limiter{
		client: client,
		config: LimiterConfig{
			IPLimit:       config.IPLimit,
			TokenLimit:    config.TokenLimit,
			BlockDuration: config.BlockDuration,
		},
	}, nil
}

func (l *Limiter) CheckLimit(ctx context.Context, key string, limitType string) (bool, error) {
	var limit int
	duration := time.Duration(l.config.BlockDuration) * time.Second

	switch limitType {
	case IPLimit:
		limit = l.config.IPLimit
	case TokenLimit:
		limit = l.config.TokenLimit
	default:
		return false, fmt.Errorf("unknown limit type: %s", limitType)
	}

	reqCount, err := l.client.Incr(ctx, key).Result()
	if err != nil {
		log.Default().Printf("Error incrementing key %s: %v", key, err)
		return false, err
	}

	if reqCount == 1 {
		if err := l.client.Expire(ctx, key, duration).Err(); err != nil {
			log.Default().Printf("Error setting expiration for key %s: %v", key, err)
			return false, err
		}
	}

	if reqCount > int64(limit) {
		l.client.Set(ctx, fmt.Sprintf("%s:block", key), "blocked", duration)
		log.Default().Printf("Rate limit exceeded for %s. Blocking for %d seconds.", key, l.config.BlockDuration)
		return false, nil
	}

	return true, nil
}
