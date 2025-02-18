package scheduler

import (
	health "checker/checker"
	"log"
	"time"
)

func StartHealthCheckScheduler(filename string, checkInterval int, logFile string, timeout int) {
	ticker := time.NewTicker(time.Duration(checkInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			go func() {
				if err := health.StartWorkers(filename, 3, checkInterval, logFile, timeout); err != nil {
					log.Printf("Health check failed: %v", err)
				}
			}()
		}
	}
}
