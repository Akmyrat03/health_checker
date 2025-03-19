package requests

type UpdateBasic struct {
	CheckInterval        int `json:"checkIntervalInSeconds"`
	Timeout              int `json:"timeoutInSeconds"`
	NotificationInterval int `json:"notificationIntervalInHours"`
}
