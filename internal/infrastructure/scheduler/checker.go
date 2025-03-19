package scheduler

import (
	"checker/internal/domain/app/repositories"
	"checker/internal/domain/app/usecases"
	"checker/internal/domain/entities"
	"checker/internal/shared"
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

const (
	logFile              = "logs/errors.log"
	consecutiveFailCount = 3
)

var (
	lastNotifyTime   = make(map[string]time.Time)
	isServerDown     = make(map[string]bool)
	consecutiveFails = make(map[string]int)
	mu               sync.Mutex
)

func CheckServer(ctx context.Context, server entities.Server, basicRepo repositories.Basic, receiverUseCase *usecases.ReceiversUseCase) error {
	basicConfig, err := basicRepo.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch basic config: %v", err)
	}

	timeout := time.Duration(basicConfig.Timeout) * time.Second
	notificationInterval := time.Duration(basicConfig.NotificationInterval) * time.Hour

	client := &http.Client{
		Timeout: timeout,
	}

	resp, err := client.Get(server.URL)
	if err != nil {
		HandleError(ctx, notificationInterval, server, receiverUseCase, 0)
		return fmt.Errorf("ERROR: Server - %s (%s) is unreachable", server.Name, server.URL)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 204 {
		HandleError(ctx, notificationInterval, server, receiverUseCase, resp.StatusCode)
		return fmt.Errorf("ERROR: Server - %s (%s) returned status %d", server.Name, server.URL, resp.StatusCode)
	}

	mu.Lock()
	if isServerDown[server.Name] {
		SendRecoveryNotification(ctx, receiverUseCase, server)
	}
	isServerDown[server.Name] = false
	consecutiveFails[server.Name] = 0
	delete(lastNotifyTime, server.Name)
	mu.Unlock()

	return nil
}

func HandleError(ctx context.Context, notificationInterval time.Duration, server entities.Server, receiverUseCase *usecases.ReceiversUseCase, statusCode int) {
	mu.Lock()
	defer mu.Unlock()

	consecutiveFails[server.Name]++

	if consecutiveFails[server.Name] < consecutiveFailCount {
		return
	}

	if time.Since(lastNotifyTime[server.Name]) < notificationInterval {
		return
	}

	subjectMessage := fmt.Sprintf("%s is failed", server.Name)
	statusText := fmt.Sprintf("%d", statusCode)
	if statusCode == 0 {
		statusText = "Unreachable"
	}

	msg := fmt.Sprintf(
		"\U0001F6A8 Service Health Alert \U0001F6A8\n\n"+
			"Name: %s\n"+
			"URL: %s\n"+
			"Status Code: %s\n\n"+
			"Please check the service for issues.",
		server.Name, server.URL, statusText,
	)

	shared.WriteLog(msg, logFile)
	receiverUseCase.SendEmailToReceiver(ctx, msg, subjectMessage)

	lastNotifyTime[server.Name] = time.Now()
	isServerDown[server.Name] = true
}

func SendRecoveryNotification(ctx context.Context, receiverUseCase *usecases.ReceiversUseCase, server entities.Server) {
	subjectMessage := fmt.Sprintf("%s is back online", server.Name)
	msg := fmt.Sprintf("\U0001F7E2 Server Recovery Notice \U0001F7E2\n\nServer - %s (%s) is back online.", server.Name, server.URL)

	shared.WriteLog(msg, logFile)
	receiverUseCase.SendEmailToReceiver(ctx, msg, subjectMessage)

	mu.Lock()
	consecutiveFails[server.Name] = 0
	mu.Unlock()
}
