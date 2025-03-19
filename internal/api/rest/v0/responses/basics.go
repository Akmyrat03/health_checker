package responses

type GetBasicConfig struct {
	CheckInterval        int `json:"checkIntervalInSeconds"`
	Timeout              int `json:"timeoutInSeconds"`
	NotificationInterval int `json:"notificationIntervalInHours"`
}
