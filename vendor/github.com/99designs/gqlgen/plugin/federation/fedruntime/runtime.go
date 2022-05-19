package fedruntime

// Service is the service object that the
// generated.go file will return for the _service
// query
type Service struct {
	SDL string `json:"sdl"`
}

// Everything with a @key implements this
type Entity interface {
	IsEntity()
}
