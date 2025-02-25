package entities

type Server struct {
	ID   int
	Name string
	Url  string
}

type Basic struct {
	CheckInterval string
	Timeout       string
	ErrorInterval string
}
