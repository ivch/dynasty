package middlewares

import (
	"context"
	"net/http"

	"github.com/ivch/dynasty/common/logger"
)

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

const (
	userIDHeader contextKey = "X-Auth-User"
	userIDCtxKey contextKey = "userID"
)

func NewIDCtx(log logger.Logger) *IDCtx { return &IDCtx{log: log} }

type IDCtx struct {
	log logger.Logger
}

func (m *IDCtx) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if v := r.Header.Get(userIDHeader.String()); v != "" {
			ctx = context.WithValue(r.Context(), userIDCtxKey, v)
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UserIDToCTX(ctx context.Context, r *http.Request) context.Context {
	if v := r.Header.Get(userIDHeader.String()); v != "" {
		ctx = context.WithValue(ctx, userIDCtxKey, v)
	}

	return ctx
}

func UserIDFromContext(ctx context.Context) (string, bool) {
	v := ctx.Value(userIDCtxKey)
	if v == nil {
		return "", false
	}
	id, ok := v.(string)
	return id, ok
}
