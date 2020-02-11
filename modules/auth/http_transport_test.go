package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ivch/dynasty/models/dto"
)

func TestHTTP_Login(t *testing.T) {
	tests := []struct {
		name     string
		request  string
		svc      Service
		wantErr  bool
		wantCode int
		want     string
	}{
		{
			name:     "error decode request",
			request:  "}{",
			wantErr:  true,
			wantCode: http.StatusInternalServerError,
		},
		{
			name:     "error empty phone",
			request:  `{"password":"123"}`,
			wantErr:  true,
			wantCode: http.StatusInternalServerError,
		},
		{
			name:     "error empty password",
			request:  `{"phone":"123"}`,
			wantErr:  true,
			wantCode: http.StatusInternalServerError,
		},
		{
			name:    "error service error",
			request: `{"password":"123", "phone":"123"}`,
			svc: &ServiceMock{
				LoginFunc: func(_ context.Context, _ *dto.AuthLoginRequest) (*dto.AuthLoginResponse, error) {
					return nil, errTestError
				},
			},
			wantErr:  true,
			wantCode: http.StatusInternalServerError,
		},
		{
			name:    "ok",
			request: `{"password":"123", "phone":"123"}`,
			svc: &ServiceMock{
				LoginFunc: func(_ context.Context, _ *dto.AuthLoginRequest) (*dto.AuthLoginResponse, error) {
					return &dto.AuthLoginResponse{
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
			h := newHTTPHandler(defaultLogger, svc)
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
		svc      Service
		wantErr  bool
		wantCode int
		want     string
	}{
		{
			name:     "error decode request",
			request:  "}{",
			wantErr:  true,
			wantCode: http.StatusInternalServerError,
		},
		{
			name:     "error no token in request",
			request:  `{"token":""}`,
			wantErr:  true,
			wantCode: http.StatusInternalServerError,
		},
		{
			name:    "error service error",
			request: `{"token":"token"}`,
			svc: &ServiceMock{
				RefreshFunc: func(_ context.Context, _ *dto.AuthRefreshTokenRequest) (*dto.AuthLoginResponse, error) {
					return nil, errTestError
				},
			},
			wantErr:  true,
			wantCode: http.StatusInternalServerError,
		},
		{
			name:    "ok",
			request: `{"token":"token"}`,
			svc: &ServiceMock{
				RefreshFunc: func(_ context.Context, _ *dto.AuthRefreshTokenRequest) (*dto.AuthLoginResponse, error) {
					return &dto.AuthLoginResponse{
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
			h := newHTTPHandler(defaultLogger, svc)
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
		svc      Service
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
			svc: &ServiceMock{
				GwfaFunc: func(string) (uint, error) {
					return 0, errTestError
				},
			},
			wantErr:  true,
			wantCode: http.StatusForbidden,
		},
		{
			name:   "ok",
			header: "Bearer token",
			svc: &ServiceMock{
				GwfaFunc: func(string) (uint, error) {
					return 1, nil
				},
			},
			wantErr:  false,
			wantCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := tt.svc
			h := newHTTPHandler(defaultLogger, svc)
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
