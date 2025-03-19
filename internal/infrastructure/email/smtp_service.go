package email

import (
	"checker/internal/config"
	"context"
	"fmt"
	"net/smtp"
	"strings"
)

type SMTPService struct {
	smtp config.SMTP
}

func NewSMTPService(cfg config.SMTP) *SMTPService {
	return &SMTPService{smtp: cfg}
}

func (s *SMTPService) SendEmail(ctx context.Context, message, subjectMessage string, receivers []string) error {
	auth := smtp.PlainAuth("", s.smtp.SMTPEmail, s.smtp.SMTPPass, s.smtp.SMTPServer)

	subject := subjectMessage

	body := message

	var msg strings.Builder
	msg.WriteString(fmt.Sprintf("From: %s\n", s.smtp.SMTPEmail))
	msg.WriteString(fmt.Sprintf("To: %s\n", strings.Join(receivers, ",")))
	msg.WriteString(fmt.Sprintf("Subject: %s\n\n", subject))
	msg.WriteString(body)

	emailMessage := msg.String()

	err := smtp.SendMail(
		fmt.Sprintf("%s:%d", s.smtp.SMTPServer, s.smtp.SMTPPort),
		auth,
		s.smtp.SMTPEmail,
		receivers,
		[]byte(emailMessage),
	)
	if err != nil {
		return fmt.Errorf("failed to send email notification: %w", err)
	}

	return nil
}
