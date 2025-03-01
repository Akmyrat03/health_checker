package services

import "context"

type SMTP interface {
	SendEmail(ctx context.Context, message string, receivers []string) error
}
