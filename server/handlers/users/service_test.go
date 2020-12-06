package users

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
	"time"

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
			s := New(defaultLogger, tt.params.repo, tt.params.verifyRegCode, tt.params.maxMembers, nil)
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
			s := New(defaultLogger, tt.params.repo, tt.params.verifyRegCode, tt.params.maxMembers, nil)
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
			s := New(defaultLogger, tt.params.repo, tt.params.verifyRegCode, tt.params.maxMembers, nil)
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
		repo          userRepository
	}

	tests := []struct {
		name    string
		params  params
		input   *UserUpdate
		wantErr bool
	}{
		{
			name: "error finding user",
			params: params{
				repo: &userRepositoryMock{
					GetUserByIDFunc: func(_ uint) (*User, error) {
						return nil, errTestError
					},
				},
			},
			input: &UserUpdate{
				ID: 1,
			},
			wantErr: true,
		},
		{
			name: "error if empty password",
			params: params{
				repo: &userRepositoryMock{
					GetUserByIDFunc: func(_ uint) (*User, error) {
						return nil, nil
					},
				},
			},
			input: &UserUpdate{
				ID:          1,
				Password:    nil,
				NewPassword: func(s string) *string { return &s }("2"),
			},
			wantErr: true,
		},
		{
			name: "error old password mismatch",
			params: params{
				repo: &userRepositoryMock{
					GetUserByIDFunc: func(_ uint) (*User, error) {
						return &User{
							Password: "1",
						}, nil
					},
				},
			},
			input: &UserUpdate{
				ID:          1,
				Password:    func(s string) *string { return &s }("2"),
				NewPassword: func(s string) *string { return &s }("2"),
			},
			wantErr: true,
		},
		{
			name: "error if pwd not changed",
			params: params{
				repo: &userRepositoryMock{
					GetUserByIDFunc: func(_ uint) (*User, error) {
						return &User{
							Password: "1",
						}, nil
					},
					UpdateUserFunc: func(u *UserUpdate) error {
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
			input: &UserUpdate{
				ID:          1,
				Password:    func(s string) *string { return &s }("1"),
				NewPassword: func(s string) *string { return &s }("2"),
			},
			wantErr: true,
		},
		{
			name: "ok w/o password change",
			params: params{
				repo: &userRepositoryMock{
					GetUserByIDFunc: func(_ uint) (*User, error) {
						return nil, nil
					},
					UpdateUserFunc: func(u *UserUpdate) error {
						return nil
					},
				},
			},
			input: &UserUpdate{
				ID: 1,
			},
			wantErr: false,
		},
		{
			name: "ok w/ password change",
			params: params{
				repo: &userRepositoryMock{
					GetUserByIDFunc: func(_ uint) (*User, error) {
						testPass, _ := hashAndSalt("1")
						return &User{
							Password: testPass,
						}, nil
					},
					UpdateUserFunc: func(u *UserUpdate) error {
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
			input: &UserUpdate{
				ID:          1,
				Password:    func(s string) *string { return &s }("1"),
				NewPassword: func(s string) *string { return &s }("3"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(defaultLogger, tt.params.repo, tt.params.verifyRegCode, tt.params.maxMembers, nil)
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
		repo          userRepository
		email         mailSender
	}

	tests := []struct {
		name    string
		params  params
		input   *User
		wantErr bool
	}{
		{
			name: "error finding user",
			params: params{
				repo: &userRepositoryMock{
					GetUserByPhoneFunc: func(_ string) (*User, error) {
						return nil, errTestError
					},
				},
			},
			input: &User{
				Phone: "1",
			},
			wantErr: true,
		},
		{
			name: "error getting count",
			params: params{
				repo: &userRepositoryMock{
					GetUserByPhoneFunc: func(_ string) (*User, error) {
						return &User{ID: 1}, nil
					},
					CountRecoveryCodesByUserIn24hFunc: func(_ uint) (int, error) {
						return 0, errTestError
					},
				},
			},
			input: &User{
				Phone: "1",
			},
			wantErr: true,
		},
		{
			name: "error limit exceeded",
			params: params{
				repo: &userRepositoryMock{
					GetUserByPhoneFunc: func(_ string) (*User, error) {
						return &User{ID: 1}, nil
					},
					CountRecoveryCodesByUserIn24hFunc: func(_ uint) (int, error) {
						return 4, nil
					},
				},
			},
			input: &User{
				Phone: "1",
			},
			wantErr: true,
		},
		{
			name: "error wrong email",
			params: params{
				repo: &userRepositoryMock{
					GetUserByPhoneFunc: func(_ string) (*User, error) {
						return &User{
							ID:    1,
							Email: "a",
						}, nil
					},
					CountRecoveryCodesByUserIn24hFunc: func(_ uint) (int, error) {
						return 2, nil
					},
				},
			},
			input: &User{
				Phone: "1",
				Email: "b",
			},
			wantErr: true,
		},
		{
			name: "error create code",
			params: params{
				repo: &userRepositoryMock{
					GetUserByPhoneFunc: func(_ string) (*User, error) {
						return &User{
							ID:    1,
							Email: "a",
						}, nil
					},
					CountRecoveryCodesByUserIn24hFunc: func(_ uint) (int, error) {
						return 2, nil
					},
					CreateRecoverCodeFunc: func(_ *PasswordRecovery) error {
						return errTestError
					},
				},
			},
			input: &User{
				Phone: "1",
				Email: "a",
			},
			wantErr: true,
		},
		{
			name: "error send email",
			params: params{
				repo: &userRepositoryMock{
					GetUserByPhoneFunc: func(_ string) (*User, error) {
						return &User{
							ID:    1,
							Email: "a",
						}, nil
					},
					CountRecoveryCodesByUserIn24hFunc: func(_ uint) (int, error) {
						return 2, nil
					},
					CreateRecoverCodeFunc: func(_ *PasswordRecovery) error {
						return nil
					},
				},
				email: &mailSenderMock{
					SendRecoveryCodeEmailFunc: func(_, _, _ string) error {
						return errTestError
					},
				},
			},
			input: &User{
				Phone: "1",
				Email: "a",
			},
			wantErr: true,
		},
		{
			name: "ok",
			params: params{
				repo: &userRepositoryMock{
					GetUserByPhoneFunc: func(_ string) (*User, error) {
						return &User{
							ID:    1,
							Email: "a",
						}, nil
					},
					CountRecoveryCodesByUserIn24hFunc: func(_ uint) (int, error) {
						return 2, nil
					},
					CreateRecoverCodeFunc: func(_ *PasswordRecovery) error {
						return nil
					},
				},
				email: &mailSenderMock{
					SendRecoveryCodeEmailFunc: func(_, _, _ string) error {
						return nil
					},
				},
			},
			input: &User{
				Phone: "1",
				Email: "a",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(defaultLogger, tt.params.repo, tt.params.verifyRegCode, tt.params.maxMembers, tt.params.email)
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
		repo          userRepository
	}

	type input struct {
		code string
		u    *UserUpdate
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
				repo: &userRepositoryMock{
					GetRecoveryCodeFunc: func(_ *PasswordRecovery) (*PasswordRecovery, error) {
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
				repo: &userRepositoryMock{
					GetRecoveryCodeFunc: func(_ *PasswordRecovery) (*PasswordRecovery, error) {
						return &PasswordRecovery{
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
				repo: &userRepositoryMock{
					GetRecoveryCodeFunc: func(_ *PasswordRecovery) (*PasswordRecovery, error) {
						return &PasswordRecovery{
							UserID:    1,
							CreatedAt: func(t time.Time) *time.Time { return &t }(time.Now()),
						}, nil
					},
					GetUserByIDFunc: func(_ uint) (*User, error) {
						return nil, errTestError
					},
				},
			},
			input: input{
				code: "1",
				u:    &UserUpdate{},
			},
			wantErr: true,
		},
		{
			name: "error on reset password",
			params: params{
				repo: &userRepositoryMock{
					GetRecoveryCodeFunc: func(_ *PasswordRecovery) (*PasswordRecovery, error) {
						return &PasswordRecovery{
							UserID:    1,
							CreatedAt: func(t time.Time) *time.Time { return &t }(time.Now()),
						}, nil
					},
					GetUserByIDFunc: func(_ uint) (*User, error) {
						return nil, nil
					},
					ResetPasswordFunc: func(_ uint, _ *UserUpdate) error {
						return errTestError
					},
				},
			},
			input: input{
				code: "1",
				u: &UserUpdate{
					NewPassword: func(s string) *string { return &s }("1"),
				},
			},
			wantErr: true,
		},
		{
			name: "ok",
			params: params{
				repo: &userRepositoryMock{
					GetRecoveryCodeFunc: func(_ *PasswordRecovery) (*PasswordRecovery, error) {
						return &PasswordRecovery{
							UserID:    1,
							CreatedAt: func(t time.Time) *time.Time { return &t }(time.Now()),
						}, nil
					},
					GetUserByIDFunc: func(_ uint) (*User, error) {
						return nil, nil
					},
					ResetPasswordFunc: func(_ uint, _ *UserUpdate) error {
						return nil
					},
				},
			},
			input: input{
				code: "1",
				u: &UserUpdate{
					NewPassword: func(s string) *string { return &s }("1"),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(defaultLogger, tt.params.repo, tt.params.verifyRegCode, tt.params.maxMembers, nil)
			err := s.ResetPassword(context.Background(), tt.input.code, tt.input.u)
			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
