package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ivch/dynasty/server/middlewares"
)

func TestRecover_Middleware(t *testing.T) {
	logMock := &mockLogger{}
	rmw := middlewares.NewRecover(logMock)
	mw := rmw.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("panic")
	}))
	w, r := httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/", nil)
	mw.ServeHTTP(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Error("invalid response status")
	}

	if strings.HasPrefix(logMock.info.f, "panic happened") {
		t.Log(logMock.error.f)
		t.Error("invalid log")
	}
}
