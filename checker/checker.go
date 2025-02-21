package checker

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// Track server states
var lastErrorTime = make(map[string]time.Time)
var lastNotifyTime = make(map[string]time.Time)
var isServerDown = make(map[string]bool)

var notificationInterval = 2 * time.Minute // Time before sending another notification

func CheckServerHealth(name string, url string, logFile string, timeout int) error {
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		handleError(name, url, logFile)
		return fmt.Errorf("ERROR: Server - %s (%s) is unreachable", name, url)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		handleError(name, url, logFile)
		return fmt.Errorf("ERROR: Server - %s (%s) returned status %d", name, url, resp.StatusCode)
	}

	// If the server was down but is now healthy, send recovery notification
	if isServerDown[name] {
		sendRecoveryNotification(name, url, logFile)
	}

	// Reset state
	isServerDown[name] = false
	delete(lastErrorTime, name)
	delete(lastNotifyTime, name)

	return nil
}

func handleError(name, url, logFile string) {
	errMsg := fmt.Sprintf("ERROR: Server - %s (%s) is unreachable or returned status 500", name, url)

	if !isServerDown[name] { // First time the server goes down
		writeLog(errMsg, logFile) // Log the first error
		SendEmailNotification(errMsg)
		lastErrorTime[name] = time.Now()
		lastNotifyTime[name] = time.Now()
		isServerDown[name] = true
	} else if time.Since(lastNotifyTime[name]) > notificationInterval { // Send every 2 hours
		writeLog(errMsg, logFile) // Log repeated errors every 2 hours
		SendEmailNotification(errMsg)
		lastNotifyTime[name] = time.Now()
	}
}

func sendRecoveryNotification(name, url, logFile string) {
	recoveryMsg := fmt.Sprintf("Server - %s (%s) is back online", name, url)
	writeLog(recoveryMsg, logFile)
	SendEmailNotification(recoveryMsg)
}

func writeLog(message string, logFile string) {
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("ERROR: Failed to open log file %s: %v", logFile, err)
		return
	}
	defer file.Close()

	logger := log.New(file, "", log.LstdFlags)
	logger.Println(message)
}
