package main

import (
	"log"

	. "github.com/redis-go/config"
	. "github.com/redis-go/internal/app"
)

func main() {
	// Configuration
	config, err := NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	Run(config)
}
