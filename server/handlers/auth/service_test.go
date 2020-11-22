package auth

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/ivch/dynasty/common/logger"
	"github.com/ivch/dynasty/server/handlers/users/transport"
)

var (
	defaultLogger *logger.StdLog
	errTestError  = errors.New("some err")
)

func TestMain(m *testing.M) {
	defaultLogger = logger.NewStdLog(logger.WithWriter(ioutil.Discard))
	os.Exit(m.Run())
}

func TestService_Gwfa(t *testing.T) {
	var (
		invalidToken = `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MTksIk5hbWUiOiJhIGIiLCJSb2xlIjo0LCJhdWQiOiJkeW5hcHAiLCJleHAiOjE2MDU4MDYxNDcsImlhdCI6MTYwNTcxOTc0NywiaXNzIjoiYXV0aC5keW5hcHAifQ.5VFdapEz5DYdWJkBausjNL7vgVJXJ96KKHmlXsGgQy4`
		validToken   = `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MTAsIk5hbWUiOiJKYW5lIERvZSIsIlJvbGUiOjEwLCJhdWQiOiJkeW5hcHAiLCJleHAiOjc4ODg0NTQ4MDgsImlhdCI6MTU4MTI1NDgwOCwiaXNzIjoiYXV0aC5keW5hcHAifQ.Wzh9zsGSFEBVOWPLHzwfRcKiFQ9GJDbcFs8lE4X-Ha4`
		expired      = `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MTksIk5hbWUiOiJhIGIiLCJSb2xlIjo0LCJhdWQiOiJkeW5hcHAiLCJleHAiOjE2MDU2MzM1NjUsImlhdCI6MTYwNTcxOTk2NSwiaXNzIjoiYXV0aC5keW5hcHAifQ.x23MsR9bxzuxUcMUNajgX8MA085fO1yhcUKEUrYSr8Y`
	)

	tests := []struct {
		name    string
		secret  string
		token   string
		wantErr bool
		want    uint
	}{
		{
			name:    "error failed to parse token",
			token:   "asd",
			wantErr: true,
		},
		{
			name:    "error invalid token signature",
			token:   invalidToken,
			secret:  "covabunga",
			wantErr: true,
		},
		{
			name:    "error expired token",
			token:   expired,
			secret:  "covabunga",
			wantErr: true,
		},
		{
			name:    "ok",
			token:   validToken,
			secret:  "covabunga",
			wantErr: false,
			want:    10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(defaultLogger, nil, nil, tt.secret)
			got, err := s.Gwfa(tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("Gwfa() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Gwfa() got = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestService_Refresh(t *testing.T) {
	type fields struct {
		repo authRepository
		usrv userService
		req  string
	}

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
		want    *Tokens
	}{
		{
			name: "error no session",
			fields: fields{
				repo: &authRepositoryMock{
					FindSessionByAccessTokenFunc: func(t string) (*Session, error) {
						return nil, errTestError
					},
				},
				req: "token",
			},
			wantErr: true,
		},
		{
			name: "error on session delete",
			fields: fields{
				repo: &authRepositoryMock{
					FindSessionByAccessTokenFunc: func(t string) (*Session, error) {
						return &Session{ID: "1"}, nil
					},
					DeleteSessionByIDFunc: func(_ string) error {
						return errTestError
					},
				},
				req: "token",
			},
			wantErr: true,
		},
		{
			name: "error expired token",
			fields: fields{
				repo: &authRepositoryMock{
					FindSessionByAccessTokenFunc: func(t string) (*Session, error) {
						return &Session{
							ID:        "1",
							ExpiresIn: time.Now().Add(-1 * time.Minute).Unix(),
						}, nil
					},
					DeleteSessionByIDFunc: func(_ string) error {
						return nil
					},
				},
				req: "token",
			},
			wantErr: true,
		},
		{
			name: "error finding user",
			fields: fields{
				repo: &authRepositoryMock{
					FindSessionByAccessTokenFunc: func(t string) (*Session, error) {
						return &Session{
							ID:        "1",
							ExpiresIn: time.Now().Add(10 * time.Minute).Unix(),
						}, nil
					},
					DeleteSessionByIDFunc: func(_ string) error {
						return nil
					},
				},
				usrv: &userServiceMock{
					UserByIDFunc: func(_ context.Context, _ uint) (*transport.UserByIDResponse, error) {
						return nil, errTestError
					},
				},
				req: "token",
			},
			wantErr: true,
		},
		{
			name: "error creating session",
			fields: fields{
				repo: &authRepositoryMock{
					FindSessionByAccessTokenFunc: func(t string) (*Session, error) {
						return &Session{
							ID:        "1",
							ExpiresIn: time.Now().Add(10 * time.Minute).Unix(),
							UserID:    1,
						}, nil
					},
					DeleteSessionByIDFunc: func(_ string) error {
						return nil
					},
					CreateSessionFunc: func(_ uint) (string, error) {
						return "", errTestError
					},
				},
				usrv: &userServiceMock{
					UserByIDFunc: func(_ context.Context, _ uint) (*transport.UserByIDResponse, error) {
						return nil, nil
					},
				},
				req: "token",
			},
			wantErr: true,
		},
		{
			name: "ok",
			fields: fields{
				repo: &authRepositoryMock{
					FindSessionByAccessTokenFunc: func(t string) (*Session, error) {
						return &Session{
							ID:        "1",
							ExpiresIn: time.Now().Add(10 * time.Minute).Unix(),
							UserID:    1,
						}, nil
					},
					DeleteSessionByIDFunc: func(_ string) error {
						return nil
					},
					CreateSessionFunc: func(_ uint) (string, error) {
						return "refresh_token", nil
					},
				},
				usrv: &userServiceMock{
					UserByIDFunc: func(_ context.Context, _ uint) (*transport.UserByIDResponse, error) {
						return &transport.UserByIDResponse{
							ID:        1,
							FirstName: "Jane",
							LastName:  "Doe",
							Role:      1,
						}, nil
					},
				},
				req: "token",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(defaultLogger, tt.fields.repo, tt.fields.usrv, "secret")
			got, err := s.Refresh(context.Background(), tt.fields.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Refresh() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil && (got.RefreshToken == "" || got.AccessToken == "") {
				t.Errorf("Refresh() empty resilt %v", got)
			}
		})
	}
}

func TestService_Login(t *testing.T) {
	type fields struct {
		usrv userService
		repo authRepository
		req  [2]string
	}

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "error no user",
			fields: fields{
				usrv: &userServiceMock{
					UserByPhoneAndPasswordFunc: func(_ context.Context, _ string, _ string) (*transport.UserByIDResponse, error) {
						return nil, errTestError
					},
				},
				req: [2]string{"1", "2"},
			},
			wantErr: true,
		},
		{
			name: "error on inactive user",
			fields: fields{
				usrv: &userServiceMock{
					UserByPhoneAndPasswordFunc: func(_ context.Context, _ string, _ string) (*transport.UserByIDResponse, error) {
						return &transport.UserByIDResponse{ID: 1, Active: false}, nil
					},
				},
				repo: &authRepositoryMock{
					CreateSessionFunc: func(_ uint) (string, error) {
						return "", errTestError
					},
				},
				req: [2]string{"1", "2"},
			},
			wantErr: true,
		},
		{
			name: "error on create session",
			fields: fields{
				usrv: &userServiceMock{
					UserByPhoneAndPasswordFunc: func(_ context.Context, _ string, _ string) (*transport.UserByIDResponse, error) {
						return &transport.UserByIDResponse{ID: 1, Active: true}, nil
					},
				},
				repo: &authRepositoryMock{
					CreateSessionFunc: func(_ uint) (string, error) {
						return "", errTestError
					},
				},
				req: [2]string{"1", "2"},
			},
			wantErr: true,
		},
		{
			name: "ok",
			fields: fields{
				usrv: &userServiceMock{
					UserByPhoneAndPasswordFunc: func(_ context.Context, _ string, _ string) (*transport.UserByIDResponse, error) {
						return &transport.UserByIDResponse{ID: 1, Role: 1, FirstName: "Jane", LastName: "Doe", Active: true}, nil
					},
				},
				repo: &authRepositoryMock{
					CreateSessionFunc: func(_ uint) (string, error) {
						return "token", nil
					},
				},
				req: [2]string{"1", "2"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(defaultLogger, tt.fields.repo, tt.fields.usrv, "secret")
			_, err := s.Login(context.Background(), tt.fields.req[0], tt.fields.req[1])
			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
