package health

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestNewHTTPTransport(t *testing.T) {
	t.Run("New", func(t *testing.T) {
		checker := testChecker{}

		want := &HTTPTransport{
			checker: &checker,
		}

		if got := NewHTTPTransport(&checker); !reflect.DeepEqual(got, want) {
			t.Errorf("NewHTTPTransport() = %v, want %v", got, want)
		}
	})
}

func TestHTTPTransport_ServeHTTP(t *testing.T) {
	type tcase struct {
		checker        Checker
		wantStatusCode int
	}
	tests := map[string]tcase{
		"nil": {
			checker: &testChecker{HealthFunc: func(ctx context.Context) error {
				return nil
			}},
			wantStatusCode: http.StatusOK,
		},
		"err": {
			checker: &testChecker{HealthFunc: func(ctx context.Context) error {
				return errors.New("test error")
			}},
			wantStatusCode: http.StatusServiceUnavailable,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			transport := &HTTPTransport{
				checker: tc.checker,
			}
			srv := httptest.NewServer(transport)
			defer srv.Close()

			resp, err := http.Get(srv.URL)
			if err != nil {
				t.Errorf("request failed: %v", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode != tc.wantStatusCode {
				t.Errorf("health check failed, wanted status: %d, got: %d", tc.wantStatusCode, resp.StatusCode)
			}
		})
	}
}
