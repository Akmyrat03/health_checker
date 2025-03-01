package scheduler

import (
	"checker/internal/domain/app/repositories"
	"checker/internal/domain/app/usecases"
	"checker/internal/domain/entities"
	"checker/internal/shared"
	"context"
	"fmt"
	"net/http"
	"time"
)

const (
	logFile = "logs/errors.log"
)

var (
	lastErrorTime  = make(map[string]time.Time)
	lastNotifyTime = make(map[string]time.Time)
	isServerDown   = make(map[string]bool)
)

func CheckServer(ctx context.Context, server entities.Server, basicRepo repositories.Basic, receiverUseCase *usecases.ReceiversUseCase) error {
	basicConfig, err := basicRepo.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch basic config: %v", err)
	}

	timeout := time.Duration(basicConfig.Timeout) * time.Second
	notificationInterval := time.Duration(basicConfig.NotificationInterval) * time.Minute

	client := &http.Client{
		Timeout: timeout,
	}

	resp, err := client.Get(server.Url)
	if err != nil {
		HandleError(ctx, notificationInterval, server, receiverUseCase)
		return fmt.Errorf("ERROR: Server - %s (%s) is unreachable", server.Name, server.Url)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		HandleError(ctx, notificationInterval, server, receiverUseCase)
		return fmt.Errorf("ERROR: Server - %s (%s) returned status %d", server.Name, server.Url, resp.StatusCode)
	}

	if isServerDown[server.Name] {
		SendRecoveryNotification(ctx, receiverUseCase, server)
	}

	isServerDown[server.Name] = false
	delete(lastErrorTime, server.Name)
	delete(lastNotifyTime, server.Name)

	return nil
}

func HandleError(ctx context.Context, notificationInterval time.Duration, server entities.Server, receiverUseCase *usecases.ReceiversUseCase) {
	msg := fmt.Sprintf("ERROR: Server - %s (%s) is unreachable or returned status 500", server.Name, server.Url)

	if !isServerDown[server.Name] { // First time the server goes down
		shared.WriteLog(msg, logFile) // Log the first error
		receiverUseCase.SendEmailToReceiver(ctx, msg)
		lastErrorTime[server.Name] = time.Now()
		lastNotifyTime[server.Name] = time.Now()
		isServerDown[server.Name] = true
	} else if time.Since(lastNotifyTime[server.Name]) > notificationInterval {
		shared.WriteLog(msg, logFile)
		receiverUseCase.SendEmailToReceiver(ctx, msg)
		lastNotifyTime[server.Name] = time.Now()
	}
}

func SendRecoveryNotification(ctx context.Context, receiverUseCase *usecases.ReceiversUseCase, server entities.Server) {
	msg := fmt.Sprintf("Server - %s (%s) is back online", server.Name, server.Url)
	shared.WriteLog(msg, logFile)
	receiverUseCase.SendEmailToReceiver(ctx, msg)
}
