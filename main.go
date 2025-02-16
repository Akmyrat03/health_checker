package main

import (
	"checker/config"
	"checker/handler"
	"checker/health"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	config, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("config.LoadConfig - Error: %v", err)
	}

	fmt.Println("Servers:", config.Servers)
	fmt.Println("Check Interval:", config.CheckInterval)

	ticker := time.NewTicker(time.Duration(config.CheckInterval) * time.Second)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		handler.HealthHandler(w, r, config)
	})

	go http.ListenAndServe(":8080", nil)

	for {
		select {
		case <-ticker.C:
			health.StartWorkers(config.Servers, 2, config.CheckInterval, config.LogFile)
		}
	}
}
