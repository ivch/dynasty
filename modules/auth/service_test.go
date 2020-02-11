package auth

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"

	"github.com/ivch/dynasty/models/dto"
	"github.com/ivch/dynasty/models/entities"
)

var (
	defaultLogger *zerolog.Logger
	errTestError  = errors.New("some err")
)

func TestMain(m *testing.M) {
	logger := zerolog.New(ioutil.Discard)
	defaultLogger = &logger
	os.Exit(m.Run())
}

func TestService_Gwfa(t *testing.T) {
	var (
		invalidToken = `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MTAsIk5hbWUiOiJKYW5lIERvZSIsIlJvbGUiOjEwLCJhdWQiOiJkeW5hcHAiLCJleHAiOjE1ODEzNDAyODYsImlhdCI6MTU4MTI1Mzg4NiwiaXNzIjoiYXV0aC5keW5hcHAifQ.VrX0Bih8Ha7HHlrRjvKD_mJmH_ND4EPKZ1HuSoeZWsg`
		validToken   = `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MTAsIk5hbWUiOiJKYW5lIERvZSIsIlJvbGUiOjEwLCJhdWQiOiJkeW5hcHAiLCJleHAiOjc4ODg0NTQ4MDgsImlhdCI6MTU4MTI1NDgwOCwiaXNzIjoiYXV0aC5keW5hcHAifQ.Wzh9zsGSFEBVOWPLHzwfRcKiFQ9GJDbcFs8lE4X-Ha4`
		expired      = `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MTAsIk5hbWUiOiJKYW5lIERvZSIsIlJvbGUiOjEwLCJhdWQiOiJkeW5hcHAiLCJleHAiOjE1ODExNjc5OTAsImlhdCI6MTU4MTI1NDM5MCwiaXNzIjoiYXV0aC5keW5hcHAifQ.axsbEPTx7sG1cB3eFIbScrT1827etn5sIyLGpL0Lwq8`
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
			s := newService(defaultLogger, nil, nil, tt.secret)
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
		req  *dto.AuthRefreshTokenRequest
	}

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
		want    *dto.AuthLoginResponse
	}{
		{
			name: "error no session",
			fields: fields{
				repo: &authRepositoryMock{
					FindSessionByAccessTokenFunc: func(t string) (*entities.Session, error) {
						return nil, errTestError
					},
				},
				req: &dto.AuthRefreshTokenRequest{Token: "token"},
			},
			wantErr: true,
		},
		{
			name: "error on session delete",
			fields: fields{
				repo: &authRepositoryMock{
					FindSessionByAccessTokenFunc: func(t string) (*entities.Session, error) {
						return &entities.Session{ID: "1"}, nil
					},
					DeleteSessionByIDFunc: func(_ string) error {
						return errTestError
					},
				},
				req: &dto.AuthRefreshTokenRequest{Token: "token"},
			},
			wantErr: true,
		},
		{
			name: "error expired token",
			fields: fields{
				repo: &authRepositoryMock{
					FindSessionByAccessTokenFunc: func(t string) (*entities.Session, error) {
						return &entities.Session{
							ID:        "1",
							ExpiresIn: time.Now().Add(-1 * time.Minute).Unix(),
						}, nil
					},
					DeleteSessionByIDFunc: func(_ string) error {
						return nil
					},
				},
				req: &dto.AuthRefreshTokenRequest{Token: "token"},
			},
			wantErr: true,
		},
		{
			name: "error finding user",
			fields: fields{
				repo: &authRepositoryMock{
					FindSessionByAccessTokenFunc: func(t string) (*entities.Session, error) {
						return &entities.Session{
							ID:        "1",
							ExpiresIn: time.Now().Add(10 * time.Minute).Unix(),
						}, nil
					},
					DeleteSessionByIDFunc: func(_ string) error {
						return nil
					},
				},
				usrv: &userServiceMock{
					UserByIDFunc: func(_ context.Context, _ uint) (*entities.User, error) {
						return nil, errTestError
					},
				},
				req: &dto.AuthRefreshTokenRequest{Token: "token"},
			},
			wantErr: true,
		},
		{
			name: "error creating session",
			fields: fields{
				repo: &authRepositoryMock{
					FindSessionByAccessTokenFunc: func(t string) (*entities.Session, error) {
						return &entities.Session{
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
					UserByIDFunc: func(_ context.Context, _ uint) (*entities.User, error) {
						return nil, nil
					},
				},
				req: &dto.AuthRefreshTokenRequest{Token: "token"},
			},
			wantErr: true,
		},
		{
			name: "ok",
			fields: fields{
				repo: &authRepositoryMock{
					FindSessionByAccessTokenFunc: func(t string) (*entities.Session, error) {
						return &entities.Session{
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
					UserByIDFunc: func(_ context.Context, _ uint) (*entities.User, error) {
						return &entities.User{
							ID:        1,
							FirstName: "Jane",
							LastName:  "Doe",
							Role:      1,
						}, nil
					},
				},
				req: &dto.AuthRefreshTokenRequest{Token: "token"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newService(defaultLogger, tt.fields.repo, tt.fields.usrv, "secret")
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
		req  *dto.AuthLoginRequest
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
					UserByPhoneAndPasswordFunc: func(_ context.Context, _ string, _ string) (*entities.User, error) {
						return nil, errTestError
					},
				},
				req: &dto.AuthLoginRequest{Phone: "1", Password: "1"},
			},
			wantErr: true,
		},
		{
			name: "error on create session",
			fields: fields{
				usrv: &userServiceMock{
					UserByPhoneAndPasswordFunc: func(_ context.Context, _ string, _ string) (*entities.User, error) {
						return &entities.User{ID: 1}, nil
					},
				},
				repo: &authRepositoryMock{
					CreateSessionFunc: func(_ uint) (string, error) {
						return "", errTestError
					},
				},
				req: &dto.AuthLoginRequest{Phone: "1", Password: "1"},
			},
			wantErr: true,
		},
		{
			name: "ok",
			fields: fields{
				usrv: &userServiceMock{
					UserByPhoneAndPasswordFunc: func(_ context.Context, _ string, _ string) (*entities.User, error) {
						return &entities.User{ID: 1, Role: 1, FirstName: "Jane", LastName: "Doe"}, nil
					},
				},
				repo: &authRepositoryMock{
					CreateSessionFunc: func(_ uint) (string, error) {
						return "token", nil
					},
				},
				req: &dto.AuthLoginRequest{Phone: "1", Password: "1"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newService(defaultLogger, tt.fields.repo, tt.fields.usrv, "secret")
			_, err := s.Login(context.Background(), tt.fields.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
