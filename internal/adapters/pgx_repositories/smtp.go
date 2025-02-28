package pgx_repositories

import (
	"checker/internal/config"
	"context"
	"fmt"
	"net/smtp"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SMTPRepository struct {
	DbPool *pgxpool.Pool
	smtp   config.SMTP
}

func NewSMTPRepository(dbPool *pgxpool.Pool, filename string) (*SMTPRepository, error) {
	cfg, err := config.LoadConfig(filename)
	if err != nil {
		return nil, err
	}
	return &SMTPRepository{DbPool: dbPool, smtp: cfg.SMTP}, nil
}

func (s *SMTPRepository) FetchReceivers(ctx context.Context) ([]string, error) {
	var receivers []string

	query := `SELECT email FROM receivers`

	rows, err := s.DbPool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch receivers: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var email string
		if err := rows.Scan(&email); err != nil {
			return nil, fmt.Errorf("failed to scan receiver email: %v", err)
		}
		receivers = append(receivers, email)
	}

	if len(receivers) == 0 {
		return nil, fmt.Errorf("no receivers found in database")
	}

	return receivers, nil
}

func (s *SMTPRepository) SendEmail(ctx context.Context, message string) error {
	auth := smtp.PlainAuth("", s.smtp.SMTPEmail, s.smtp.SMTPPass, s.smtp.SMTPServer)

	receivers, err := s.FetchReceivers(ctx)
	if err != nil {
		return fmt.Errorf("failed to retrieve receivers: %v", err)
	}

	subject := s.smtp.SubjectPrefix + " - Error Notification"
	body := message

	var msg strings.Builder
	msg.WriteString(fmt.Sprintf("From: %s\n", s.smtp.SMTPEmail))
	msg.WriteString(fmt.Sprintf("To: %s\n", strings.Join(receivers, ",")))
	msg.WriteString(fmt.Sprintf("Subject: %s\n\n", subject))
	msg.WriteString(body)

	emailMessage := msg.String()

	err = smtp.SendMail(
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
