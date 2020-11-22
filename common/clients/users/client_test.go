package users

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/ivch/dynasty/server/handlers/users"
	"github.com/ivch/dynasty/server/handlers/users/transport"
)

var errTestError = errors.New("some err")

func TestClient_UserByID(t *testing.T) {
	tests := []struct {
		name    string
		usrv    userService
		id      uint
		wantErr bool
		want    *transport.UserByIDResponse
	}{
		{
			name:    "error wrong id",
			id:      0,
			wantErr: true,
		},
		{
			name: "error from service",
			id:   1,
			usrv: &userServiceMock{
				UserByIDFunc: func(_ context.Context, _ uint) (*users.User, error) {
					return nil, errTestError
				},
			},
			wantErr: true,
		},
		{
			name: "ok",
			id:   1,
			usrv: &userServiceMock{
				UserByIDFunc: func(_ context.Context, id uint) (*users.User, error) {
					if id != 1 {
						return nil, errTestError
					}
					return &users.User{
						ID:        1,
						FirstName: "1",
						LastName:  "1",
						Phone:     "1",
						Email:     "1",
						Role:      1,
					}, nil
				},
			},
			wantErr: false,
			want: &transport.UserByIDResponse{
				ID:        1,
				Email:     "1",
				Phone:     "1",
				FirstName: "1",
				LastName:  "1",
				Role:      1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(tt.usrv)
			got, err := s.UserByID(context.Background(), tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserByID() got = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestClient_UserByPhoneAndPassword(t *testing.T) {
	tests := []struct {
		name     string
		usrv     userService
		phone    string
		password string
		wantErr  bool
		want     *transport.UserByIDResponse
	}{
		{
			name:    "error empty phone",
			phone:   "",
			wantErr: true,
		},
		{
			name:  "error from service",
			phone: "1",
			usrv: &userServiceMock{
				UserByPhoneAndPasswordFunc: func(_ context.Context, _ string, _ string) (*users.User, error) {
					return nil, errTestError
				},
			},
			wantErr: true,
		},
		{
			name:     "ok",
			phone:    "1",
			password: "1",
			usrv: &userServiceMock{
				UserByPhoneAndPasswordFunc: func(_ context.Context, phone string, password string) (*users.User, error) {
					if phone != "1" || password != "1" {
						return nil, errTestError
					}
					return &users.User{
						ID:        1,
						FirstName: "1",
						LastName:  "1",
						Role:      1,
					}, nil
				},
			},
			wantErr: false,
			want: &transport.UserByIDResponse{
				ID:        1,
				FirstName: "1",
				LastName:  "1",
				Role:      1,
				Building:  &users.Building{},
				Entry:     &users.Entry{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(tt.usrv)
			got, err := s.UserByPhoneAndPassword(context.Background(), tt.phone, tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserByPhoneAndPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserByPhoneAndPassword() got = %#v, want %#v", got, tt.want)
			}
		})
	}
}
