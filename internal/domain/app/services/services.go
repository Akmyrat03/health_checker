package services

type MessageSender interface {
	SendEmail(message string) error
}
