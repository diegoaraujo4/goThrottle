package main

import (
	"context"
	"goThrottle/config"
	"goThrottle/limiter"
	"goThrottle/middleware"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func responder(w http.ResponseWriter, _ *http.Request) {
	log.Print("Executing Response Handler")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("OK")); err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

func main() {
	ctx := context.Background()
	// Loading environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}
	log.Default().Printf("Configuration loaded: %+v\n", cfg)

	//Initialize Redis client
	redisClient := limiter.NewRedisClient(cfg.RedisAddress)
	log.Default().Printf("Redis client initialized...")
	// Create a new limiter
	rateLimiter, err := limiter.NewLimiter(redisClient, cfg)
	if err != nil {
		log.Fatalf("Error creating limiter: %v", err)
	}
	log.Default().Printf("Rate limiter initialized...")

	// Set up the HTTP server
	mux := http.NewServeMux()
	mux.Handle("/", middleware.RateLimiter(ctx, rateLimiter, http.HandlerFunc(responder)))

	// Start the server
	port := "8080"
	log.Default().Printf("Listening on port %s...\n", port)
	err = http.ListenAndServe(":"+port, mux)
	log.Fatal(err)
}
