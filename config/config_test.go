package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	os.Clearenv()
	os.Setenv("IP_LIMIT", "10")
	os.Setenv("TOKEN_LIMIT", "20")
	os.Setenv("BLOCK_DURATION", "600")
	os.Setenv("REDIS_ADDRESS", "127.0.0.1:6379")

	config, err := LoadConfig()
	assert.Nil(t, err, "Error should be nil")

	assert.Equal(t, 10, config.IPLimit)
	assert.Equal(t, 20, config.TokenLimit)
	assert.Equal(t, 600, config.BlockDuration)
	assert.Equal(t, "127.0.0.1:6379", config.RedisAddress)
}

func TestLoadInvalidIPConfig(t *testing.T) {
	os.Clearenv()
	os.Setenv("IP_LIMIT", "A")
	config, err := LoadConfig()
	assert.Error(t, err)
	assert.ErrorContains(t, err, "invalid IP_LIMIT")
	assert.Equal(t, config, Config{})
}

func TestLoadInvalidTokenLimitConfig(t *testing.T) {
	os.Clearenv()
	os.Setenv("IP_LIMIT", "10")
	os.Setenv("TOKEN_LIMIT", "A")
	config, err := LoadConfig()
	assert.Error(t, err)
	assert.ErrorContains(t, err, "invalid TOKEN_LIMIT ")
	assert.Equal(t, config, Config{})
}

func TestLoadInvalidDurationConfig(t *testing.T) {
	os.Clearenv()
	os.Setenv("IP_LIMIT", "10")
	os.Setenv("TOKEN_LIMIT", "20")
	os.Setenv("BLOCK_DURATION", "A")
	config, err := LoadConfig()
	assert.Error(t, err)
	assert.ErrorContains(t, err, "invalid BLOCK_DURATION ")
	assert.Equal(t, config, Config{})
}

func TestLoadConfigWithDefaults(t *testing.T) {
	os.Clearenv()

	config, err := LoadConfig()
	assert.Nil(t, err, "Error should be nil")

	assert.Equal(t, 5, config.IPLimit)
	assert.Equal(t, 10, config.TokenLimit)
	assert.Equal(t, 300, config.BlockDuration)
	assert.Equal(t, "localhost:6379", config.RedisAddress)
}
