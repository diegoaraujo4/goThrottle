package main

import (
	"goThrottle/config"
	"log"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()
	log.Default().Printf("Config: %+v", cfg)
}
