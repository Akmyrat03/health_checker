package responses

type CreateServer struct {
	ID int `json:"id"`
}

type GetServer struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}
