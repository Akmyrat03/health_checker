package services

import "context"

type SMTP interface {
	SendEmail(ctx context.Context, message string, subjectMessage string, receivers []string) error
}
