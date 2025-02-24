package usecases

import (
	"checker/internal/domain/app/services"
	"log"
)

func SendEmailNotification(sender services.MessageSender, message string) {
	err := sender.SendEmail(message)
	if err != nil {
		log.Printf("Failed to send email notification: %v", err)
	}
}
