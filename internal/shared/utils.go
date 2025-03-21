package shared

import (
	"fmt"
	"log"
	"os"
	"time"
)

func WriteLog(message string, logFile string) {
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("ERROR: Failed to open log file %s: %v", logFile, err)
		return
	}
	defer file.Close()

	logger := log.New(file, "", log.LstdFlags)
	logger.Println(message)
}

func FormatDuration(d time.Duration) string {
	h := d / time.Hour
	m := (d % time.Hour) / time.Minute

	if h > 0 && m > 0 {
		return fmt.Sprintf("%d hours %d minutes", h, m)
	} else if h > 0 {
		return fmt.Sprintf("%d hours", h)
	}

	return fmt.Sprintf("%d minutes", m)
}

// GetOrdinalSuffix returns the ordinal suffix for a number (e.g., 1st, 2nd, 3rd, 4th, etc.)
func GetOrdinalSuffix(n int) string {
	if n%100 >= 11 && n%100 <= 13 {
		return "th"
	}
	switch n % 10 {
	case 1:
		return "st"
	case 2:
		return "nd"
	case 3:
		return "rd"
	default:
		return "th"
	}
}
