package requests

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/ivch/dynasty/models"
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
			request:  "bad json }}}",
			id:       "1",
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
			name:     "error wrong user id",
			request:  "{}",
			id:       "1",
			header:   "asd",
			wantErr:  true,
			wantCode: http.StatusInternalServerError,
		},
		{
			name:     "error no user",
			request:  "{}",
			id:       "1",
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
			rq, err := http.NewRequest("PUT", "/v1/request/"+tt.id, strings.NewReader(tt.request))
			fmt.Println(err)
			rq.Header.Add("X-Auth-User", tt.header)
			h.ServeHTTP(rr, rq)
			if (rr.Code != http.StatusOK) != tt.wantErr {
				t.Errorf("Request error. status = %d, wantErr %v", rr.Code, tt.wantErr)
			}
		})
	}
}

func TestHTTP_My(t *testing.T) {
	tests := []struct {
		name     string
		svc      Service
		query    string
		header   string
		wantErr  bool
		wantCode int
		want     string
	}{
		{
			name:     "error no user",
			header:   "0",
			wantErr:  true,
			wantCode: http.StatusInternalServerError,
		},
		{
			name:     "error wrong id",
			header:   "}{",
			wantErr:  true,
			wantCode: http.StatusInternalServerError,
		},
		{
			name:     "error no offset",
			header:   "1",
			query:    "?limit=1",
			wantErr:  true,
			wantCode: http.StatusInternalServerError,
		},
		{
			name:     "error no limit",
			header:   "1",
			query:    "?offset=1",
			wantErr:  true,
			wantCode: http.StatusInternalServerError,
		},
		{
			name:     "error bad offset",
			header:   "1",
			query:    "?offset=a&limit=1",
			wantErr:  true,
			wantCode: http.StatusInternalServerError,
		},
		{
			name:     "error bad limit",
			header:   "1",
			query:    "?offset=1&limit=as",
			wantErr:  true,
			wantCode: http.StatusInternalServerError,
		},
		{
			name:     "error too low limit",
			header:   "1",
			query:    "?offset=1&limit=0",
			wantErr:  true,
			wantCode: http.StatusInternalServerError,
		},
		{
			name:     "error too big limit",
			header:   "1",
			query:    "?offset=1&limit=300",
			wantErr:  true,
			wantCode: http.StatusInternalServerError,
		},
		{
			name:   "error service",
			query:  "?offset=1&limit=1",
			header: "1",
			svc: &ServiceMock{
				MyFunc: func(_ context.Context, _ *myRequest) (*myResponse, error) {
					return nil, errTestError
				},
			},
			wantErr:  true,
			wantCode: http.StatusInternalServerError,
		},
		{
			name:   "ok",
			query:  `?offset=0&limit=1`,
			header: "1",
			svc: &ServiceMock{
				MyFunc: func(_ context.Context, r *myRequest) (*myResponse, error) {
					return &myResponse{
						Data: []*models.Request{
							{
								ID:          1,
								Type:        "1",
								UserID:      1,
								Time:        1,
								Description: "1",
								Status:      "1",
							},
						},
					}, nil
				},
			},
			wantErr:  false,
			wantCode: http.StatusOK,
			want:     `{"data":[{"id":1,"type":"1","user_id":1,"time":1,"description":"1","status":"1"}]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := tt.svc
			h := newHTTPHandler(defaultLogger, svc)
			rr := httptest.NewRecorder()
			rq, _ := http.NewRequest("GET", "/v1/my"+tt.query, nil)
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

func TestHTTP_Delete(t *testing.T) {
	tests := []struct {
		name     string
		svc      Service
		id       string
		header   string
		wantErr  bool
		wantCode int
	}{
		{
			name:     "error no id",
			id:       "0",
			header:   "0",
			wantErr:  true,
			wantCode: http.StatusInternalServerError,
		},
		{
			name:     "error wrong id",
			id:       "asd",
			header:   "0",
			wantErr:  true,
			wantCode: http.StatusInternalServerError,
		},
		{
			name:     "error no user",
			header:   "0",
			wantErr:  true,
			wantCode: http.StatusInternalServerError,
		},
		{
			name:   "error service",
			header: "1",
			svc: &ServiceMock{
				DeleteFunc: func(_ context.Context, _ *byIDRequest) error {
					return errTestError
				},
			},
			wantErr:  true,
			wantCode: http.StatusInternalServerError,
		},
		{
			name:   "ok",
			id:     "1",
			header: "1",
			svc: &ServiceMock{
				DeleteFunc: func(_ context.Context, r *byIDRequest) error {
					if r.ID != 1 || r.UserID != 1 {
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
			rq, _ := http.NewRequest("DELETE", "/v1/request/"+tt.id, nil)
			rq.Header.Add("X-Auth-User", tt.header)
			h.ServeHTTP(rr, rq)
			if (rr.Code != http.StatusOK) != tt.wantErr {
				t.Errorf("Request error. status = %d, wantErr %v", rr.Code, tt.wantErr)
			}
		})
	}
}

func TestHTTP_Get(t *testing.T) {
	tests := []struct {
		name     string
		svc      Service
		id       string
		header   string
		wantErr  bool
		wantCode int
		want     string
	}{
		{
			name:     "error no id",
			id:       "0",
			header:   "0",
			wantErr:  true,
			wantCode: http.StatusInternalServerError,
		},
		{
			name:     "error wrong id",
			id:       "asd",
			header:   "0",
			wantErr:  true,
			wantCode: http.StatusInternalServerError,
		},
		{
			name:     "error no user",
			id:       "1",
			header:   "0",
			wantErr:  true,
			wantCode: http.StatusInternalServerError,
		},
		{
			name:     "error wrong user",
			id:       "1",
			header:   "asd",
			wantErr:  true,
			wantCode: http.StatusInternalServerError,
		},
		{
			name:   "error service",
			header: "1",
			svc: &ServiceMock{
				GetFunc: func(_ context.Context, _ *byIDRequest) (*getResponse, error) {
					return nil, errTestError
				},
			},
			wantErr:  true,
			wantCode: http.StatusInternalServerError,
		},
		{
			name:   "ok",
			id:     "1",
			header: "1",
			svc: &ServiceMock{
				GetFunc: func(_ context.Context, r *byIDRequest) (*getResponse, error) {
					if r.ID != 1 || r.UserID != 1 {
						return nil, errTestError
					}

					return &getResponse{
						&models.Request{
							ID:          1,
							Type:        "1",
							UserID:      1,
							Time:        1,
							Description: "1",
							Status:      "1",
						},
					}, nil
				},
			},
			wantErr:  false,
			wantCode: http.StatusOK,
			want:     `{"id":1,"type":"1","user_id":1,"time":1,"description":"1","status":"1"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := tt.svc
			h := newHTTPHandler(defaultLogger, svc)
			rr := httptest.NewRecorder()
			rq, _ := http.NewRequest("GET", "/v1/request/"+tt.id, nil)
			rq.Header.Add("X-Auth-User", tt.header)
			h.ServeHTTP(rr, rq)
			if (rr.Code != http.StatusOK) != tt.wantErr {
				t.Errorf("Request error. status = %d, wantErr %v", rr.Code, tt.wantErr)
			}
		})
	}
}
