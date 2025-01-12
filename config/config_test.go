package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	os.Setenv("IP_LIMIT", "10")
	os.Setenv("TOKEN_LIMIT", "20")
	os.Setenv("BLOCK_DURATION", "600")
	os.Setenv("REDIS_ADDRESS", "127.0.0.1:6379")

	config := LoadConfig()

	assert.Equal(t, 10, config.IPLimit, "IPLimit should be 10")
	assert.Equal(t, 20, config.TokenLimit, "TokenLimit should be 20")
	assert.Equal(t, 600, config.BlockDuration, "BlockDuration should be 600")
	assert.Equal(t, "127.0.0.1:6379", config.RedisAddress, "RedisAddress should be 127.0.0.1:6379")
}

func TestLoadConfigWithDefaults(t *testing.T) {
	os.Clearenv()

	config := LoadConfig()

	assert.Equal(t, 5, config.IPLimit, "IPLimit should be 5")
	assert.Equal(t, 10, config.TokenLimit, "TokenLimit should be 10")
	assert.Equal(t, 300, config.BlockDuration, "BlockDuration should be 300")
	assert.Equal(t, "localhost:6379", config.RedisAddress, "RedisAddress should be localhost:6379")
}
