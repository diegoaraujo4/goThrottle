package limiter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRedisClient(t *testing.T) {

	address := "localhost:6379"
	client := NewRedisClient(address)
	assert.NotNil(t, client)
}
