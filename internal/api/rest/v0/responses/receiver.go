package responses

type CreateReceiver struct {
	ID int `json:"id"`
}

type GetReceivers struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}
