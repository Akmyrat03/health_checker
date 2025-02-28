package smtp

import (
	"checker/internal/domain/app/repositories"
	"context"
	"fmt"
)

type SMTPUseCase struct {
	smtpRepository repositories.SMTP
}

func NewSMTPUseCase(smtpRepo repositories.SMTP) *SMTPUseCase {
	return &SMTPUseCase{smtpRepository: smtpRepo}
}

func (s *SMTPUseCase) SendEmail(ctx context.Context, message string) error {
	err := s.smtpRepository.SendEmail(ctx, message)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}
