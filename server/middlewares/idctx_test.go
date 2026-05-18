package middlewares_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/ivch/dynasty/server/middlewares"
	"github.com/stretchr/testify/assert"
)

func TestUserIDToCTX(t *testing.T) {
	ctx := context.Background()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)
	r.Header.Add(middlewares.UserIDHeader.String(), "2")
	newCtx := middlewares.UserIDToCTX(ctx, r)
	val := newCtx.Value(middlewares.UserIDCtxKey)
	assert.Equal(t, "2", val)
}

func TestUserIDFromContext(t *testing.T) {
	ctx := context.Background()
	_, ok := middlewares.UserIDFromContext(ctx)
	assert.False(t, ok)
	ctx = context.WithValue(ctx, middlewares.UserIDCtxKey, "2")
	val, ok := middlewares.UserIDFromContext(ctx)
	assert.True(t, ok)
	assert.Equal(t, "2", val)
	ctx = context.WithValue(context.Background(), middlewares.UserIDCtxKey, []string{"a", "b"})
	_, ok = middlewares.UserIDFromContext(ctx)
	assert.False(t, ok)
}
