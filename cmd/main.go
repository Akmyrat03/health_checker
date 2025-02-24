package main

import (
	"checker/internal/config"
	"checker/internal/domain/app"
	"checker/internal/infrastructure/fiber"
	"log"
)

func main() {
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("[config.LoadConfig]: failed to load config file: %v", err)
	}

	go app.TimeScheduler(cfg, 2)

	fiber.RunFiberServer(cfg)
}
