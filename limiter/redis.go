package limiter

import (
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(address string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: address,
	})

	return client
}
