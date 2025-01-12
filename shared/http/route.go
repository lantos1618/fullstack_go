package http

// Route represents an API route with its associated request and response types
type Route[Req any, Res any] struct {
	Path   string
	Method string
}

// NewRoute creates a new route with the given path and method
func NewRoute[Req any, Res any](path, method string) Route[Req, Res] {
	return Route[Req, Res]{
		Path:   path,
		Method: method,
	}
}

// Common HTTP methods
const (
	MethodGet    = "GET"
	MethodPost   = "POST"
	MethodPut    = "PUT"
	MethodDelete = "DELETE"
	MethodPatch  = "PATCH"
)
