package main

import (
	_ "checker/docs"
	"checker/internal/config"
	"checker/internal/infrastructure/fiber"
	"log"
)

// @title Health Checker API
// @version 1.0
// @description API for managing servers
// @schemas http, https
// @securityDefinitions.apikey BearerAuth
// @in header
// @name authorization
func main() {
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("[config.LoadConfig]: failed to load config file: %v", err)
	}

	fiber.RunFiberServer(cfg)
}
