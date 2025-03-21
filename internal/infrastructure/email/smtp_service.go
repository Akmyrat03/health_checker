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

	var msg strings.Builder
	msg.WriteString(fmt.Sprintf("From: %s\n", s.smtp.SMTPEmail))
	msg.WriteString(fmt.Sprintf("To: %s\n", strings.Join(receivers, ",")))
	msg.WriteString(fmt.Sprintf("Subject: %s\n", subjectMessage))
	msg.WriteString("MIME-Version: 1.0\n")
	msg.WriteString("Content-Type: text/html; charset=UTF-8\n\n") // <-- Set Content-Type to HTML
	msg.WriteString(message)                                      // <-- Send HTML content instead of plain text

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
