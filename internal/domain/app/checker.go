package app

import (
	"checker/internal/infrastructure/smtp"
	"checker/internal/shared"
	"fmt"
	"net/http"
	"time"
)

var (
	lastErrorTime  = make(map[string]time.Time)
	lastNotifyTime = make(map[string]time.Time)
	isServerDown   = make(map[string]bool)
)

var notificationInterval = 2 * time.Minute // Time before sending another notification

func CheckServer(name string, url string, logFile string, timeout int) error {
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		HandleError(name, url, logFile)
		return fmt.Errorf("ERROR: Server - %s (%s) is unreachable", name, url)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		HandleError(name, url, logFile)
		return fmt.Errorf("ERROR: Server - %s (%s) returned status %d", name, url, resp.StatusCode)
	}

	// If the server was down but is now healthy, send recovery notification
	if isServerDown[name] {
		SendRecoveryNotification(name, url, logFile)
	}

	// Reset state
	isServerDown[name] = false
	delete(lastErrorTime, name)
	delete(lastNotifyTime, name)

	return nil
}

func HandleError(name, url, logFile string) {
	errMsg := fmt.Sprintf("ERROR: Server - %s (%s) is unreachable or returned status 500", name, url)

	if !isServerDown[name] { // First time the server goes down

		shared.WriteLog(errMsg, logFile) // Log the first error
		smtp.SendEmailNotification(errMsg)
		lastErrorTime[name] = time.Now()
		lastNotifyTime[name] = time.Now()
		isServerDown[name] = true

	} else if time.Since(lastNotifyTime[name]) > notificationInterval {

		shared.WriteLog(errMsg, logFile) // Log repeated errors every 2 hours
		smtp.SendEmailNotification(errMsg)
		lastNotifyTime[name] = time.Now()

	}
}

func SendRecoveryNotification(name, url, logFile string) {
	onlineMsg := fmt.Sprintf("Server - %s (%s) is back online", name, url)
	shared.WriteLog(onlineMsg, logFile)
	smtp.SendEmailNotification(onlineMsg)
}
