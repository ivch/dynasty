package server

import (
	"context"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	log := &testLogger{}
	health := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	hh := map[string]http.Handler{"/health": health}
	t.Run("tt.name", func(t *testing.T) {
		got, err := New(":8080", log, hh)
		if err != nil {
			t.Fatal(err)
		}

		if got.server.Addr != ":8080" {
			t.Errorf("expected :8080 server addres, but got: %s", got.server.Addr)
		}
		if reflect.DeepEqual(got.router, http.NewServeMux()) {
			t.Errorf("expected mux: %v, but got: %v", http.NewServeMux(), got.router)
		}
	})
}

func TestServer_Serve(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	hh := map[string]http.Handler{"/handle": handler}
	srv, err := New(":3000", &testLogger{}, hh)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	srvExitWaiter := func(ch chan struct{}) {
		err := srv.Serve(ctx)
		if err != nil {
			t.Error(err)
		}
		ch <- struct{}{}
	}
	onSrvExit := make(chan struct{})
	go srvExitWaiter(onSrvExit)
	cancel()
	select {
	case <-time.After(time.Second * 1):
		t.Error("timeout exit")
	case <-onSrvExit:
	}
}

type testLogger struct{}

func (l *testLogger) Debug(s string, v ...interface{}) {}

func (l *testLogger) Info(s string, v ...interface{}) {}

func (l *testLogger) Warn(s string, v ...interface{}) {}

func (l *testLogger) Error(s string, v ...interface{}) {}
