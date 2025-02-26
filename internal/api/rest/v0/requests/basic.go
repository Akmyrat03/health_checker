package requests

type UpdateBasic struct {
	CheckInterval        int `json:"check_interval"`
	Timeout              int `json:"timeout"`
	NotificationInterval int `json:"notification_interval"`
}
