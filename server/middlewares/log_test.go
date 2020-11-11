package middlewares

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestLogging_Middleware(t *testing.T) {
	logMock := &mockLogger{}
	m := &Logging{
		log: logMock,
	}
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	})
	mw := m.Middleware(next)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/hello", nil)
	r.RemoteAddr = "10.10.10.10"

	mw.ServeHTTP(w, r)

	if w.Code != http.StatusAccepted {
		t.Error("invalid status")
	}
	argsToCheck := []interface{}{"10.10.10.10", http.MethodPost, r.RequestURI,
		r.Proto, http.StatusAccepted}
	if len(logMock.info.v) != 6 {
		t.Error("invalid log data len")
	}
	for i := range argsToCheck {
		if !reflect.DeepEqual(argsToCheck[i], logMock.info.v[i]) {
			t.Error("invalid log data", argsToCheck[i], logMock.info.v[i])
		}
	}
}
