package requests

type CreateServer struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
