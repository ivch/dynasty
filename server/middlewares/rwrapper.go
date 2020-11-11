package middlewares

import (
	"net/http"
)

// NewResponseWrapper returns new instance of Wrapper struct.
func NewResponseWrapper(w http.ResponseWriter) *Wrapper {
	wrapper := &Wrapper{
		w:    w,
		code: 0,
	}
	return wrapper
}

// Wrapper represents http.ResponseWriter wrapper which holds response code.
type Wrapper struct {
	w    http.ResponseWriter
	code int
}

// Code returns http.Response code after response.
func (w *Wrapper) Code() int {
	return w.code
}

func (w *Wrapper) Header() http.Header {
	return w.w.Header()
}

func (w *Wrapper) Write(bytes []byte) (int, error) {
	return w.w.Write(bytes)
}

func (w *Wrapper) WriteHeader(statusCode int) {
	w.code = statusCode
	w.w.WriteHeader(statusCode)
}
