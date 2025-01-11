package limiter

import (
	"github.com/go-redis/redis/v8"
)

func NewRedisClient(address string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: address,
	})

	return client
}
