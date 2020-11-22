package server

import (
	"net/http"
)

// Middleware represents logic that should be implemented
// by any type which needs to middleware for http.Handler.
type Middleware interface {
	// Middleware takes http.Handler wraps it with additional functionality
	// and returns http.Handler.
	Middleware(next http.Handler) http.Handler
}
