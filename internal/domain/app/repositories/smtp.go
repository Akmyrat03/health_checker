package repositories

import "context"

type SMTP interface {
	SendEmail(ctx context.Context, message string) error
}
