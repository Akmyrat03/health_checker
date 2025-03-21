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

func (receiverUseCase *ReceiversUseCase) MuteStatus(ctx context.Context, email string, mute bool) error {
	err := receiverUseCase.receiversRepository.MuteStatus(ctx, email, mute)
	if err != nil {
		return err
	}

	return nil
}

func (receiverUseCase *ReceiversUseCase) GetAllUnmuted(ctx context.Context) ([]entities.Receiver, error) {
	unmutedReceivers, err := receiverUseCase.receiversRepository.GetAllUnmuted(ctx)
	if err != nil {
		return nil, err
	}

	return unmutedReceivers, nil
}

func (receiverUseCase *ReceiversUseCase) SendEmailToReceiver(ctx context.Context, message, subjectMessage string) error {
	unmuted, err := receiverUseCase.receiversRepository.GetAllUnmuted(ctx)
	if err != nil {
		return err
	}

	if len(unmuted) == 0 {
		return fmt.Errorf("no unmuted receivers found")
	}

	var emails []string
	for _, receiver := range unmuted {
		emails = append(emails, receiver.Email)
	}

	err = receiverUseCase.smtpService.SendEmail(ctx, message, subjectMessage, emails)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}
