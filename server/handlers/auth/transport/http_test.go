package transport

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/ivch/dynasty/common/logger"
	"github.com/ivch/dynasty/server/handlers/auth"
	"github.com/ivch/dynasty/server/middlewares"
)

var (
	defaultLogger *logger.StdLog
	errTestError  = errors.New("some err")
)

func TestMain(m *testing.M) {
	defaultLogger = logger.NewStdLog(logger.WithWriter(ioutil.Discard))
	os.Exit(m.Run())
}

func TestHTTP_Login(t *testing.T) {
	tests := []struct {
		name     string
		request  string
		svc      AuthService
		wantErr  bool
		wantCode int
		want     string
	}{
		{
			name:     "error decode request",
			request:  "}{",
			wantErr:  true,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "error empty password",
			request:  `{"phone":"123"}`,
			wantErr:  true,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "error empty phone",
			request:  `{"password":"123456"}`,
			wantErr:  true,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "error bad phone",
			request:  `{"phone":"123qwe123qwe", "password":"123456"}`,
			wantErr:  true,
			wantCode: http.StatusBadRequest,
		},
		{
			name:    "error service error",
			request: `{"password":"123456", "phone":"123123123123"}`,
			svc: &AuthServiceMock{
				LoginFunc: func(_ context.Context, _, _ string) (*auth.Tokens, error) {
					return nil, errTestError
				},
			},
			wantErr:  true,
			wantCode: http.StatusInternalServerError,
		},
		{
			name:    "ok",
			request: `{"password":"123456", "phone":"123123123123"}`,
			svc: &AuthServiceMock{
				LoginFunc: func(_ context.Context, _, _ string) (*auth.Tokens, error) {
					return &auth.Tokens{
						AccessToken:  "at",
						RefreshToken: "rt",
					}, nil
				},
			},
			wantErr:  false,
			want:     `{"access_token":"at","refresh_token":"rt"}`,
			wantCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := tt.svc
			h := NewHTTPTransport(defaultLogger, svc, middlewares.NewIDCtx(defaultLogger).Middleware)
			rr := httptest.NewRecorder()
			rq, _ := http.NewRequest("POST", "/v1/login", strings.NewReader(tt.request))
			h.ServeHTTP(rr, rq)
			if (rr.Code != tt.wantCode) && tt.wantErr {
				t.Errorf("Request error. status = %d, wantCode = %d, wantErr %v", rr.Code, tt.wantCode, tt.wantErr)
			}
			if !tt.wantErr && tt.want != strings.TrimSpace(rr.Body.String()) {
				t.Errorf("Response error, got = %s, want = %s", rr.Body.String(), tt.want)
			}
		})
	}
}

func TestHTTP_Refresh(t *testing.T) {
	tests := []struct {
		name     string
		request  string
		svc      AuthService
		wantErr  bool
		wantCode int
		want     string
	}{
		{
			name:     "error decode request",
			request:  "}{",
			wantErr:  true,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "error no token in request",
			request:  `{"token":""}`,
			wantErr:  true,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "error bad token in request",
			request:  `{"token":"asdsadasdas"}`,
			wantErr:  true,
			wantCode: http.StatusBadRequest,
		},
		{
			name:    "error service error",
			request: `{"token":"d3ffebcf-1cec-441b-93d6-a984b7647d48"}`,
			svc: &AuthServiceMock{
				RefreshFunc: func(_ context.Context, _ string) (*auth.Tokens, error) {
					return nil, errTestError
				},
			},
			wantErr:  true,
			wantCode: http.StatusInternalServerError,
		},
		{
			name:    "ok",
			request: `{"token":"d3ffebcf-1cec-441b-93d6-a984b7647d48"}`,
			svc: &AuthServiceMock{
				RefreshFunc: func(_ context.Context, _ string) (*auth.Tokens, error) {
					return &auth.Tokens{
						AccessToken:  "at",
						RefreshToken: "rt",
					}, nil
				},
			},
			want:     `{"access_token":"at","refresh_token":"rt"}`,
			wantErr:  true,
			wantCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := tt.svc
			h := NewHTTPTransport(defaultLogger, svc, middlewares.NewIDCtx(defaultLogger).Middleware)
			rr := httptest.NewRecorder()
			rq, _ := http.NewRequest("POST", "/v1/refresh", strings.NewReader(tt.request))
			h.ServeHTTP(rr, rq)
			if (rr.Code != tt.wantCode) && tt.wantErr {
				t.Errorf("Request error. status = %d, wantCode = %d, wantErr %v", rr.Code, tt.wantCode, tt.wantErr)
			}
			if !tt.wantErr && tt.want != strings.TrimSpace(rr.Body.String()) {
				t.Errorf("Response error, got = %s, want = %s", rr.Body.String(), tt.want)
			}
		})
	}
}

func TestHTTP_Gwfa(t *testing.T) {
	tests := []struct {
		name     string
		header   string
		svc      AuthService
		wantErr  bool
		wantCode int
		want     string
	}{
		{
			name:     "error empty auth header",
			header:   "",
			wantErr:  true,
			wantCode: http.StatusUnauthorized,
		},
		{
			name:   "error service error",
			header: "Bearer token",
			svc: &AuthServiceMock{
				GwfaFunc: func(string) (uint, error) {
					return 0, errTestError
				},
			},
			wantErr:  true,
			wantCode: http.StatusUnauthorized,
		},
		{
			name:   "ok",
			header: "Bearer token",
			svc: &AuthServiceMock{
				GwfaFunc: func(string) (uint, error) {
					return 1, nil
				},
			},
			wantErr:  false,
			wantCode: http.StatusOK,
			want:     "null",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := tt.svc
			h := NewHTTPTransport(defaultLogger, svc, middlewares.NewIDCtx(defaultLogger).Middleware)
			rr := httptest.NewRecorder()
			rq, _ := http.NewRequest("GET", "/v1/gwfa", nil)
			rq.Header.Add("Authorization", tt.header)
			h.ServeHTTP(rr, rq)
			if (rr.Code != tt.wantCode) && tt.wantErr {
				t.Errorf("Request error. status = %d, wantCode = %d, wantErr %v", rr.Code, tt.wantCode, tt.wantErr)
			}
			if !tt.wantErr && tt.want != strings.TrimSpace(rr.Body.String()) {
				t.Errorf("Response error, got = %s, want = %s", rr.Body.String(), tt.want)
			}
			if rr.Code == http.StatusOK && rr.Header().Get("X-Auth-User") != "1" {
				t.Errorf("Wrogn header, got = %s, want = %s", rr.Header().Get("X-Auth-User"), "1")
			}
		})
	}
}

func TestHTTP_Logout(t *testing.T) {
	tests := []struct {
		name     string
		header   string
		svc      AuthService
		wantErr  bool
		wantCode int
		want     string
	}{
		{
			name:     "error no user",
			wantErr:  true,
			wantCode: http.StatusUnauthorized,
		},
		{
			name:     "error bad user id #1",
			header:   "aswd",
			wantErr:  true,
			wantCode: http.StatusUnauthorized,
		},
		{
			name:     "error bad user id #2",
			header:   "0",
			wantErr:  true,
			wantCode: http.StatusUnauthorized,
		},
		{
			name:   "error service error",
			header: "1",
			svc: &AuthServiceMock{
				LogoutFunc: func(_ context.Context, _ uint) error {
					return errTestError
				},
			},
			wantErr:  true,
			wantCode: http.StatusInternalServerError,
		},
		{
			name:   "ok",
			header: "1",
			svc: &AuthServiceMock{
				LogoutFunc: func(_ context.Context, _ uint) error {
					return nil
				},
			},
			wantErr:  false,
			wantCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := tt.svc
			h := NewHTTPTransport(defaultLogger, svc, middlewares.NewIDCtx(defaultLogger).Middleware)
			rr := httptest.NewRecorder()
			rq, _ := http.NewRequest("GET", "/v1/logout", nil)
			rq.Header.Add("X-Auth-User", tt.header)
			h.ServeHTTP(rr, rq)
			if (rr.Code != tt.wantCode) && tt.wantErr {
				t.Errorf("Request error. status = %d, wantCode = %d, wantErr %v", rr.Code, tt.wantCode, tt.wantErr)
			}
		})
	}
}
