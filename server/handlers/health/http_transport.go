package health

import (
	"net/http"
)

// NewHTTPTransport returns a new instance of HTTPTransport.
// Implements http.Handler interface.
// Takes Checker interface which will be used to perform health checks.
func NewHTTPTransport(checker Checker) http.Handler { return &HTTPTransport{checker: checker} }

// HTTPTransport holds a checker which will be used to perform health checks.
type HTTPTransport struct {
	checker Checker
}

func (t *HTTPTransport) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := t.checker.Health(r.Context()); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
}
