package main

import (
	_ "checker/docs"
	"checker/internal/config"
	"checker/internal/infrastructure/fiber"
)

// @title Health Checker API
// @version 1.0
// @description API for managing servers
// @schemas http, https
// @in header
func main() {
	cfg := config.LoadConfig()
	fiber.RunFiberServer(cfg)
}
