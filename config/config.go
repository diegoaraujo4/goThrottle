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
	ipLimit, err := strconv.Atoi(getEnv("IP_LIMIT", "5"))
	if err != nil {
		log.Fatalf("Invalid IP_LIMIT: %v", err)
	}

	tokenLimit, err := strconv.Atoi(getEnv("TOKEN_LIMIT", "10"))
	if err != nil {
		log.Fatalf("Invalid TOKEN_LIMIT: %v", err)
	}

	blockDuration, err := strconv.Atoi(getEnv("BLOCK_IN_SECONDS", "300"))
	if err != nil {
		log.Fatalf("Invalid BLOCK_IN_SECONDS: %v", err)
	}

	return Config{
		IPLimit:       ipLimit,
		TokenLimit:    tokenLimit,
		BlockDuration: blockDuration,
		RedisAddress:  getEnv("REDIS_ADDRESS", "localhost:6379"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
