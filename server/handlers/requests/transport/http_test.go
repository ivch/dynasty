package transport

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/microcosm-cc/bluemonday"

	"github.com/ivch/dynasty/common/logger"
	"github.com/ivch/dynasty/server/handlers/requests"
	"github.com/ivch/dynasty/server/handlers/users"
	"github.com/ivch/dynasty/server/middlewares"
)

var (
	defaultLogger *logger.StdLog
	errTestError  = errors.New("some err")
	defaultPolicy = bluemonday.StrictPolicy()
)

func TestMain(m *testing.M) {
	defaultLogger = logger.NewStdLog(logger.WithWriter(ioutil.Discard))
	os.Exit(m.Run())
}

func TestHTTP_Create(t *testing.T) {
	tests := []struct {
		name     string
		svc      RequestsService
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
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "error no user",
			request:  "{}",
			header:   "0",
			wantErr:  true,
			wantCode: http.StatusUnauthorized,
		},
		{
			name:     "error wrong type",
			request:  `{"time":1,"description":"abc", "type":"test"}`,
			header:   "1",
			wantErr:  true,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "error no time",
			request:  `{"type":"taxi","description":"abc"}`,
			header:   "1",
			wantErr:  true,
			wantCode: http.StatusBadRequest,
		},
		{
			name:    "error service",
			request: `{"type":"taxi","description":"abc","time":1}`,
			header:  "1",
			svc: &RequestsServiceMock{
				CreateFunc: func(_ context.Context, _ *requests.Request) (*requests.Request, error) {
					return nil, errTestError
				},
			},
			wantErr:  true,
			wantCode: http.StatusInternalServerError,
		},
		{
			name:    "ok",
			request: `{"type":"taxi","description":"abc","time":1}`,
			header:  "1",
			svc: &RequestsServiceMock{
				CreateFunc: func(_ context.Context, _ *requests.Request) (*requests.Request, error) {
					return &requests.Request{ID: 1}, nil
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
			h := NewHTTPTransport(defaultLogger, svc, defaultPolicy, middlewares.NewIDCtx(defaultLogger).Middleware)
			rr := httptest.NewRecorder()
			rq, _ := http.NewRequest("POST", "/v1/request", strings.NewReader(tt.request))
			rq.Header.Add("X-Auth-User", tt.header)
			h.ServeHTTP(rr, rq)
			if (rr.Code != tt.wantCode) && tt.wantErr {
				t.Errorf("Request error. status = %d, expected %v", rr.Code, tt.wantCode)
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
		svc      RequestsService
		id       string
		request  string
		header   string
		wantErr  bool
		wantCode int
	}{
		{
			name:     "error wrong user id",
			request:  "{}",
			id:       "1",
			header:   "asd",
			wantErr:  true,
			wantCode: http.StatusUnauthorized,
		},
		{
			name:     "error parsing request",
			request:  "bad json }}}",
			id:       "1",
			header:   "1",
			wantErr:  true,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "error no id",
			request:  "{}",
			id:       "0",
			header:   "1",
			wantErr:  true,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "error wrong id",
			request:  "{}",
			id:       "asd",
			header:   "1",
			wantErr:  true,
			wantCode: http.StatusBadRequest,
		},
		{
			name:    "error service",
			request: `{"type":"1","description":"abc","time":1}`,
			id:      "1",
			header:  "1",
			svc: &RequestsServiceMock{
				UpdateFunc: func(_ context.Context, _ *requests.Request) error {
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
			svc: &RequestsServiceMock{
				UpdateFunc: func(_ context.Context, r *requests.Request) error {
					expected := &requests.Request{
						ID:          1,
						UserID:      1,
						Type:        "1",
						Time:        1,
						Description: "abc",
						Status:      "new",
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
			h := NewHTTPTransport(defaultLogger, svc, defaultPolicy, middlewares.NewIDCtx(defaultLogger).Middleware)
			rr := httptest.NewRecorder()
			rq, _ := http.NewRequest("PUT", "/v1/request/"+tt.id, strings.NewReader(tt.request))
			rq.Header.Add("X-Auth-User", tt.header)
			h.ServeHTTP(rr, rq)
			if (rr.Code != tt.wantCode) && tt.wantErr {
				t.Errorf("Request error. status = %d, expected %v", rr.Code, tt.wantCode)
			}
		})
	}
}

func TestHTTP_My(t *testing.T) {
	tests := []struct {
		name     string
		svc      RequestsService
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
			wantCode: http.StatusUnauthorized,
		},
		{
			name:     "error wrong id",
			header:   "}{",
			wantErr:  true,
			wantCode: http.StatusUnauthorized,
		},
		{
			name:     "error no offset",
			header:   "1",
			query:    "?limit=1",
			wantErr:  true,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "error no limit",
			header:   "1",
			query:    "?offset=1",
			wantErr:  true,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "error bad offset",
			header:   "1",
			query:    "?offset=a&limit=1",
			wantErr:  true,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "error bad limit",
			header:   "1",
			query:    "?offset=1&limit=as",
			wantErr:  true,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "error too low limit",
			header:   "1",
			query:    "?offset=1&limit=0",
			wantErr:  true,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "error too big limit",
			header:   "1",
			query:    "?offset=1&limit=300",
			wantErr:  true,
			wantCode: http.StatusBadRequest,
		},
		{
			name:   "error service",
			query:  "?offset=1&limit=1",
			header: "1",
			svc: &RequestsServiceMock{
				MyFunc: func(_ context.Context, _ *requests.RequestListFilter) ([]*requests.Request, error) {
					return nil, errTestError
				},
			},
			wantErr:  true,
			wantCode: http.StatusInternalServerError,
		},
		{
			name:   "ok w/o images",
			query:  `?offset=0&limit=1`,
			header: "1",
			svc: &RequestsServiceMock{
				MyFunc: func(ctx context.Context, r *requests.RequestListFilter) (response []*requests.Request, err error) {
					return []*requests.Request{
						{
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
			want:     `{"data":[{"id":1,"type":"1","user_id":1,"time":1,"description":"1","status":"1"}]}`,
		},
		{
			name:   "ok w/ images",
			query:  `?offset=0&limit=1`,
			header: "1",
			svc: &RequestsServiceMock{
				MyFunc: func(ctx context.Context, r *requests.RequestListFilter) (response []*requests.Request, err error) {
					return []*requests.Request{
						{
							ID:          1,
							Type:        "1",
							UserID:      1,
							Time:        1,
							Description: "1",
							Status:      "1",
							ImagesURL: []map[string]string{
								{
									"img":   "a",
									"thumb": "b",
								},
							},
						},
					}, nil
				},
			},
			wantErr:  false,
			wantCode: http.StatusOK,
			want:     `{"data":[{"id":1,"type":"1","user_id":1,"time":1,"description":"1","status":"1","images":[{"img":"a","thumb":"b"}]}]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := tt.svc
			h := NewHTTPTransport(defaultLogger, svc, defaultPolicy, middlewares.NewIDCtx(defaultLogger).Middleware)
			rr := httptest.NewRecorder()
			rq, _ := http.NewRequest("GET", "/v1/my"+tt.query, nil)
			rq.Header.Add("X-Auth-User", tt.header)
			h.ServeHTTP(rr, rq)
			if (rr.Code != tt.wantCode) && tt.wantErr {
				t.Errorf("Request error. status = %d, expected %v", rr.Code, tt.wantCode)
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
		svc      RequestsService
		id       string
		header   string
		wantErr  bool
		wantCode int
	}{
		{
			name:     "error no user",
			id:       "1",
			header:   "0",
			wantErr:  true,
			wantCode: http.StatusUnauthorized,
		},
		{
			name:     "error no id",
			id:       "0",
			header:   "1",
			wantErr:  true,
			wantCode: http.StatusBadRequest,
		},
		{
			name:   "error service",
			id:     "1",
			header: "1",
			svc: &RequestsServiceMock{
				DeleteFunc: func(_ context.Context, _ *requests.Request) error {
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
			svc: &RequestsServiceMock{
				DeleteFunc: func(_ context.Context, r *requests.Request) error {
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
			h := NewHTTPTransport(defaultLogger, svc, defaultPolicy, middlewares.NewIDCtx(defaultLogger).Middleware)
			rr := httptest.NewRecorder()
			rq, _ := http.NewRequest("DELETE", "/v1/request/"+tt.id, nil)
			rq.Header.Add("X-Auth-User", tt.header)
			h.ServeHTTP(rr, rq)
			if (rr.Code != tt.wantCode) && tt.wantErr {
				t.Errorf("Request error. status = %d, expected %v", rr.Code, tt.wantCode)
			}
		})
	}
}

func TestHTTP_Get(t *testing.T) {
	tests := []struct {
		name     string
		svc      RequestsService
		id       string
		header   string
		wantErr  bool
		wantCode int
		want     string
	}{
		{
			name:     "error no user",
			id:       "1",
			header:   "0",
			wantErr:  true,
			wantCode: http.StatusUnauthorized,
		},
		{
			name:     "error no id",
			id:       "0",
			header:   "1",
			wantErr:  true,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "error wrong id",
			id:       "asd",
			header:   "1",
			wantErr:  true,
			wantCode: http.StatusBadRequest,
		},
		{
			name:   "error service",
			id:     "1",
			header: "1",
			svc: &RequestsServiceMock{
				GetFunc: func(_ context.Context, _ *requests.Request) (*requests.Request, error) {
					return nil, errTestError
				},
			},
			wantErr:  true,
			wantCode: http.StatusInternalServerError,
		},
		{
			name:   "ok w/o images",
			id:     "1",
			header: "1",
			svc: &RequestsServiceMock{
				GetFunc: func(_ context.Context, r *requests.Request) (*requests.Request, error) {
					if r.ID != 1 || r.UserID != 1 {
						return nil, errTestError
					}

					return &requests.Request{
						ID:          1,
						Type:        "1",
						UserID:      1,
						Time:        1,
						Description: "1",
						Status:      "1",
					}, nil
				},
			},
			wantErr:  false,
			wantCode: http.StatusOK,
			want:     `{"id":1,"type":"1","user_id":1,"time":1,"description":"1","status":"1"}`,
		},
		{
			name:   "ok w/ images",
			id:     "1",
			header: "1",
			svc: &RequestsServiceMock{
				GetFunc: func(_ context.Context, r *requests.Request) (*requests.Request, error) {
					if r.ID != 1 || r.UserID != 1 {
						return nil, errTestError
					}

					return &requests.Request{
						ID:          1,
						Type:        "1",
						UserID:      1,
						Time:        1,
						Description: "1",
						Status:      "1",
						ImagesURL: []map[string]string{
							{
								"img":   "a",
								"thumb": "b",
							},
						},
					}, nil
				},
			},
			wantErr:  false,
			wantCode: http.StatusOK,
			want:     `{"id":1,"type":"1","user_id":1,"time":1,"description":"1","status":"1","images":[{"img":"a","thumb":"b"}]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := tt.svc
			h := NewHTTPTransport(defaultLogger, svc, defaultPolicy, middlewares.NewIDCtx(defaultLogger).Middleware)
			rr := httptest.NewRecorder()
			rq, _ := http.NewRequest("GET", "/v1/request/"+tt.id, nil)
			rq.Header.Add("X-Auth-User", tt.header)
			h.ServeHTTP(rr, rq)
			if (rr.Code != tt.wantCode) && tt.wantErr {
				t.Errorf("Request error. status = %d, expected %v", rr.Code, tt.wantCode)
			}
		})
	}
}

func TestHTTP_GuardUpdate(t *testing.T) {
	tests := []struct {
		name     string
		svc      RequestsService
		id       string
		request  string
		wantErr  bool
		wantCode int
	}{
		{
			name:     "error parsing request",
			request:  "bad json }}}",
			id:       "1",
			wantErr:  true,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "error no id",
			request:  "{}",
			id:       "0",
			wantErr:  true,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "error wrong id",
			request:  "{}",
			id:       "asd",
			wantErr:  true,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "error bad request",
			request:  `{"status":"asd"}`,
			id:       "2",
			wantErr:  true,
			wantCode: http.StatusBadRequest,
		},
		{
			name:    "error service",
			request: `{"status":"closed"}`,
			id:      "1",
			svc: &RequestsServiceMock{
				GuardUpdateRequestFunc: func(_ context.Context, _ *requests.Request) error {
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
			svc: &RequestsServiceMock{
				GuardUpdateRequestFunc: func(_ context.Context, _ *requests.Request) error {
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
			h := NewHTTPTransport(defaultLogger, svc, defaultPolicy, middlewares.NewIDCtx(defaultLogger).Middleware)
			rr := httptest.NewRecorder()
			rq, _ := http.NewRequest("PUT", "/v1/guard/request/"+tt.id, strings.NewReader(tt.request))
			h.ServeHTTP(rr, rq)
			if (rr.Code != tt.wantCode) && tt.wantErr {
				t.Errorf("Request error. status = %d, expected %v", rr.Code, tt.wantCode)
			}
		})
	}
}

func TestHTTP_GuardList(t *testing.T) {
	tests := []struct {
		name     string
		svc      RequestsService
		query    string
		wantErr  bool
		wantCode int
		want     string
	}{
		{
			name:     "error no offset",
			query:    "?limit=1",
			wantErr:  true,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "error no limit",
			query:    "?offset=1",
			wantErr:  true,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "error bad offset",
			query:    "?offset=a&limit=1",
			wantErr:  true,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "error bad limit",
			query:    "?offset=1&limit=as",
			wantErr:  true,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "error too low limit",
			query:    "?offset=1&limit=0",
			wantErr:  true,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "error too big limit",
			query:    "?offset=1&limit=300",
			wantErr:  true,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "error bad type",
			query:    "?offset=1&limit=10&type=asd",
			wantErr:  true,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "error bad place",
			query:    "?offset=1&limit=10&place=asd",
			wantErr:  true,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "error bad status",
			query:    "?offset=1&limit=10&status=asd",
			wantErr:  true,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "error bad apartment",
			query:    "?offset=1&limit=10&apartment=asde",
			wantErr:  true,
			wantCode: http.StatusBadRequest,
		},
		{
			name:  "error service",
			query: "?offset=1&limit=1",
			svc: &RequestsServiceMock{
				GuardRequestListFunc: func(_ context.Context, _ *requests.RequestListFilter) ([]*requests.Request, int, error) {
					return nil, 0, errTestError
				},
			},
			wantErr:  true,
			wantCode: http.StatusInternalServerError,
		},
		{
			name:  "ok w/o images",
			query: `?offset=0&limit=1`,
			svc: &RequestsServiceMock{
				GuardRequestListFunc: func(_ context.Context, _ *requests.RequestListFilter) ([]*requests.Request, int, error) {
					return []*requests.Request{
						{
							ID:          1,
							Type:        "1",
							UserID:      1,
							Time:        1,
							Description: "1",
							Status:      "1",
							User:        &users.User{},
						},
					}, 1, nil
				},
			},
			wantErr:  false,
			wantCode: http.StatusOK,
			want:     `{"data":[{"id":1,"user_id":1,"type":"1","time":1,"description":"1","status":"1","user_name":" ","phone":"","address":"","apartment":0}],"count":1}`,
		},
		{
			name:  "ok w images",
			query: `?offset=0&limit=1`,
			svc: &RequestsServiceMock{
				GuardRequestListFunc: func(_ context.Context, _ *requests.RequestListFilter) ([]*requests.Request, int, error) {
					return []*requests.Request{
						{
							ID:          1,
							Type:        "1",
							UserID:      1,
							Time:        1,
							Description: "1",
							Status:      "1",
							User:        &users.User{},
							ImagesURL: []map[string]string{
								{
									"img":   "a",
									"thumb": "b",
								},
							},
						},
					}, 1, nil
				},
			},
			wantErr:  false,
			wantCode: http.StatusOK,
			want:     `{"data":[{"id":1,"user_id":1,"type":"1","time":1,"description":"1","status":"1","user_name":" ","phone":"","address":"","apartment":0,"images":[{"img":"a","thumb":"b"}]}],"count":1}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := tt.svc
			h := NewHTTPTransport(defaultLogger, svc, defaultPolicy, middlewares.NewIDCtx(defaultLogger).Middleware)
			rr := httptest.NewRecorder()
			rq, _ := http.NewRequest("GET", "/v1/guard/list"+tt.query, nil)
			h.ServeHTTP(rr, rq)
			if (rr.Code != tt.wantCode) && tt.wantErr {
				t.Errorf("Request error. status = %d, expected %v", rr.Code, tt.wantCode)
			}

			if !tt.wantErr && tt.want != strings.TrimSpace(rr.Body.String()) {
				t.Errorf("Response error, got = %s, want = %s", rr.Body.String(), tt.want)
			}
		})

	}
}

func TestHTTP_UploadImage(t *testing.T) {
	tests := []struct {
		name     string
		want     string
		header   string
		req      string
		filename string
		svc      RequestsService
		wantErr  bool
		wantCode int
	}{
		{
			name:     "error no user",
			req:      "",
			header:   "0",
			wantErr:  true,
			wantCode: http.StatusUnauthorized,
		},
		{
			name:     "error no id",
			req:      "",
			header:   "1",
			wantErr:  true,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "error bad id",
			req:      "asd",
			header:   "1",
			wantErr:  true,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "error parse multipart",
			req:      "1",
			header:   "1",
			wantErr:  true,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "error no file",
			req:      "1",
			header:   "1",
			wantErr:  true,
			filename: "errfile",
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "error service",
			req:      "1",
			header:   "1",
			wantErr:  true,
			filename: "ok",
			svc: &RequestsServiceMock{
				UploadImageFunc: func(_ context.Context, _ *requests.Image) (*requests.Image, error) {
					return nil, errTestError
				},
			},
			wantCode: http.StatusInternalServerError,
		},
		{
			name:     "ok",
			req:      "1",
			header:   "1",
			wantErr:  false,
			filename: "ok",
			svc: &RequestsServiceMock{
				UploadImageFunc: func(_ context.Context, _ *requests.Image) (*requests.Image, error) {
					return &requests.Image{
						URL:   "path",
						Thumb: "path",
					}, nil
				},
			},
			wantCode: http.StatusOK,
			want:     `{"img":"path","thumb":"path"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := tt.svc
			h := NewHTTPTransport(defaultLogger, svc, defaultPolicy, middlewares.NewIDCtx(defaultLogger).Middleware)
			rr := httptest.NewRecorder()

			bb := &bytes.Buffer{}
			contentHeader := "multipart/form-data"
			if tt.filename != "" {
				filename := "../../../../test_image.jpeg"
				f, err := os.Open(filename)
				if err != nil {
					t.Fatal(err)
				}

				paramName := "photo"
				if tt.filename == "errfile" {
					paramName = "err"
				}

				writer := multipart.NewWriter(bb)
				part, err := writer.CreateFormFile(paramName, filepath.Base(filename))
				if err != nil {
					t.Fatal(err)
				}
				_, err = io.Copy(part, f)
				if err != nil {
					t.Fatal(err)
				}
				contentHeader = writer.FormDataContentType()
				writer.Close()
			}

			rq, _ := http.NewRequest("POST", fmt.Sprintf("/v1/request/%s/file", tt.req), bb)
			rq.Header.Add("X-Auth-User", tt.header)
			rq.Header.Add("Content-Type", contentHeader)
			h.ServeHTTP(rr, rq)
			if (rr.Code != tt.wantCode) && tt.wantErr {
				t.Errorf("Request error. status = %d, expected %v", rr.Code, tt.wantCode)
			}

			if !tt.wantErr && tt.want != strings.TrimSpace(rr.Body.String()) {
				t.Errorf("Response error, got = %s, want = %s", rr.Body.String(), tt.want)
			}
		})
	}
}

func TestHTTP_DeleteImage(t *testing.T) {
	tests := []struct {
		name     string
		svc      RequestsService
		id       string
		request  string
		header   string
		wantErr  bool
		wantCode int
	}{
		{
			name:     "error no user",
			id:       "1",
			header:   "0",
			wantErr:  true,
			wantCode: http.StatusUnauthorized,
		},
		{
			name:     "error no id",
			id:       "0",
			header:   "1",
			wantErr:  true,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "error wrong id",
			id:       "asd",
			header:   "1",
			wantErr:  true,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "error bad request",
			id:       "1",
			header:   "1",
			request:  "}{",
			wantErr:  true,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "error empty request",
			id:       "1",
			header:   "1",
			request:  "{}",
			wantErr:  true,
			wantCode: http.StatusBadRequest,
		},
		{
			name:    "error service",
			id:      "1",
			header:  "1",
			request: `{"filepath":"somepath"}`,
			svc: &RequestsServiceMock{
				DeleteImageFunc: func(ctx context.Context, _ *requests.Image) error {
					return errTestError
				},
			},
			wantErr:  true,
			wantCode: http.StatusInternalServerError,
		},
		{
			name:    "ok",
			id:      "1",
			header:  "1",
			request: `{"filepath":"somepath"}`,
			svc: &RequestsServiceMock{
				DeleteImageFunc: func(ctx context.Context, r *requests.Image) error {
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
			h := NewHTTPTransport(defaultLogger, svc, defaultPolicy, middlewares.NewIDCtx(defaultLogger).Middleware)
			rr := httptest.NewRecorder()
			rq, _ := http.NewRequest("DELETE", "/v1/request/"+tt.id+"/file", strings.NewReader(tt.request))
			rq.Header.Add("X-Auth-User", tt.header)
			h.ServeHTTP(rr, rq)
			if (rr.Code != tt.wantCode) && tt.wantErr {
				t.Errorf("Request error. status = %d, expected %v", rr.Code, tt.wantCode)
			}
		})
	}
}
