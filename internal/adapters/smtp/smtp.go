package smtp

import (
	"checker/internal/config"
	"fmt"
	"net/smtp"
	"strings"
)

type SMTP struct {
	smtp config.SMTP
}

func NewSMTP(filename string) (*SMTP, error) {
	cfg, err := config.LoadConfig(filename)
	if err != nil {
		return nil, err
	}
	return &SMTP{smtp: cfg.SMTP}, nil
}

func (s *SMTP) SendEmail(message string) error {
	auth := smtp.PlainAuth("", s.smtp.SMTPEmail, s.smtp.SMTPPass, s.smtp.SMTPServer)

	usersTo := s.smtp.Receivers
	subject := s.smtp.SubjectPrefix + " - Error Notification"
	body := message

	var msg strings.Builder
	msg.WriteString(fmt.Sprintf("From: %s\n", s.smtp.SMTPEmail))
	msg.WriteString(fmt.Sprintf("To: %s\n", strings.Join(usersTo, ",")))
	msg.WriteString(fmt.Sprintf("Subject: %s\n\n", subject))
	msg.WriteString(body)

	emailMessage := msg.String()

	err := smtp.SendMail(
		fmt.Sprintf("%s:%d", s.smtp.SMTPServer, s.smtp.SMTPPort),
		auth,
		s.smtp.SMTPEmail,
		usersTo,
		[]byte(emailMessage),
	)

	if err != nil {
		return fmt.Errorf("failed to send email notification: %w", err)
	}

	return nil
}
