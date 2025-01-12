package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"goThrottle/config"
	"goThrottle/limiter"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
)

func TestRateLimiter(t *testing.T) {
	client, mockRedis := redismock.NewClientMock()
	tests := []struct {
		name           string
		redisKey       string
		apiKey         string
		remoteAddr     string
		iplimit        int
		tokenLimit     int
		expectedStatus int
		tokenValidaton bool
		skipRedis      bool
	}{
		{
			name:           "Valid API token limit",
			tokenLimit:     5,
			redisKey:       "token:valid-api-key",
			apiKey:         "valid-api-key",
			remoteAddr:     "192.168.1.1:1234",
			tokenValidaton: true,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Exceeded API token limit",
			tokenLimit:     0,
			redisKey:       "token:invalid-api-key",
			apiKey:         "invalid-api-key",
			remoteAddr:     "192.168.1.1:1234",
			tokenValidaton: true,
			expectedStatus: http.StatusTooManyRequests,
		},
		{
			name:           "Valid IP address limit",
			iplimit:        5,
			redisKey:       "ip:192.168.1.1",
			remoteAddr:     "192.168.1.1:1234",
			tokenValidaton: false,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Exceeded IP address limit",
			iplimit:        0,
			redisKey:       "ip:192.168.1.1",
			remoteAddr:     "192.168.1.1:1234",
			tokenValidaton: false,
			expectedStatus: http.StatusTooManyRequests,
		},
		{
			name:           "Invalid remote address",
			remoteAddr:     "invalid-remote-addr",
			tokenValidaton: false,
			expectedStatus: http.StatusBadRequest,
			skipRedis:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new rate limiter with the mock client
			rateLimiter, err := limiter.NewLimiter(client, config.Config{
				IPLimit:       tt.iplimit,
				TokenLimit:    tt.tokenLimit,
				BlockDuration: 1,
			})
			assert.NoError(t, err)
			handler := RateLimiter(context.Background(), rateLimiter, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))
			duration := time.Duration(1) * time.Second

			if !tt.skipRedis {
				mockRedis.ExpectIncr(tt.redisKey).SetVal(int64(1))
				mockRedis.ExpectExpire(tt.redisKey, duration).SetVal(true)
			}

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.RemoteAddr = tt.remoteAddr

			if tt.tokenValidaton {
				req.Header.Set("API_KEY", tt.apiKey)
			}
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.NoError(t, mockRedis.ExpectationsWereMet())
		})
	}
}
