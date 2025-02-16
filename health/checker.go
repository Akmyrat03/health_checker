package health

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func CheckServerHealth(url string, logFile string, timeout int) error {

	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		writeLog(fmt.Sprintf("ERROR: Server %s is unreachable or returned status %d", url, 404), logFile)
		return fmt.Errorf("ERROR: Server %s is unreachable or returned status 404", url)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		errMsg := fmt.Sprintf("ERROR: Server %s is unreachable or returned status %d", url, resp.StatusCode)
		writeLog(errMsg, logFile)
		return fmt.Errorf(errMsg)
	}
	return nil
}

func writeLog(message string, logFile string) {
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("os.OpenFile - Error: %v", err)
	}
	defer file.Close()

	logger := log.New(file, "", log.LstdFlags)
	logger.Println(message)
}
