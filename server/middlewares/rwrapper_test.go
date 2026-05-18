package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ivch/dynasty/server/middlewares"
)

func TestNewResponseWrapper(t *testing.T) {
	w := httptest.NewRecorder()
	ww := middlewares.NewResponseWrapper(w)

	ww.WriteHeader(http.StatusContinue)
	ww.Header().Set("hello", "world")

	// check response recorder
	if w.Code != 100 {
		t.Error("response recorder invalid code passed")
	}
	if w.Header().Get("hello") != "world" {
		t.Error("response recorder invalid header key")
	}
}
