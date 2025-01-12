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

	// Create a new limiter
	rateLimiter := limiter.NewLimiter(redisClient, cfg)
	log.Default().Println("Rate limiter initialized %d", rateLimiter)
}
