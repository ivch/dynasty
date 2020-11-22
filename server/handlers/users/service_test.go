package users

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/ivch/dynasty/common/logger"
)

var (
	defaultLogger *logger.StdLog
	errTestError  = errors.New("some err")
)

func TestMain(m *testing.M) {
	defaultLogger = logger.NewStdLog(logger.WithWriter(ioutil.Discard))
	os.Exit(m.Run())
}

func TestService_Register(t *testing.T) {
	type params struct {
		verifyRegCode bool
		maxMembers    int
		repo          userRepository
	}

	tests := []struct {
		name    string
		params  params
		input   *User
		wantErr bool
		want    *User
	}{
		{
			name: "error failed to check user",
			params: params{
				repo: &userRepositoryMock{
					GetUserByPhoneFunc: func(_ string) (*User, error) {
						return nil, errTestError
					},
				},
			},
			input:   &User{},
			wantErr: true,
			want:    nil,
		},
		{
			name: "error user exists",
			params: params{
				repo: &userRepositoryMock{
					GetUserByPhoneFunc: func(_ string) (*User, error) {
						return &User{
							ID: 1,
						}, nil
					},
				},
			},
			input:   &User{},
			wantErr: true,
			want:    nil,
		},
		{
			name: "error wrong reg code",
			params: params{
				verifyRegCode: true,
				repo: &userRepositoryMock{
					GetUserByPhoneFunc: func(_ string) (*User, error) {
						return nil, nil
					},
					ValidateRegCodeFunc: func(_ string) error {
						return errTestError
					},
				},
			},
			input:   &User{},
			wantErr: true,
			want:    nil,
		},
		{
			name: "error failed to find user by apartment",
			params: params{
				repo: &userRepositoryMock{
					GetUserByPhoneFunc: func(_ string) (*User, error) {
						return nil, nil
					},
					FindUserByApartmentFunc: func(_ uint, _ uint) (*User, error) {
						return nil, errTestError
					},
				},
			},
			input:   &User{},
			wantErr: true,
		},
		{
			name: "error master account already exists",
			params: params{
				repo: &userRepositoryMock{
					GetUserByPhoneFunc: func(_ string) (*User, error) {
						return nil, nil
					},
					FindUserByApartmentFunc: func(_ uint, _ uint) (*User, error) {
						return &User{}, nil
					},
				},
			},
			input:   &User{},
			wantErr: true,
		},
		{
			name: "error user not created",
			params: params{
				repo: &userRepositoryMock{
					GetUserByPhoneFunc: func(_ string) (*User, error) {
						return nil, nil
					},
					FindUserByApartmentFunc: func(_ uint, _ uint) (*User, error) {
						return nil, nil
					},
					CreateUserFunc: func(_ *User) error {
						return errTestError
					},
				},
			},
			input:   &User{},
			wantErr: true,
			want:    nil,
		},
		{
			name: "error failed to use reg code, user deleted",
			params: params{
				verifyRegCode: true,
				repo: &userRepositoryMock{
					GetUserByPhoneFunc: func(_ string) (*User, error) {
						return nil, nil
					},
					FindUserByApartmentFunc: func(_ uint, _ uint) (*User, error) {
						return nil, nil
					},
					CreateUserFunc: func(*User) error {
						return nil
					},
					ValidateRegCodeFunc: func(_ string) error {
						return nil
					},
					UseRegCodeFunc: func(_ string) error {
						return errTestError
					},
					DeleteUserFunc: func(_ *User) error {
						return nil
					},
				},
			},
			input:   &User{},
			wantErr: true,
			want:    nil,
		},
		{
			name: "error failed to use reg code, user not deleted",
			params: params{
				verifyRegCode: true,
				repo: &userRepositoryMock{
					GetUserByPhoneFunc: func(_ string) (*User, error) {
						return nil, nil
					},
					FindUserByApartmentFunc: func(_ uint, _ uint) (*User, error) {
						return nil, nil
					},
					CreateUserFunc: func(*User) error {
						return nil
					},
					ValidateRegCodeFunc: func(_ string) error {
						return nil
					},
					UseRegCodeFunc: func(_ string) error {
						return errTestError
					},
					DeleteUserFunc: func(_ *User) error {
						return errTestError
					},
				},
			},
			input:   &User{},
			wantErr: true,
			want:    nil,
		},
		{
			name: "ok",
			params: params{
				repo: &userRepositoryMock{
					GetUserByPhoneFunc: func(_ string) (*User, error) {
						return nil, nil
					},
					FindUserByApartmentFunc: func(_ uint, _ uint) (*User, error) {
						return nil, nil
					},
					CreateUserFunc: func(u *User) error {
						u.ID = 1
						return nil
					},
				},
			},
			input:   &User{Phone: "1", Password: "1"},
			wantErr: false,
			want: &User{
				ID:     1,
				Phone:  "1",
				Role:   defaultUserRole,
				Active: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(defaultLogger, tt.params.repo, tt.params.verifyRegCode, tt.params.maxMembers)
			got, err := s.Register(context.Background(), tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil {
				got.Password = ""
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Register() got = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestService_UserByPhoneAndPassword(t *testing.T) {
	type params struct {
		verifyRegCode bool
		maxMembers    int
		repo          userRepository
	}

	type input struct {
		phone    string
		password string
	}

	tests := []struct {
		name    string
		params  params
		input   input
		wantErr bool
		want    *User
	}{
		{
			name: "error no user",
			params: params{
				repo: &userRepositoryMock{
					GetUserByPhoneFunc: func(_ string) (*User, error) {
						return nil, errTestError
					},
				},
			},
			input: input{
				phone:    "1",
				password: "2",
			},
			wantErr: true,
			want:    nil,
		},
		{
			name: "error wrong password",
			params: params{
				repo: &userRepositoryMock{
					GetUserByPhoneFunc: func(_ string) (*User, error) {
						p, _ := hashAndSalt("1")
						return &User{Password: p}, nil
					},
				},
			},
			input: input{
				phone:    "1",
				password: "2",
			},
			wantErr: true,
			want:    nil,
		},
		{
			name: "ok",
			params: params{
				repo: &userRepositoryMock{
					GetUserByPhoneFunc: func(_ string) (*User, error) {
						p, _ := hashAndSalt("1")
						return &User{ID: 1, FirstName: "a", LastName: "b", Role: 1, Password: p}, nil
					},
				},
			},
			input: input{
				phone:    "1",
				password: "1",
			},
			wantErr: false,
			want: &User{
				ID:        1,
				FirstName: "a",
				LastName:  "b",
				Role:      1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(defaultLogger, tt.params.repo, tt.params.verifyRegCode, tt.params.maxMembers)
			got, err := s.UserByPhoneAndPassword(context.Background(), tt.input.phone, tt.input.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserByPhoneAndPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != nil {
				got.Password = ""
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserByPhoneAndPassword() got = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestService_UserByID(t *testing.T) {
	type params struct {
		verifyRegCode bool
		maxMembers    int
		repo          userRepository
	}

	tests := []struct {
		name    string
		params  params
		input   uint
		wantErr bool
		want    *User
	}{
		{
			name: "error no user",
			params: params{
				repo: &userRepositoryMock{
					GetUserByIDFunc: func(id uint) (*User, error) {
						return nil, errTestError
					},
				},
			},
			input:   1,
			wantErr: true,
			want:    nil,
		},
		{
			name: "ok",
			params: params{
				repo: &userRepositoryMock{
					GetUserByIDFunc: func(id uint) (*User, error) {
						return &User{
							ID: 1,
							Building: Building{
								ID:      1,
								Name:    "a",
								Address: "b",
							},
							Entry: Entry{
								ID:         1,
								Name:       "1",
								BuildingID: 1,
							},
							Apartment:  1,
							Email:      "a",
							Password:   "b",
							Phone:      "c",
							FirstName:  "d",
							LastName:   "e",
							BuildingID: 1,
						}, nil
					},
				},
			},
			wantErr: false,
			want: &User{
				ID: 1,
				Building: Building{
					ID:      1,
					Name:    "a",
					Address: "b",
				},
				Entry: Entry{
					ID:         1,
					Name:       "1",
					BuildingID: 1,
				},
				Apartment:  1,
				Email:      "a",
				Password:   "b",
				Phone:      "c",
				FirstName:  "d",
				LastName:   "e",
				BuildingID: 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(defaultLogger, tt.params.repo, tt.params.verifyRegCode, tt.params.maxMembers)
			got, err := s.UserByID(context.Background(), tt.input)
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
