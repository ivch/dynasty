package middlewares

import (
	"net/http"
	"strings"
	"time"

	"github.com/ivch/dynasty/common/logger"
)

// NewLogging returns new instance of Logging struct.
func NewLogging(log logger.Logger) *Logging { return &Logging{log: log} }

// Logging holds properties for middleware.
type Logging struct {
	log logger.Logger
}

// Middleware Middleware implements server.Middleware interface adding
// request logging functionality.
func (m *Logging) Middleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.RequestURI, "/ui/") {
			next.ServeHTTP(w, r)
			return
		}
		start := time.Now()
		ww := NewResponseWrapper(w)
		next.ServeHTTP(ww, r)
		duration := time.Since(start).Milliseconds()
		m.log.Info("%s - - \"%s %s %s\" %d %d", r.RemoteAddr, r.Method, r.RequestURI, r.Proto, ww.Code(), duration)
	}
	return http.HandlerFunc(fn)
}
