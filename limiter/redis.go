package limiter

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// Define an interface for the Redis client
// mockgen -source=limiter/redis.go -destination=limiter/mock_redis_client.go -package=limiter
type RedisClientInterface interface {
	Incr(ctx context.Context, key string) *redis.IntCmd
	Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Get(ctx context.Context, key string) *redis.StringCmd
}

func NewRedisClient(address string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: address,
	})

	return client
}
