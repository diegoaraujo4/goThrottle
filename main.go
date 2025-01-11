package main

import (
	"goThrottle/config"
	"goThrottle/limiter"
	"log"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize Redis client
	redisClient := limiter.NewRedisClient(cfg.RedisAddress)
	log.Default().Printf("Config: %+v", redisClient.Options().Addr)
}
