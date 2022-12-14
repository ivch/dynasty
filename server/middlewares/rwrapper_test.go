package middlewares

import (
	"net/http/httptest"
	"testing"
)

func TestNewResponseWrapper(t *testing.T) {
	w := httptest.NewRecorder()
	ww := NewResponseWrapper(w)

	ww.WriteHeader(100)
	ww.Header().Set("hello", "world")
	n, err := ww.Write([]byte("Garry Goodini"))
	if err != nil {
		t.Fatal(err)
	}
	if n != 13 {
		t.Error("invalid bytes len written")
	}

	// check response recorder
	if w.Code != 100 {
		t.Error("response recorder invalid code passed")
	}
	if w.Header().Get("hello") != "world" {
		t.Error("response recorder invalid header key")
	}
	if w.Body.String() != "Garry Goodini" {
		t.Error("response recorder invalid body")
	}
}
