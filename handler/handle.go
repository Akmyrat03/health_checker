package handler

import (
	health "checker/checker"
	"checker/config"
	"encoding/json"
	"net/http"
)

type HealthStatus struct {
	ServerURL string `json:"server_url"`
	Status    string `json:"status"`
}

func HealthHandler(w http.ResponseWriter, r *http.Request, cfg *config.Config) {
	servers, err := config.LoadServers("servers.json")
	if err != nil {
		http.Error(w, "failed to load servers", http.StatusInternalServerError)
		return
	}

	status := make([]HealthStatus, 0)

	for _, server := range servers {
		err := health.CheckServerHealth(server, cfg.LogFile, cfg.Timeout)
		status = append(status, HealthStatus{
			ServerURL: server,
			Status:    "Healthy",
		})

		if err != nil {
			status[len(status)-1].Status = "Unhealthy"
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}
