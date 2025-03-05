package usecases

import (
	"checker/internal/domain/app/inputs"
	"checker/internal/domain/app/repositories"
	"checker/internal/domain/app/services"
	"checker/internal/domain/entities"
	"context"
	"fmt"
)

type ReceiversUseCase struct {
	receiversRepository repositories.Receivers
	smtpService         services.SMTP
}

func NewReceiversUseCase(receiverRepository repositories.Receivers, smtpService services.SMTP) *ReceiversUseCase {
	return &ReceiversUseCase{receiversRepository: receiverRepository, smtpService: smtpService}
}

func (receiverUseCase *ReceiversUseCase) Create(ctx context.Context, receiver inputs.CreateReceiver) (int, error) {
	id, err := receiverUseCase.receiversRepository.Create(ctx, receiver)
	if err != nil {
		return 0, fmt.Errorf("failed to create receiver: %v", err)
	}

	return id, nil
}

func (receiverUseCase *ReceiversUseCase) Delete(ctx context.Context, id int) error {
	err := receiverUseCase.receiversRepository.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (receiverUseCase *ReceiversUseCase) List(ctx context.Context) ([]entities.Receiver, error) {
	receivers, err := receiverUseCase.receiversRepository.List(ctx)
	if err != nil {
		return nil, err
	}

	return receivers, nil
}

func (receiverUseCase *ReceiversUseCase) SendEmailToReceiver(ctx context.Context, message string) error {
	receivers, err := receiverUseCase.receiversRepository.List(ctx)
	if err != nil {
		return fmt.Errorf("failed to retrieve receivers: %v", err)
	}

	if len(receivers) == 0 {
		return fmt.Errorf("no receivers found")
	}

	var emails []string
	for _, receiver := range receivers {
		emails = append(emails, receiver.Email)
	}

	err = receiverUseCase.smtpService.SendEmail(ctx, message, emails)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}
