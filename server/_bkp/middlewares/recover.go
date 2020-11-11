package middlewares

import (
	"net/http"
	"runtime/debug"

	"svc/pkg/logger"
)

// NewRecover returns new instance of Recover struct.
func NewRecover(log logger.Logger) *Recover { return &Recover{log: log} }

// Recover holds properties for middleware.
type Recover struct {
	log logger.Logger
}

// Middleware implements server.Middleware interface adding recover
// functionality.
func (m *Recover) Middleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				m.log.Error("panic happened: %v\nstacktrace: %s", r, debug.Stack())
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
