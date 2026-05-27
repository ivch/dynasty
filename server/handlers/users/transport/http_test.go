package transport_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/microcosm-cc/bluemonday"

	"github.com/ivch/dynasty/common/errs"
	"github.com/ivch/dynasty/common/logger"
	"github.com/ivch/dynasty/server/handlers/users"
	"github.com/ivch/dynasty/server/handlers/users/transport"
	"github.com/ivch/dynasty/server/middlewares"
)

var (
	defaultLogger *logger.StdLog
	errTestError  = errors.New("some err")
	defaultPolicy = bluemonday.StrictPolicy()
)

func TestMain(m *testing.M) {
	defaultLogger = logger.NewStdLog(logger.WithWriter(io.Discard))
	os.Exit(m.Run())
}

func TestHTTP_GetUser(t *testing.T) {
	tests := []struct {
		name     string
		svc      transport.UsersService
		header   string
		wantErr  bool
		want     string
		wantCode int
	}{
		{
			name:     "error no id",
			svc:      nil,
			header:   "0",
			wantErr:  true,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "error wrong id",
			svc:      nil,
			header:   "a",
			wantErr:  true,
			wantCode: http.StatusBadRequest,
		},
		{
			name: "error service error",
			svc: &transport.UsersServiceMock{
				UserByIDFunc: func(_ context.Context, _ uint) (*users.User, error) {
					return nil, errTestError
				},
			},
			header:   "1",
			wantErr:  true,
			wantCode: http.StatusInternalServerError,
		},
		{
			name: "error data not found",
			svc: &transport.UsersServiceMock{
				UserByIDFunc: func(_ context.Context, _ uint) (*users.User, error) {
					return nil, nil
				},
			},
			header:   "1",
			wantErr:  true,
			wantCode: http.StatusNotFound,
		},
		{
			name: "ok",
			svc: &transport.UsersServiceMock{
				UserByIDFunc: func(_ context.Context, _ uint) (*users.User, error) {
					return &users.User{
						ID:        1,
						Apartment: 1,
						FirstName: "1",
						LastName:  "1",
						Phone:     "1",
						Email:     "1",
						Role:      1,
						Active:    true,
						Building: users.Building{
							ID:      1,
							Name:    "1",
							Address: "1",
						},
					}, nil
				},
			},
			header:   "1",
			wantErr:  false,
			want:     `{"id":1,"apartment":1,"first_name":"1","last_name":"1","phone":"1","email":"1","role":1,"building":{"id":1,"name":"1","address":"1"},"entry":{"id":0,"name":"","building_id":0},"active":true}`,
			wantCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := tt.svc
			h := transport.NewHTTPTransport(defaultLogger, svc, defaultPolicy, middlewares.NewIDCtx(defaultLogger).Middleware)
			rr := httptest.NewRecorder()
			rq, _ := http.NewRequest(http.MethodGet, "/v1/user", nil)
			rq.Header.Add("X-Auth-User", tt.header)
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

func TestHTTP_Register(t *testing.T) {
	tests := []struct {
		name    string
		svc     transport.UsersService
		request string
		wantErr bool
		want    string
	}{
		{
			name:    "error unmarshal request",
			svc:     nil,
			request: "}",
			wantErr: true,
		},
		{
			name:    "error invalid password",
			svc:     nil,
			request: `{"email":"test@test.com","first_name":"John","last_name":"Doe","apartment":1,"phone":"380671234567","building_id": 2, "code":"1231"}`,
			wantErr: true,
		},
		{
			name:    "error password too short",
			svc:     nil,
			request: `{"email":"test@test.com","first_name":"John","last_name":"Doe","apartment":1,"password":"1213", "phone":"380671234567","building_id": 2, "code":"1231"}`,
			wantErr: true,
		},
		{
			name:    "error to short phone",
			svc:     nil,
			request: `{"email":"test@test.com","first_name":"John","phone":"12345","last_name":"Doe","apartment":1,"password":"1213123","building_id": 2, "code":"1231"}`,
			wantErr: true,
		},
		{
			name:    "error wrong phone",
			svc:     nil,
			request: `{"email":"test@test.com","first_name":"John","phone":"123asd123asd","last_name":"Doe","apartment":1,"password":"1213123","building_id": 2, "code":"1231"}`,
			wantErr: true,
		},
		{
			name:    "error invalid first name",
			svc:     nil,
			request: `{"email":"test@test.com","last_name":"Doe","apartment":1,"password":"1213123", "phone":"380671234567","building_id": 2, "code":"1231"}`,
			wantErr: true,
		},
		{
			name:    "error invalid last name",
			svc:     nil,
			request: `{"email":"test@test.com","first_name":"John","apartment":1,"password":"1213123", "phone":"380671234567","building_id": 2, "code":"1231"}`,
			wantErr: true,
		},
		{
			name:    "error invalid building",
			svc:     nil,
			request: `{"email":"test@test.com","first_name":"John","last_name":"Doe","apartment":1,"password":"1213123", "phone":"380671234567", "code":"1231"}`,
			wantErr: true,
		},
		{
			name:    "error invalid entry",
			svc:     nil,
			request: `{"email":"test@test.com","first_name":"John","last_name":"Doe","password":"1213123", "phone":"380671234567","building_id": 2,"code":"1231"}`,
			wantErr: true,
		},
		{
			name:    "error invalid apartment",
			svc:     nil,
			request: `{"email":"test@test.com","first_name":"John","last_name":"Doe","password":"1213123", "phone":"380671234567","building_id": 2, "entry_id": 1, "code":"1231"}`,
			wantErr: true,
		},
		{
			name:    "error invalid apartment #2",
			svc:     nil,
			request: `{"email":"test@test.com","first_name":"John","last_name":"Doe","password":"1213123", "phone":"380671234567","building_id": 2, "entry_id": 1, "code":"1231", "apartment":12312312"}`,
			wantErr: true,
		},
		{
			name:    "error invalid email",
			svc:     nil,
			request: `{"first_name":"John","last_name":"Doe","apartment":1, "entry_id": 1,"password":"1213123", "phone":"380671234567","building_id": 2, "code":"1231"}`,
			wantErr: true,
		},
		{
			name:    "error invalid email#2",
			svc:     nil,
			request: `{"email":"testst.com","first_name":"John","last_name":"Doe","apartment":1, "entry_id": 1,"password":"1213123", "phone":"380671234567","building_id": 2, "code":"1231"}`,
			wantErr: true,
		},
		{
			name:    "error invalid email#3",
			svc:     nil,
			request: `{"email":"te","first_name":"John","last_name":"Doe","apartment":1, "entry_id": 1,"password":"1213123", "phone":"380671234567","building_id": 2, "code":"1231"}`,
			wantErr: true,
		},
		{
			name:    "error invalid email#4",
			svc:     nil,
			request: `{"email":"test@test.com","first_name":"John","last_name":"Doe","apartment":1, "entry_id": 1,"password":"1213123", "phone":"380671234567","building_id": 2, "code":"1231"}`,
			wantErr: true,
		},
		{
			name: "error service",
			svc: &transport.UsersServiceMock{
				RegisterFunc: func(_ context.Context, _ *users.User) (*users.User, error) {
					return nil, errTestError
				},
			},
			request: `{"email":"test@mail.com","first_name":"John","last_name":"Doe","apartment":1, "entry_id": 1,"password":"1213123", "phone":"380671234567","building_id": 2, "code":"1231"}`,
			wantErr: true,
		},
		{
			name: "ok",
			svc: &transport.UsersServiceMock{
				RegisterFunc: func(_ context.Context, _ *users.User) (*users.User, error) {
					return &users.User{
						ID:    1,
						Phone: "380671234567",
					}, nil
				},
			},
			request: `{"email":"test@mail.com","first_name":"John", "entry_id": 1,"last_name":"Doe","apartment":1,"password":"1213123", "phone":"+380671234567","building_id": 2, "code":"1231", "entry_id":1}`,
			wantErr: false,
			want:    `{"id":1,"phone":"380671234567"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := tt.svc
			h := transport.NewHTTPTransport(defaultLogger, svc, defaultPolicy, middlewares.NewIDCtx(defaultLogger).Middleware)
			rr := httptest.NewRecorder()
			rq, _ := http.NewRequest(http.MethodPost, "/v1/register", strings.NewReader(tt.request))
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

func TestHTTP_AddFamilyMember(t *testing.T) {
	tests := []struct {
		name    string
		svc     transport.UsersService
		header  string
		request string
		wantErr bool
		want    string
	}{
		{
			name:    "error no user",
			svc:     nil,
			header:  "",
			wantErr: true,
		},
		{
			name:    "error unmarshal request",
			svc:     nil,
			header:  "1",
			request: "}",
			wantErr: true,
		},
		{
			name:    "error empty phone",
			svc:     nil,
			header:  "1",
			request: `{"phone":""}`,
			wantErr: true,
		},
		{
			name:    "error bad phone",
			svc:     nil,
			header:  "1",
			request: `{"phone":"123asd123asd"}`,
			wantErr: true,
		},
		{
			name:   "error service",
			header: "1",
			svc: &transport.UsersServiceMock{
				AddFamilyMemberFunc: func(_ context.Context, _ *users.User) (*users.User, error) {
					return nil, errTestError
				},
			},
			request: `{"phone":"123456789012"}`,
			wantErr: true,
		},
		{
			name:   "ok",
			header: "1",
			svc: &transport.UsersServiceMock{
				AddFamilyMemberFunc: func(_ context.Context, r *users.User) (*users.User, error) {
					return &users.User{RegCode: "123"}, nil
				},
			},
			request: `{"phone":"380671234567"}`,
			wantErr: false,
			want:    `{"code":"123"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := tt.svc
			h := transport.NewHTTPTransport(defaultLogger, svc, defaultPolicy, middlewares.NewIDCtx(defaultLogger).Middleware)
			rr := httptest.NewRecorder()
			rq, _ := http.NewRequest(http.MethodPost, "/v1/member", strings.NewReader(tt.request))
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

func TestHTTP_ListFamilyMembers(t *testing.T) {
	tests := []struct {
		name    string
		svc     transport.UsersService
		header  string
		wantErr bool
		want    string
	}{
		{
			name:    "error no user",
			svc:     nil,
			header:  "",
			wantErr: true,
		},
		{
			name:   "error service",
			header: "1",
			svc: &transport.UsersServiceMock{
				ListFamilyMembersFunc: func(_ context.Context, _ uint) ([]*users.User, error) {
					return nil, errTestError
				},
			},
			wantErr: true,
		},
		{
			name:   "ok",
			header: "1",
			svc: &transport.UsersServiceMock{
				ListFamilyMembersFunc: func(_ context.Context, _ uint) ([]*users.User, error) {
					return []*users.User{
						{
							ID:        1,
							FirstName: "1",
							LastName:  "2",
							Phone:     "1",
							RegCode:   "1",
							Active:    false,
						},
					}, nil
				},
			},
			wantErr: false,
			want:    `{"data":[{"id":1,"name":"1 2","phone":"1","code":"1","active":false}]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := tt.svc
			h := transport.NewHTTPTransport(defaultLogger, svc, defaultPolicy, middlewares.NewIDCtx(defaultLogger).Middleware)
			rr := httptest.NewRecorder()
			rq, _ := http.NewRequest(http.MethodGet, "/v1/members", nil)
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

func TestHTTP_DeleteFamilyMember(t *testing.T) {
	tests := []struct {
		name    string
		svc     transport.UsersService
		header  string
		request string
		wantErr bool
	}{
		{
			name:    "error no user",
			svc:     nil,
			header:  "",
			request: "1",
			wantErr: true,
		},
		{
			name:    "error bad request",
			svc:     nil,
			header:  "1",
			request: "asd",
			wantErr: true,
		},
		{
			name:    "error zero member",
			svc:     nil,
			header:  "1",
			request: "0",
			wantErr: true,
		},
		{
			name:    "error same ids",
			svc:     nil,
			header:  "1",
			request: "1",
			wantErr: true,
		},
		{
			name:    "error service",
			header:  "1",
			request: "2",
			svc: &transport.UsersServiceMock{
				DeleteFamilyMemberFunc: func(_ context.Context, _, _ uint) error {
					return errTestError
				},
			},
			wantErr: true,
		},
		{
			name:    "ok",
			header:  "1",
			request: "2",
			svc: &transport.UsersServiceMock{
				DeleteFamilyMemberFunc: func(_ context.Context, _, _ uint) error {
					return nil
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := tt.svc
			h := transport.NewHTTPTransport(defaultLogger, svc, defaultPolicy, middlewares.NewIDCtx(defaultLogger).Middleware)
			rr := httptest.NewRecorder()
			rq, _ := http.NewRequest(http.MethodDelete, "/v1/member/"+tt.request, nil)
			rq.Header.Add("X-Auth-User", tt.header)
			h.ServeHTTP(rr, rq)
			if (rr.Code != http.StatusOK) != tt.wantErr {
				t.Errorf("Request error. status = %d, wantErr %v", rr.Code, tt.wantErr)
			}
		})
	}
}

func TestHTTP_UpdateUser(t *testing.T) {
	tests := []struct {
		name    string
		svc     transport.UsersService
		header  string
		request string
		wantErr bool
	}{
		{
			name:    "error no user",
			svc:     nil,
			header:  "",
			request: "1",
			wantErr: true,
		},
		{
			name:    "error bad request",
			svc:     nil,
			header:  "1",
			request: "asd",
			wantErr: true,
		},
		{
			name: "email invalid",
			svc: &transport.UsersServiceMock{
				UpdateFunc: func(_ context.Context, req *users.UserUpdate) error {
					if *req.Email != "email" {
						return errTestError
					}
					return nil
				},
			},
			header:  "1",
			request: `{"email":"asd"}`,
			wantErr: true,
		},
		{
			name: "email not updated",
			svc: &transport.UsersServiceMock{
				UpdateFunc: func(_ context.Context, req *users.UserUpdate) error {
					if *req.Email != "email@mail.com" {
						return errTestError
					}
					return nil
				},
			},
			header:  "1",
			request: `{"email":"test@mail.com"}`,
			wantErr: true,
		},
		{
			name: "email updated",
			svc: &transport.UsersServiceMock{
				UpdateFunc: func(_ context.Context, req *users.UserUpdate) error {
					if *req.Email != "email@mail.com" {
						return errTestError
					}
					return nil
				},
			},
			header:  "1",
			request: `{"email":"email@mail.com"}`,
			wantErr: false,
		},
		{
			name: "first name not updated",
			svc: &transport.UsersServiceMock{
				UpdateFunc: func(_ context.Context, req *users.UserUpdate) error {
					if *req.FirstName != "a" {
						return errTestError
					}
					return nil
				},
			},
			header:  "1",
			request: `{"first_name":"b"}`,
			wantErr: true,
		},
		{
			name: "first name updated",
			svc: &transport.UsersServiceMock{
				UpdateFunc: func(_ context.Context, req *users.UserUpdate) error {
					if *req.FirstName != "a" {
						return errTestError
					}
					return nil
				},
			},
			header:  "1",
			request: `{"first_name":"a"}`,
			wantErr: false,
		},
		{
			name: "last name not updated",
			svc: &transport.UsersServiceMock{
				UpdateFunc: func(_ context.Context, req *users.UserUpdate) error {
					if *req.LastName != "a" {
						return errTestError
					}
					return nil
				},
			},
			header:  "1",
			request: `{"last_name":"b"}`,
			wantErr: true,
		},
		{
			name: "last name updated",
			svc: &transport.UsersServiceMock{
				UpdateFunc: func(_ context.Context, req *users.UserUpdate) error {
					if *req.LastName != "a" {
						return errTestError
					}
					return nil
				},
			},
			header:  "1",
			request: `{"last_name":"a"}`,
			wantErr: false,
		},
		{
			name:    "error bad password",
			svc:     nil,
			header:  "1",
			request: `{"new_password":"a"}`,
			wantErr: true,
		},
		{
			name:    "error no password confirm",
			svc:     nil,
			header:  "1",
			request: `{"new_password":"a", "password":"1"}`,
			wantErr: true,
		},
		{
			name:    "error password mismatch",
			svc:     nil,
			header:  "1",
			request: `{"new_password":"a", "password":"1", "new_password_confirm":"b"}`,
			wantErr: true,
		},
		{
			name:    "error invalid password",
			svc:     nil,
			header:  "1",
			request: `{"new_password":"a", "password":"1", "new_password_confirm":"a"}`,
			wantErr: true,
		},
		{
			name: "error password not updated",
			svc: &transport.UsersServiceMock{
				UpdateFunc: func(_ context.Context, req *users.UserUpdate) error {
					if *req.NewPassword != "1234567" {
						return errTestError
					}
					return nil
				},
			},
			header:  "1",
			request: `{"new_password":"1234567", "password":"1", "new_password_confirm":"1234567"}`,
			wantErr: false,
		},
		{
			name:    "error service",
			header:  "1",
			request: "{}",
			svc: &transport.UsersServiceMock{
				UpdateFunc: func(_ context.Context, _ *users.UserUpdate) error {
					return errTestError
				},
			},
			wantErr: true,
		},
		{
			name:    "ok",
			header:  "1",
			request: "{}",
			svc: &transport.UsersServiceMock{
				UpdateFunc: func(_ context.Context, _ *users.UserUpdate) error {
					return nil
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := tt.svc
			h := transport.NewHTTPTransport(defaultLogger, svc, defaultPolicy, middlewares.NewIDCtx(defaultLogger).Middleware)
			rr := httptest.NewRecorder()
			rq, _ := http.NewRequest(http.MethodPut, "/v1/user", strings.NewReader(tt.request))
			rq.Header.Add("X-Auth-User", tt.header)
			h.ServeHTTP(rr, rq)
			if (rr.Code != http.StatusOK) != tt.wantErr {
				t.Errorf("Request error. status = %d, wantErr %v", rr.Code, tt.wantErr)
			}
		})
	}
}

func TestHTTPTransport_PasswordRecoveryRequest(t *testing.T) {
	tests := []struct {
		name    string
		svc     transport.UsersService
		request string
		wantErr bool
	}{
		{
			name:    "error bad json",
			svc:     nil,
			request: "}{",
			wantErr: true,
		},
		{
			name:    "error empty phone",
			svc:     nil,
			request: `{"phone":""}`,
			wantErr: true,
		},
		{
			name:    "error bad phone",
			svc:     nil,
			request: `{"phone":"123asd123asd"}`,
			wantErr: true,
		},
		{
			name:    "email invalid",
			svc:     nil,
			request: `{"email":"12", "phone":"123456789012"}`,
			wantErr: true,
		},
		{
			name:    "email invalid",
			svc:     nil,
			request: `{"email":"test@test.com", "phone":"123456789012"}`,
			wantErr: true,
		},
		{
			name:    "error service",
			request: `{"email":"test@mail.com", "phone":"123456789012"}`,
			svc: &transport.UsersServiceMock{
				RecoveryCodeFunc: func(_ context.Context, _ *users.User) error {
					return errTestError
				},
			},
			wantErr: true,
		},
		{
			name:    "ok",
			request: `{"email":"test@mail.com", "phone":"123456789012"}`,
			svc: &transport.UsersServiceMock{
				RecoveryCodeFunc: func(_ context.Context, _ *users.User) error {
					return nil
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := tt.svc
			h := transport.NewHTTPTransport(defaultLogger, svc, defaultPolicy, middlewares.NewIDCtx(defaultLogger).Middleware)
			rr := httptest.NewRecorder()
			rq, _ := http.NewRequest(http.MethodPost, "/v1/password-recovery", strings.NewReader(tt.request))
			h.ServeHTTP(rr, rq)
			if (rr.Code != http.StatusOK) != tt.wantErr {
				t.Errorf("Request error. status = %d, wantErr %v", rr.Code, tt.wantErr)
			}
		})
	}
}

func TestHTTPTransport_PasswordResetRequest(t *testing.T) {
	tests := []struct {
		name    string
		svc     transport.UsersService
		request string
		wantErr bool
	}{
		{
			name:    "error bad json",
			svc:     nil,
			request: "}{",
			wantErr: true,
		},
		{
			name:    "error short code",
			svc:     nil,
			request: `{"code":"1"}`,
			wantErr: true,
		},
		{
			name:    "error no new password",
			svc:     nil,
			request: `{"code":"1234567890"}`,
			wantErr: true,
		},
		{
			name:    "error no new password confirm",
			svc:     nil,
			request: `{"code":"1234567890","new_password":"passwd"}`,
			wantErr: true,
		},
		{
			name:    "error password confirm mismatch",
			svc:     nil,
			request: `{"code":"1234567890","new_password":"passwd", "new_password_confirm":"passwd1"}`,
			wantErr: true,
		},
		{
			name:    "error bad password",
			svc:     nil,
			request: `{"code":"1234567890","new_password":"pass", "new_password_confirm":"pass"}`,
			wantErr: true,
		},
		{
			name:    "error service",
			request: `{"code":"1234567890","new_password":"passwd", "new_password_confirm":"passwd"}`,
			svc: &transport.UsersServiceMock{
				ResetPasswordFunc: func(_ context.Context, _ string, _ *users.UserUpdate) error {
					return errTestError
				},
			},
			wantErr: true,
		},
		{
			name:    "ok",
			request: `{"code":"1234567890","new_password":"passwd", "new_password_confirm":"passwd"}`,
			svc: &transport.UsersServiceMock{
				ResetPasswordFunc: func(_ context.Context, _ string, _ *users.UserUpdate) error {
					return nil
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := tt.svc
			h := transport.NewHTTPTransport(defaultLogger, svc, defaultPolicy, middlewares.NewIDCtx(defaultLogger).Middleware)
			rr := httptest.NewRecorder()
			rq, _ := http.NewRequest(http.MethodPost, "/v1/password-reset", strings.NewReader(tt.request))
			h.ServeHTTP(rr, rq)
			if (rr.Code != http.StatusOK) != tt.wantErr {
				t.Errorf("Request error. status = %d, wantErr %v", rr.Code, tt.wantErr)
			}
		})
	}
}

func TestHTTP_AdminResetApartment(t *testing.T) {
	tests := []struct {
		name     string
		svc      transport.UsersService
		header   string
		request  string
		wantCode int
	}{
		{
			name:     "missing admin header",
			svc:      nil,
			header:   "",
			request:  `{"building_id":1,"apartment_number":123}`,
			wantCode: http.StatusUnauthorized,
		},
		{
			name:     "bad admin id in header",
			svc:      nil,
			header:   "0",
			request:  `{"building_id":1,"apartment_number":123}`,
			wantCode: http.StatusUnauthorized,
		},
		{
			name:     "missing building_id",
			svc:      nil,
			header:   "1",
			request:  `{"apartment_number":123}`,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "missing apartment_number",
			svc:      nil,
			header:   "1",
			request:  `{"building_id":1}`,
			wantCode: http.StatusBadRequest,
		},
		{
			name:    "non-admin forbidden",
			header:  "1",
			request: `{"building_id":1,"apartment":123}`,
			svc: &transport.UsersServiceMock{
				AdminResetApartmentFunc: func(_ context.Context, _, _, _ uint) (string, error) {
					return "", errs.InsufficientPermissions
				},
			},
			wantCode: http.StatusForbidden,
		},
		{
			name:    "service error",
			header:  "1",
			request: `{"building_id":1,"apartment":123}`,
			svc: &transport.UsersServiceMock{
				AdminResetApartmentFunc: func(_ context.Context, _, _, _ uint) (string, error) {
					return "", errTestError
				},
			},
			wantCode: http.StatusInternalServerError,
		},
		{
			name:    "ok",
			header:  "1",
			request: `{"building_id":1,"apartment":123}`,
			svc: &transport.UsersServiceMock{
				AdminResetApartmentFunc: func(_ context.Context, _, _, _ uint) (string, error) {
					return "abc123", nil
				},
			},
			wantCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := transport.NewHTTPTransport(defaultLogger, tt.svc, defaultPolicy, middlewares.NewIDCtx(defaultLogger).Middleware)
			rr := httptest.NewRecorder()
			rq, _ := http.NewRequest(http.MethodPost, "/v1/admin/apartment/reset", strings.NewReader(tt.request))
			if tt.header != "" {
				rq.Header.Set("X-Auth-User", tt.header)
			}
			h.ServeHTTP(rr, rq)
			if rr.Code != tt.wantCode {
				t.Errorf("AdminResetApartment() status = %d, want %d", rr.Code, tt.wantCode)
			}
			if tt.name == "ok" && !strings.Contains(rr.Body.String(), `"reg_code":"abc123"`) {
				t.Errorf("AdminResetApartment() body = %s, want reg_code in response", rr.Body.String())
			}
		})
	}
}
