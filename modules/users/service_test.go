package users

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/rs/zerolog"

	"github.com/ivch/dynasty/models"
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

func TestService_Register(t *testing.T) {
	type fields struct {
		verifyRegCode bool
		repo          userRepository
	}

	type args struct {
		r *userRegisterRequest
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		want    *userRegisterResponse
	}{
		{
			name: "error failed to check user",
			fields: fields{
				repo: &userRepositoryMock{
					GetUserByPhoneFunc: func(_ string) (*models.User, error) {
						return nil, errTestError
					},
				},
			},
			args:    args{r: &userRegisterRequest{}},
			wantErr: true,
			want:    nil,
		},
		{
			name: "error user exists",
			fields: fields{
				repo: &userRepositoryMock{
					GetUserByPhoneFunc: func(_ string) (*models.User, error) {
						return &models.User{
							ID: 1,
						}, nil
					},
				},
			},
			args:    args{r: &userRegisterRequest{}},
			wantErr: true,
			want:    nil,
		},
		{
			name: "error wrong reg code",
			fields: fields{
				verifyRegCode: true,
				repo: &userRepositoryMock{
					GetUserByPhoneFunc: func(_ string) (*models.User, error) {
						return nil, nil
					},
					ValidateRegCodeFunc: func(_ string) error {
						return errTestError
					},
				},
			},
			args:    args{r: &userRegisterRequest{}},
			wantErr: true,
			want:    nil,
		},
		{
			name: "error user not created",
			fields: fields{
				repo: &userRepositoryMock{
					GetUserByPhoneFunc: func(_ string) (*models.User, error) {
						return nil, nil
					},
					CreateUserFunc: func(_ *models.User) error {
						return errTestError
					},
				},
			},
			args:    args{r: &userRegisterRequest{}},
			wantErr: true,
			want:    nil,
		},
		{
			name: "error failed to use reg code, user deleted",
			fields: fields{
				verifyRegCode: true,
				repo: &userRepositoryMock{
					GetUserByPhoneFunc: func(_ string) (*models.User, error) {
						return nil, nil
					},
					CreateUserFunc: func(*models.User) error {
						return nil
					},
					ValidateRegCodeFunc: func(_ string) error {
						return nil
					},
					UseRegCodeFunc: func(_ string) error {
						return errTestError
					},
					DeleteUserFunc: func(_ *models.User) error {
						return nil
					},
				},
			},
			args:    args{r: &userRegisterRequest{}},
			wantErr: true,
			want:    nil,
		},
		{
			name: "error failed to use reg code, user not deleted",
			fields: fields{
				verifyRegCode: true,
				repo: &userRepositoryMock{
					GetUserByPhoneFunc: func(_ string) (*models.User, error) {
						return nil, nil
					},
					CreateUserFunc: func(*models.User) error {
						return nil
					},
					ValidateRegCodeFunc: func(_ string) error {
						return nil
					},
					UseRegCodeFunc: func(_ string) error {
						return errTestError
					},
					DeleteUserFunc: func(_ *models.User) error {
						return errTestError
					},
				},
			},
			args:    args{r: &userRegisterRequest{}},
			wantErr: true,
			want:    nil,
		},
		{
			name: "ok",
			fields: fields{
				repo: &userRepositoryMock{
					GetUserByPhoneFunc: func(_ string) (*models.User, error) {
						return nil, nil
					},
					CreateUserFunc: func(u *models.User) error {
						u.ID = 1
						return nil
					},
				},
			},
			args:    args{r: &userRegisterRequest{Phone: "1"}},
			wantErr: false,
			want: &userRegisterResponse{
				ID:    1,
				Phone: "1",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newService(defaultLogger, tt.fields.repo, tt.fields.verifyRegCode)
			got, err := s.Register(context.Background(), tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Register() got = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestService_UserByPhoneAndPassword(t *testing.T) {
	type fields struct {
		verifyRegCode bool
		repo          userRepository
	}

	type args struct {
		phone    string
		password string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		want    *userAuthResponse
	}{
		{
			name: "error no user",
			fields: fields{
				repo: &userRepositoryMock{
					GetUserByPhoneFunc: func(_ string) (*models.User, error) {
						return nil, errTestError
					},
				},
			},
			args: args{
				phone:    "1",
				password: "2",
			},
			wantErr: true,
			want:    nil,
		},
		{
			name: "error wrong password",
			fields: fields{
				repo: &userRepositoryMock{
					GetUserByPhoneFunc: func(_ string) (*models.User, error) {
						p, _ := hashAndSalt("1")
						return &models.User{Password: p}, nil
					},
				},
			},
			args: args{
				phone:    "1",
				password: "2",
			},
			wantErr: true,
			want:    nil,
		},
		{
			name: "ok",
			fields: fields{
				repo: &userRepositoryMock{
					GetUserByPhoneFunc: func(_ string) (*models.User, error) {
						p, _ := hashAndSalt("1")
						return &models.User{ID: 1, FirstName: "a", LastName: "b", Role: 1, Password: p}, nil
					},
				},
			},
			args: args{
				phone:    "1",
				password: "1",
			},
			wantErr: false,
			want: &userAuthResponse{
				ID:        1,
				FirstName: "a",
				LastName:  "b",
				Role:      1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newService(defaultLogger, tt.fields.repo, tt.fields.verifyRegCode)
			got, err := s.UserByPhoneAndPassword(context.Background(), tt.args.phone, tt.args.password)
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

func TestService_UserByID(t *testing.T) {
	type fields struct {
		verifyRegCode bool
		repo          userRepository
	}

	type args struct {
		id uint
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		want    *userByIDResponse
	}{
		{
			name: "error no user",
			fields: fields{
				repo: &userRepositoryMock{
					GetUserByIDFunc: func(id uint) (*models.User, error) {
						return nil, errTestError
					},
				},
			},
			args:    args{id: 1},
			wantErr: true,
			want:    nil,
		},
		{
			name: "ok",
			fields: fields{
				repo: &userRepositoryMock{
					GetUserByIDFunc: func(id uint) (*models.User, error) {
						return &models.User{
							ID: 1,
							Building: models.Building{
								ID:      1,
								Name:    "a",
								Address: "b",
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
			want: &userByIDResponse{
				ID:        1,
				Apartment: 1,
				FirstName: "d",
				LastName:  "e",
				Phone:     "c",
				Email:     "a",
				Building: models.Building{
					ID:      1,
					Name:    "a",
					Address: "b",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newService(defaultLogger, tt.fields.repo, tt.fields.verifyRegCode)
			got, err := s.UserByID(context.Background(), tt.args.id)
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
