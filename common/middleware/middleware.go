package middleware

import (
	"context"
	"net/http"
)

const (
	userIDHeader = "X-Auth-User"
	userIDCtxKey = "userID"
)

func UserIDToCTX(ctx context.Context, r *http.Request) context.Context {
	if v := r.Header.Get(userIDHeader); v != "" {
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
