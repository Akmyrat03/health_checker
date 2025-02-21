package handler

import (
	"checker/checker"
	"checker/config"
	"encoding/json"
	"net/http"
	"sync"
)

type HealthStatus struct {
	ServerName string `json:"server_name"`
	ServerURL  string `json:"server_url"`
	Status     string `json:"status"`
}

func HealthHandler(w http.ResponseWriter, r *http.Request, cfg *config.Config) {
	serversChan, err := config.LoadServers("config.json")
	if err != nil {
		http.Error(w, "failed to load servers", http.StatusInternalServerError)
		return
	}

	var wg sync.WaitGroup
	var status []HealthStatus

	for server := range serversChan {
		wg.Add(1)

		go func(server config.Server) {
			defer wg.Done()

			err := checker.CheckServerHealth(server.Name, server.Server, cfg.HealthChecker.LogFile, cfg.HealthChecker.Timeout)
			serverStatus := "Healthy"
			if err != nil {
				serverStatus = "Unhealthy"
			}

			status = append(status, HealthStatus{
				ServerName: server.Name,
				ServerURL:  server.Server,
				Status:     serverStatus,
			})
		}(server)
	}

	wg.Wait()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}
