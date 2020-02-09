package requests

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestHTTP_Create(t *testing.T) {
	tests := []struct {
		name     string
		svc      Service
		request  string
		header   string
		wantErr  bool
		want     string
		wantCode int
	}{
		{
			name:     "error parsing request",
			request:  "}{",
			wantErr:  true,
			wantCode: http.StatusInternalServerError,
		},
		{
			name:     "error no user",
			request:  "{}",
			header:   "0",
			wantErr:  true,
			wantCode: http.StatusInternalServerError,
		},
		{
			name:     "error no type",
			request:  `{"time":1,"description":"abc"}`,
			header:   "1",
			wantErr:  true,
			wantCode: http.StatusInternalServerError,
		},
		{
			name:     "error no time",
			request:  `{"type":"1","description":"abc"}`,
			header:   "1",
			wantErr:  true,
			wantCode: http.StatusInternalServerError,
		},
		{
			name:    "error service",
			request: `{"type":"1","description":"abc","time":1}`,
			header:  "1",
			svc: &ServiceMock{
				CreateFunc: func(_ context.Context, _ *createRequest) (*createResponse, error) {
					return nil, errTestError
				},
			},
			wantErr:  true,
			wantCode: http.StatusInternalServerError,
		},
		{
			name:    "ok",
			request: `{"type":"1","description":"abc","time":1}`,
			header:  "1",
			svc: &ServiceMock{
				CreateFunc: func(_ context.Context, _ *createRequest) (*createResponse, error) {
					return &createResponse{ID: 1}, nil
				},
			},
			want:     `{"id":1}`,
			wantErr:  false,
			wantCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := tt.svc
			h := newHTTPHandler(defaultLogger, svc)
			rr := httptest.NewRecorder()
			rq, _ := http.NewRequest("POST", "/v1/request", strings.NewReader(tt.request))
			rq.Header.Add("X-Auth-User", tt.header)
			h.ServeHTTP(rr, rq)
			if (rr.Code != http.StatusOK) != tt.wantErr {
				t.Errorf("Request error. status = %d, wantErr %v", rr.Code, tt.wantErr)
			}

			if !tt.wantErr && tt.want != strings.TrimSpace(rr.Body.String()) {
				t.Errorf("Response error, got = %s, want = %s", rr.Body.String(), tt.want)
			}
		})
	}
}

func TestHTTP_Update(t *testing.T) {
	tests := []struct {
		name     string
		svc      Service
		id       string
		request  string
		header   string
		wantErr  bool
		wantCode int
	}{
		{
			name:     "error parsing request",
			request:  "}{",
			wantErr:  true,
			wantCode: http.StatusInternalServerError,
		},
		{
			name:     "error no id",
			request:  "{}",
			id:       "0",
			header:   "0",
			wantErr:  true,
			wantCode: http.StatusInternalServerError,
		},
		{
			name:     "error wrong id",
			request:  "{}",
			id:       "asd",
			header:   "0",
			wantErr:  true,
			wantCode: http.StatusInternalServerError,
		},
		{
			name:     "error no user",
			request:  "{}",
			header:   "0",
			wantErr:  true,
			wantCode: http.StatusInternalServerError,
		},
		{
			name:    "error service",
			request: `{"type":"1","description":"abc","time":1}`,
			header:  "1",
			svc: &ServiceMock{
				UpdateFunc: func(_ context.Context, _ *updateRequest) error {
					return errTestError
				},
			},
			wantErr:  true,
			wantCode: http.StatusInternalServerError,
		},
		{
			name:    "ok",
			request: `{"type":"1","description":"abc","time":1,"status":"new"}`,
			id:      "1",
			header:  "1",
			svc: &ServiceMock{
				UpdateFunc: func(_ context.Context, r *updateRequest) error {
					expected := &updateRequest{
						ID:          1,
						UserID:      1,
						Type:        func(s string) *string { return &s }("1"),
						Time:        func(s int64) *int64 { return &s }(1),
						Description: func(s string) *string { return &s }("abc"),
						Status:      func(s string) *string { return &s }("new"),
					}

					if !reflect.DeepEqual(expected, r) {
						return errTestError
					}
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
			h := newHTTPHandler(defaultLogger, svc)
			rr := httptest.NewRecorder()
			rq, _ := http.NewRequest("PUT", "/v1/request/"+tt.id, strings.NewReader(tt.request))
			rq.Header.Add("X-Auth-User", tt.header)
			h.ServeHTTP(rr, rq)
			if (rr.Code != http.StatusOK) != tt.wantErr {
				t.Errorf("Request error. status = %d, wantErr %v", rr.Code, tt.wantErr)
			}
		})
	}
}
