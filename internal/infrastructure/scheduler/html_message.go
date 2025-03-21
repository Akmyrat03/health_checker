package scheduler

import (
	"checker/internal/domain/app/usecases"
	"checker/internal/domain/entities"
	"checker/internal/shared"
	"context"
	"fmt"
)

func SendFailureNotificationHTML(ctx context.Context, server entities.Server, receiver entities.Receiver, receiverUseCase *usecases.ReceiversUseCase, statusCode int) {
	notificationNum := notifyCount[server.Name]
	subject := fmt.Sprintf("❗⚠️ %s - Failed x%d", server.Name, notificationNum)

	statusText := fmt.Sprintf("%d Internal Server Error", statusCode)
	if statusCode == 0 {
		statusText = "Unreachable"
	}

	intervalText := shared.FormatDuration(notifyInterval[server.Name])
	apiBaseURL := "http://199.40.7.56:3030/api/v0/receiver"

	message := fmt.Sprintf(`
	<html>
	<head>
		<style>
			body { font-family: Arial, sans-serif; background-color: #f4f4f4; padding: 20px; }
			.container { background: white; padding: 20px; border-radius: 10px; box-shadow: 0px 0px 10px rgba(0, 0, 0, 0.1); }
			h2 { color: #d9534f; }
			p { font-size: 16px; }
			.status { font-weight: bold; color: #d9534f; }
			.buttons { margin-top: 20px; }
			.button {
				display: inline-block; padding: 10px 20px; color: white; border-radius: 5px;
				text-decoration: none; font-weight: bold; margin-right: 10px;
			}
			.mute { background-color: #ffc107; }
			.unmute { background-color: #17a2b8; }
		</style>
	</head>
	<body>
		<div class="container">
			<p><strong>Name:</strong> %s</p>
			<p><strong>URL:</strong> <a href="%s">%s</a></p>
			<p><strong>Status Code:</strong> <span class="status">%s</span></p>
			<p><strong>Next notification in:</strong> %s</p>

			<div class="buttons">
				<a href="%s/mute?email=%s" class="button mute">🔕 Mute</a>
				<a href="%s/unmute?email=%s" class="button unmute">🔔 Unmute</a>
			</div>
		</div>
	</body>
	</html>`,
		server.Name, server.URL, server.URL, statusText, intervalText,
		apiBaseURL, receiver.Email, apiBaseURL, receiver.Email)

	receiverUseCase.SendEmailToReceiver(ctx, message, subject)
}

func SendRecoveryNotificationHTML(ctx context.Context, receiver entities.Receiver, receiverUseCase *usecases.ReceiversUseCase, server entities.Server) {
	subjectMessage := fmt.Sprintf("✅ %s - Back Online", server.Name)

	apiBaseURL := "http://199.40.7.56:3030/api/v0/receiver"

	muteURL := fmt.Sprintf("%s/mute?email=%s", apiBaseURL, receiver.Email)
	unmuteURL := fmt.Sprintf("%s/unmute?email=%s", apiBaseURL, receiver.Email)

	htmlMessage := fmt.Sprintf(`
	<html>
	<head>
		<style>
			body { font-family: Arial, sans-serif; background-color: #f4f4f4; padding: 20px; }
			.container { background: white; padding: 20px; border-radius: 10px; box-shadow: 0px 0px 10px rgba(0, 0, 0, 0.1); text-align: center; }
			h2 { color: #28a745; font-size: 24px; }
			p { font-size: 18px; color: #333; }
			.status { font-weight: bold; color: #28a745; }
			hr { border: none; border-top: 1px solid #ddd; margin: 20px 0; }
			.buttons { margin-top: 20px; }
			.button {
				display: inline-block; padding: 10px 20px; color: white; border-radius: 5px;
				text-decoration: none; font-weight: bold; margin: 5px;
			}
			.mute { background-color: #ffc107; }
			.unmute { background-color: #17a2b8; }
		</style>
	</head>
	<body>
		<div class="container">
			<h2>✅ %s is Back Online</h2>
			<p><strong>Server Name:</strong> %s</p>
			<p><strong>Server URL:</strong> <a href="%s" style="color: #007bff;">%s</a></p>
			<hr>
			<p>Manage this server:</p>
			<div class="buttons">
				<a href="%s" class="button mute">🔕 Mute</a>
				<a href="%s" class="button unmute">🔔 Unmute</a>
			</div>
		</div>
	</body>
	</html>`, server.Name, server.Name, server.URL, server.URL, muteURL, unmuteURL)

	receiverUseCase.SendEmailToReceiver(ctx, htmlMessage, subjectMessage)

	mu.Lock()
	consecutiveFails[server.Name] = 0
	delete(notifyInterval, server.Name)
	delete(notifyCount, server.Name)
	mu.Unlock()
}
