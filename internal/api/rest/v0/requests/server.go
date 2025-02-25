package requests

type CreateServer struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}
