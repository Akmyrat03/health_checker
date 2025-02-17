package main

import (
	config "checker/config"
	"checker/handler"
	"checker/scheduler"
	"fmt"
	"log"
	"net/http"
)

func main() {
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("config.LoadConfig - Error: %v", err)
	}

	servers, err := config.LoadServers("servers.json")
	if err != nil {
		log.Fatalf("config.LoadServers - Error: %v", err)
	}

	fmt.Printf("Servers: %+v\n", servers)
	fmt.Println("Check Interval:", cfg.CheckInterval)
	fmt.Println("Timeout:", cfg.Timeout)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		handler.HealthHandler(w, r, cfg)
	})

	go http.ListenAndServe(":8080", nil)

	scheduler.StartHealthCheckScheduler(servers, cfg.CheckInterval, cfg.LogFile, cfg.Timeout)
}
