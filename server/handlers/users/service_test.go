package users_test

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/ivch/dynasty/common/logger"
	"github.com/ivch/dynasty/server/handlers/users"
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
		repo          users.UserRepository
	}

	tests := []struct {
		name    string
		params  params
		input   *users.User
		wantErr bool
		want    *users.User
	}{
		{
			name: "error getUserByEmail",
			params: params{
				repo: &users.UserRepositoryMock{
					GetUserByEmailFunc: func(_ string) (*users.User, error) {
						return nil, errTestError
					},
				},
			},
			input:   &users.User{},
			wantErr: true,
			want:    nil,
		},
		{
			name: "error email in use",
			params: params{
				repo: &users.UserRepositoryMock{
					GetUserByEmailFunc: func(_ string) (*users.User, error) {
						return &users.User{}, nil
					},
				},
			},
			input:   &users.User{},
			wantErr: true,
			want:    nil,
		},
		{
			name: "error failed to check user",
			params: params{
				repo: &users.UserRepositoryMock{
					GetUserByEmailFunc: func(_ string) (*users.User, error) {
						return nil, nil
					},
					GetUserByPhoneFunc: func(_ string) (*users.User, error) {
						return nil, errTestError
					},
				},
			},
			input:   &users.User{},
			wantErr: true,
			want:    nil,
		},
		{
			name: "error user exists",
			params: params{
				repo: &users.UserRepositoryMock{
					GetUserByEmailFunc: func(_ string) (*users.User, error) {
						return nil, nil
					},
					GetUserByPhoneFunc: func(_ string) (*users.User, error) {
						return &users.User{
							ID: 1,
						}, nil
					},
				},
			},
			input:   &users.User{},
			wantErr: true,
			want:    nil,
		},
		{
			name: "error failed to find user by apartment",
			params: params{
				repo: &users.UserRepositoryMock{
					GetUserByEmailFunc: func(_ string) (*users.User, error) {
						return nil, nil
					},
					GetUserByPhoneFunc: func(_ string) (*users.User, error) {
						return nil, nil
					},
					FindUserByApartmentFunc: func(_ uint, _ uint) (*users.User, error) {
						return nil, errTestError
					},
				},
			},
			input:   &users.User{},
			wantErr: true,
		},
		{
			name: "error master account already exists",
			params: params{
				repo: &users.UserRepositoryMock{
					GetUserByEmailFunc: func(_ string) (*users.User, error) {
						return nil, nil
					},
					GetUserByPhoneFunc: func(_ string) (*users.User, error) {
						return nil, nil
					},
					FindUserByApartmentFunc: func(_ uint, _ uint) (*users.User, error) {
						return &users.User{}, nil
					},
				},
			},
			input:   &users.User{},
			wantErr: true,
		},
		{
			name: "error wrong reg code",
			params: params{
				verifyRegCode: true,
				repo: &users.UserRepositoryMock{
					GetUserByEmailFunc: func(_ string) (*users.User, error) {
						return nil, nil
					},
					GetUserByPhoneFunc: func(_ string) (*users.User, error) {
						return nil, nil
					},
					FindUserByApartmentFunc: func(_ uint, _ uint) (*users.User, error) {
						return &users.User{}, nil
					},
					ValidateRegCodeFunc: func(_ string) error {
						return errTestError
					},
				},
			},
			input:   &users.User{},
			wantErr: true,
			want:    nil,
		},
		{
			name: "error user not created",
			params: params{
				repo: &users.UserRepositoryMock{
					GetUserByEmailFunc: func(_ string) (*users.User, error) {
						return nil, nil
					},
					GetUserByPhoneFunc: func(_ string) (*users.User, error) {
						return nil, nil
					},
					FindUserByApartmentFunc: func(_ uint, _ uint) (*users.User, error) {
						return nil, nil
					},
					CreateUserFunc: func(_ *users.User) error {
						return errTestError
					},
				},
			},
			input:   &users.User{},
			wantErr: true,
			want:    nil,
		},
		{
			name: "error failed to use reg code, user deleted",
			params: params{
				verifyRegCode: true,
				repo: &users.UserRepositoryMock{
					GetUserByEmailFunc: func(_ string) (*users.User, error) {
						return nil, nil
					},
					GetUserByPhoneFunc: func(_ string) (*users.User, error) {
						return nil, nil
					},
					FindUserByApartmentFunc: func(_ uint, _ uint) (*users.User, error) {
						return nil, nil
					},
					CreateUserFunc: func(*users.User) error {
						return nil
					},
					ValidateRegCodeFunc: func(_ string) error {
						return nil
					},
					UseRegCodeFunc: func(_ string) error {
						return errTestError
					},
					DeleteUserFunc: func(_ *users.User) error {
						return nil
					},
				},
			},
			input:   &users.User{},
			wantErr: true,
			want:    nil,
		},
		{
			name: "error failed to use reg code, user not deleted",
			params: params{
				verifyRegCode: true,
				repo: &users.UserRepositoryMock{
					GetUserByEmailFunc: func(_ string) (*users.User, error) {
						return nil, nil
					},
					GetUserByPhoneFunc: func(_ string) (*users.User, error) {
						return nil, nil
					},
					FindUserByApartmentFunc: func(_ uint, _ uint) (*users.User, error) {
						return nil, nil
					},
					CreateUserFunc: func(*users.User) error {
						return nil
					},
					ValidateRegCodeFunc: func(_ string) error {
						return nil
					},
					UseRegCodeFunc: func(_ string) error {
						return errTestError
					},
					DeleteUserFunc: func(_ *users.User) error {
						return errTestError
					},
				},
			},
			input:   &users.User{},
			wantErr: true,
			want:    nil,
		},
		{
			name: "ok",
			params: params{
				repo: &users.UserRepositoryMock{
					GetUserByEmailFunc: func(_ string) (*users.User, error) {
						return nil, nil
					},
					GetUserByPhoneFunc: func(_ string) (*users.User, error) {
						return nil, nil
					},
					FindUserByApartmentFunc: func(_ uint, _ uint) (*users.User, error) {
						return nil, nil
					},
					CreateUserFunc: func(u *users.User) error {
						u.ID = 1
						return nil
					},
				},
			},
			input:   &users.User{Phone: "1", Password: "1"},
			wantErr: false,
			want: &users.User{
				ID:     1,
				Phone:  "1",
				Role:   users.DefaultUserRole,
				Active: true,
			},
		},
		{
			name: "error predefined user wrong code",
			params: params{
				repo: &users.UserRepositoryMock{
					GetUserByEmailFunc: func(_ string) (*users.User, error) {
						return nil, nil
					},
					GetUserByPhoneFunc: func(_ string) (*users.User, error) {
						return nil, nil
					},
					FindUserByApartmentFunc: func(_ uint, _ uint) (*users.User, error) {
						return &users.User{Role: users.PredefinedUserRole, RegCode: "abc"}, nil
					},
				},
			},
			input:   &users.User{Phone: "1", Password: "1", RegCode: "cba"},
			wantErr: true,
			want:    nil,
		},
		{
			name: "error predefined user update error",
			params: params{
				repo: &users.UserRepositoryMock{
					GetUserByEmailFunc: func(_ string) (*users.User, error) {
						return nil, nil
					},
					GetUserByPhoneFunc: func(_ string) (*users.User, error) {
						return nil, nil
					},
					FindUserByApartmentFunc: func(_ uint, _ uint) (*users.User, error) {
						return &users.User{Role: users.PredefinedUserRole, RegCode: "abc"}, nil
					},
					UpdateUserFunc: func(_ *users.UserUpdate) error {
						return errTestError
					},
				},
			},
			input:   &users.User{Phone: "1", Password: "1", RegCode: "abc"},
			wantErr: true,
			want:    nil,
		},
		{
			name: "ok predefined user",
			params: params{
				repo: &users.UserRepositoryMock{
					GetUserByEmailFunc: func(_ string) (*users.User, error) {
						return nil, nil
					},
					GetUserByPhoneFunc: func(_ string) (*users.User, error) {
						return nil, nil
					},
					FindUserByApartmentFunc: func(_ uint, _ uint) (*users.User, error) {
						return &users.User{ID: 1, Role: users.PredefinedUserRole, RegCode: "abc"}, nil
					},
					UpdateUserFunc: func(_ *users.UserUpdate) error {
						return nil
					},
				},
			},
			input:   &users.User{Phone: "1", Password: "1", RegCode: "abc"},
			wantErr: false,
			want: &users.User{
				ID:     1,
				Phone:  "1",
				Role:   users.DefaultUserRole,
				Active: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := users.New(defaultLogger, tt.params.repo, tt.params.verifyRegCode, tt.params.maxMembers, nil)
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
		repo          users.UserRepository
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
		want    *users.User
	}{
		{
			name: "error no user",
			params: params{
				repo: &users.UserRepositoryMock{
					GetUserByPhoneFunc: func(_ string) (*users.User, error) {
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
				repo: &users.UserRepositoryMock{
					GetUserByPhoneFunc: func(_ string) (*users.User, error) {
						p, _ := users.HashAndSalt("1")
						return &users.User{Password: p}, nil
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
				repo: &users.UserRepositoryMock{
					GetUserByPhoneFunc: func(_ string) (*users.User, error) {
						p, _ := users.HashAndSalt("1")
						return &users.User{ID: 1, FirstName: "a", LastName: "b", Role: 1, Password: p}, nil
					},
				},
			},
			input: input{
				phone:    "1",
				password: "1",
			},
			wantErr: false,
			want: &users.User{
				ID:        1,
				FirstName: "a",
				LastName:  "b",
				Role:      1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := users.New(defaultLogger, tt.params.repo, tt.params.verifyRegCode, tt.params.maxMembers, nil)
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
		repo          users.UserRepository
	}

	tests := []struct {
		name    string
		params  params
		input   uint
		wantErr bool
		want    *users.User
	}{
		{
			name: "error no user",
			params: params{
				repo: &users.UserRepositoryMock{
					GetUserByIDFunc: func(id uint) (*users.User, error) {
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
				repo: &users.UserRepositoryMock{
					GetUserByIDFunc: func(id uint) (*users.User, error) {
						return &users.User{
							ID: 1,
							Building: users.Building{
								ID:      1,
								Name:    "a",
								Address: "b",
							},
							Entry: users.Entry{
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
			want: &users.User{
				ID: 1,
				Building: users.Building{
					ID:      1,
					Name:    "a",
					Address: "b",
				},
				Entry: users.Entry{
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
			s := users.New(defaultLogger, tt.params.repo, tt.params.verifyRegCode, tt.params.maxMembers, nil)
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

func Test_ServiceUpdate(t *testing.T) {
	type params struct {
		verifyRegCode bool
		maxMembers    int
		repo          users.UserRepository
	}

	tests := []struct {
		name    string
		params  params
		input   *users.UserUpdate
		wantErr bool
	}{
		{
			name: "error finding user",
			params: params{
				repo: &users.UserRepositoryMock{
					GetUserByIDFunc: func(_ uint) (*users.User, error) {
						return nil, errTestError
					},
				},
			},
			input: &users.UserUpdate{
				ID: 1,
			},
			wantErr: true,
		},
		{
			name: "error if empty password",
			params: params{
				repo: &users.UserRepositoryMock{
					GetUserByIDFunc: func(_ uint) (*users.User, error) {
						return nil, nil
					},
				},
			},
			input: &users.UserUpdate{
				ID:          1,
				Password:    nil,
				NewPassword: func(s string) *string { return &s }("2"),
			},
			wantErr: true,
		},
		{
			name: "error old password mismatch",
			params: params{
				repo: &users.UserRepositoryMock{
					GetUserByIDFunc: func(_ uint) (*users.User, error) {
						return &users.User{
							Password: "1",
						}, nil
					},
				},
			},
			input: &users.UserUpdate{
				ID:          1,
				Password:    func(s string) *string { return &s }("2"),
				NewPassword: func(s string) *string { return &s }("2"),
			},
			wantErr: true,
		},
		{
			name: "error if pwd not changed",
			params: params{
				repo: &users.UserRepositoryMock{
					GetUserByIDFunc: func(_ uint) (*users.User, error) {
						return &users.User{
							Password: "1",
						}, nil
					},
					UpdateUserFunc: func(u *users.UserUpdate) error {
						if u.NewPassword != nil {
							if u.Password == u.NewPassword {
								return errTestError
							}

							if *u.NewPassword == "2" {
								return errTestError
							}
						}

						return nil
					},
				},
			},
			input: &users.UserUpdate{
				ID:          1,
				Password:    func(s string) *string { return &s }("1"),
				NewPassword: func(s string) *string { return &s }("2"),
			},
			wantErr: true,
		},
		{
			name: "ok w/o password change",
			params: params{
				repo: &users.UserRepositoryMock{
					GetUserByIDFunc: func(_ uint) (*users.User, error) {
						return nil, nil
					},
					UpdateUserFunc: func(u *users.UserUpdate) error {
						return nil
					},
				},
			},
			input: &users.UserUpdate{
				ID: 1,
			},
			wantErr: false,
		},
		{
			name: "ok w/ password change",
			params: params{
				repo: &users.UserRepositoryMock{
					GetUserByIDFunc: func(_ uint) (*users.User, error) {
						testPass, _ := users.HashAndSalt("1")
						return &users.User{
							Password: testPass,
						}, nil
					},
					UpdateUserFunc: func(u *users.UserUpdate) error {
						if u.NewPassword != nil {
							if u.Password == u.NewPassword {
								return errTestError
							}

							if *u.NewPassword == "2" {
								return errTestError
							}
						}

						return nil
					},
				},
			},
			input: &users.UserUpdate{
				ID:          1,
				Password:    func(s string) *string { return &s }("1"),
				NewPassword: func(s string) *string { return &s }("3"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := users.New(defaultLogger, tt.params.repo, tt.params.verifyRegCode, tt.params.maxMembers, nil)
			err := s.Update(context.Background(), tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_ServiceRecovery(t *testing.T) {
	type params struct {
		verifyRegCode bool
		maxMembers    int
		repo          users.UserRepository
		email         users.MailSender
	}

	tests := []struct {
		name    string
		params  params
		input   *users.User
		wantErr bool
	}{
		{
			name: "error finding user",
			params: params{
				repo: &users.UserRepositoryMock{
					GetUserByPhoneFunc: func(_ string) (*users.User, error) {
						return nil, errTestError
					},
				},
			},
			input: &users.User{
				Phone: "1",
			},
			wantErr: true,
		},
		{
			name: "error empty user",
			params: params{
				repo: &users.UserRepositoryMock{
					GetUserByPhoneFunc: func(_ string) (*users.User, error) {
						return nil, nil
					},
				},
			},
			input: &users.User{
				Phone: "1",
			},
			wantErr: true,
		},
		{
			name: "error getting count",
			params: params{
				repo: &users.UserRepositoryMock{
					GetUserByPhoneFunc: func(_ string) (*users.User, error) {
						return &users.User{ID: 1}, nil
					},
					CountRecoveryCodesByUserIn24hFunc: func(_ uint) (int, error) {
						return 0, errTestError
					},
				},
			},
			input: &users.User{
				Phone: "1",
			},
			wantErr: true,
		},
		{
			name: "error limit exceeded",
			params: params{
				repo: &users.UserRepositoryMock{
					GetUserByPhoneFunc: func(_ string) (*users.User, error) {
						return &users.User{ID: 1}, nil
					},
					CountRecoveryCodesByUserIn24hFunc: func(_ uint) (int, error) {
						return 4, nil
					},
				},
			},
			input: &users.User{
				Phone: "1",
			},
			wantErr: true,
		},
		{
			name: "error wrong email",
			params: params{
				repo: &users.UserRepositoryMock{
					GetUserByPhoneFunc: func(_ string) (*users.User, error) {
						return &users.User{
							ID:    1,
							Email: "a",
						}, nil
					},
					CountRecoveryCodesByUserIn24hFunc: func(_ uint) (int, error) {
						return 2, nil
					},
				},
			},
			input: &users.User{
				Phone: "1",
				Email: "b",
			},
			wantErr: true,
		},
		{
			name: "error create code",
			params: params{
				repo: &users.UserRepositoryMock{
					GetUserByPhoneFunc: func(_ string) (*users.User, error) {
						return &users.User{
							ID:    1,
							Email: "a",
						}, nil
					},
					CountRecoveryCodesByUserIn24hFunc: func(_ uint) (int, error) {
						return 2, nil
					},
					CreateRecoverCodeFunc: func(_ *users.PasswordRecovery) error {
						return errTestError
					},
				},
			},
			input: &users.User{
				Phone: "1",
				Email: "a",
			},
			wantErr: true,
		},
		{
			name: "error send email",
			params: params{
				repo: &users.UserRepositoryMock{
					GetUserByPhoneFunc: func(_ string) (*users.User, error) {
						return &users.User{
							ID:    1,
							Email: "a",
						}, nil
					},
					CountRecoveryCodesByUserIn24hFunc: func(_ uint) (int, error) {
						return 2, nil
					},
					CreateRecoverCodeFunc: func(_ *users.PasswordRecovery) error {
						return nil
					},
				},
				email: &users.MailSenderMock{
					SendRecoveryCodeEmailFunc: func(_, _, _ string) error {
						return errTestError
					},
				},
			},
			input: &users.User{
				Phone: "1",
				Email: "a",
			},
			wantErr: true,
		},
		{
			name: "ok",
			params: params{
				repo: &users.UserRepositoryMock{
					GetUserByPhoneFunc: func(_ string) (*users.User, error) {
						return &users.User{
							ID:    1,
							Email: "a",
						}, nil
					},
					CountRecoveryCodesByUserIn24hFunc: func(_ uint) (int, error) {
						return 2, nil
					},
					CreateRecoverCodeFunc: func(_ *users.PasswordRecovery) error {
						return nil
					},
				},
				email: &users.MailSenderMock{
					SendRecoveryCodeEmailFunc: func(_, _, _ string) error {
						return nil
					},
				},
			},
			input: &users.User{
				Phone: "1",
				Email: "a",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := users.New(defaultLogger, tt.params.repo, tt.params.verifyRegCode, tt.params.maxMembers, tt.params.email)
			err := s.RecoveryCode(context.Background(), tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_ServiceResetPassword(t *testing.T) {
	type params struct {
		verifyRegCode bool
		maxMembers    int
		repo          users.UserRepository
	}

	type input struct {
		code string
		u    *users.UserUpdate
	}
	tests := []struct {
		name    string
		params  params
		input   input
		wantErr bool
	}{
		{
			name: "error finding code",
			params: params{
				repo: &users.UserRepositoryMock{
					GetRecoveryCodeFunc: func(_ *users.PasswordRecovery) (*users.PasswordRecovery, error) {
						return nil, errTestError
					},
				},
			},
			input: input{
				code: "1",
			},
			wantErr: true,
		},
		{
			name: "error code outdated",
			params: params{
				repo: &users.UserRepositoryMock{
					GetRecoveryCodeFunc: func(_ *users.PasswordRecovery) (*users.PasswordRecovery, error) {
						return &users.PasswordRecovery{
							CreatedAt: func(t time.Time) *time.Time { return &t }(time.Now().Add(-5 * time.Hour)),
						}, nil
					},
				},
			},
			input: input{
				code: "1",
			},
			wantErr: true,
		},
		{
			name: "error no user",
			params: params{
				repo: &users.UserRepositoryMock{
					GetRecoveryCodeFunc: func(_ *users.PasswordRecovery) (*users.PasswordRecovery, error) {
						return &users.PasswordRecovery{
							UserID:    1,
							CreatedAt: func(t time.Time) *time.Time { return &t }(time.Now()),
						}, nil
					},
					GetUserByIDFunc: func(_ uint) (*users.User, error) {
						return nil, errTestError
					},
				},
			},
			input: input{
				code: "1",
				u:    &users.UserUpdate{},
			},
			wantErr: true,
		},
		{
			name: "error on reset password",
			params: params{
				repo: &users.UserRepositoryMock{
					GetRecoveryCodeFunc: func(_ *users.PasswordRecovery) (*users.PasswordRecovery, error) {
						return &users.PasswordRecovery{
							UserID:    1,
							CreatedAt: func(t time.Time) *time.Time { return &t }(time.Now()),
						}, nil
					},
					GetUserByIDFunc: func(_ uint) (*users.User, error) {
						return nil, nil
					},
					ResetPasswordFunc: func(_ uint, _ *users.UserUpdate) error {
						return errTestError
					},
				},
			},
			input: input{
				code: "1",
				u: &users.UserUpdate{
					NewPassword: func(s string) *string { return &s }("1"),
				},
			},
			wantErr: true,
		},
		{
			name: "ok",
			params: params{
				repo: &users.UserRepositoryMock{
					GetRecoveryCodeFunc: func(_ *users.PasswordRecovery) (*users.PasswordRecovery, error) {
						return &users.PasswordRecovery{
							UserID:    1,
							CreatedAt: func(t time.Time) *time.Time { return &t }(time.Now()),
						}, nil
					},
					GetUserByIDFunc: func(_ uint) (*users.User, error) {
						return nil, nil
					},
					ResetPasswordFunc: func(_ uint, _ *users.UserUpdate) error {
						return nil
					},
				},
			},
			input: input{
				code: "1",
				u: &users.UserUpdate{
					NewPassword: func(s string) *string { return &s }("1"),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := users.New(defaultLogger, tt.params.repo, tt.params.verifyRegCode, tt.params.maxMembers, nil)
			err := s.ResetPassword(context.Background(), tt.input.code, tt.input.u)
			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
