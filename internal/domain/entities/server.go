package entities

type Server struct {
	ID   int
	Name string
	Url  string
}

type Basic struct {
	CheckInterval int
	Timeout       int
	ErrorInterval int
}
