package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	IPLimit       int
	TokenLimit    int
	BlockDuration int
	RedisAddress  string
}

func LoadConfig() Config {
	ipLimitStr := getEnv("IP_LIMIT", "5")
	tokenLimitStr := getEnv("TOKEN_LIMIT", "10")
	redisAddress := getEnv("REDIS_ADDRESS", "localhost:6379")
	blockDurationStr := getEnv("BLOCK_DURATION", "300")

	ipLimit, err := strconv.Atoi(ipLimitStr)
	if err != nil {
		log.Fatalf("Invalid IP_LIMIT: %v", err)
	}

	tokenLimit, err := strconv.Atoi(tokenLimitStr)
	if err != nil {
		log.Fatalf("Invalid TOKEN_LIMIT '%s': %v", tokenLimitStr, err)
	}

	blockDuration, err := strconv.Atoi(blockDurationStr)
	if err != nil {
		log.Fatalf("Invalid BLOCK_DURATION '%s': %v", blockDurationStr, err)
	}

	return Config{
		IPLimit:       ipLimit,
		TokenLimit:    tokenLimit,
		BlockDuration: blockDuration,
		RedisAddress:  redisAddress,
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" && os.Getenv(key) == "" {
		return defaultValue
	}
	return value
}
