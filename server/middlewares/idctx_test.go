package middlewares

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserIDToCTX(t *testing.T) {
	ctx := context.Background()
	r, _ := http.NewRequest("GET", "/", nil)
	r.Header.Add(userIDHeader, "2")
	newCtx := UserIDToCTX(ctx, r)
	val := newCtx.Value(userIDCtxKey)
	assert.Equal(t, val, "2")
}

func TestUserIDFromContext(t *testing.T) {
	ctx := context.Background()
	_, ok := UserIDFromContext(ctx)
	assert.Equal(t, ok, false)
	ctx = context.WithValue(ctx, userIDCtxKey, "2")
	val, ok := UserIDFromContext(ctx)
	assert.Equal(t, ok, true)
	assert.Equal(t, val, "2")
	ctx = context.WithValue(context.Background(), userIDCtxKey, []string{"a", "b"})
	_, ok = UserIDFromContext(ctx)
	assert.Equal(t, ok, false)
}
