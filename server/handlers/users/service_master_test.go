package users

import (
	"context"
	"reflect"
	"testing"
)

func TestService_DeleteFamilyMember(t *testing.T) {
	type fields struct {
		verifyRegCode bool
		maxMembers    int
		repo          userRepository
	}

	type input struct {
		ownerID  uint
		memberID uint
	}

	tests := []struct {
		name    string
		fields  fields
		input   *input
		wantErr bool
	}{
		{
			name: "error getting user",
			fields: fields{
				verifyRegCode: false,
				maxMembers:    0,
				repo: &userRepositoryMock{
					GetUserByIDFunc: func(_ uint) (*User, error) {
						return nil, errTestError
					},
				},
			},
			input:   &input{memberID: 1},
			wantErr: true,
		},
		{
			name: "error wrong owner or member id",
			fields: fields{
				verifyRegCode: false,
				maxMembers:    0,
				repo: &userRepositoryMock{
					GetUserByIDFunc: func(_ uint) (*User, error) {
						var id uint = 2
						return &User{ParentID: &id}, nil
					},
				},
			},
			input:   &input{memberID: 1, ownerID: 3},
			wantErr: true,
		},
		{
			name: "error delete user",
			fields: fields{
				verifyRegCode: false,
				maxMembers:    0,
				repo: &userRepositoryMock{
					GetUserByIDFunc: func(_ uint) (*User, error) {
						var id uint = 2
						return &User{ParentID: &id}, nil
					},
					DeleteUserFunc: func(_ *User) error {
						return errTestError
					},
				},
			},
			input:   &input{memberID: 1, ownerID: 2},
			wantErr: true,
		},
		{
			name: "ok",
			fields: fields{
				verifyRegCode: false,
				maxMembers:    0,
				repo: &userRepositoryMock{
					GetUserByIDFunc: func(_ uint) (*User, error) {
						var id uint = 2
						return &User{ParentID: &id}, nil
					},
					DeleteUserFunc: func(_ *User) error {
						return nil
					},
				},
			},
			input:   &input{memberID: 1, ownerID: 2},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(defaultLogger, tt.fields.repo, tt.fields.verifyRegCode, tt.fields.maxMembers)
			err := s.DeleteFamilyMember(context.Background(), tt.input.ownerID, tt.input.memberID)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteFamilyMember() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestService_ListFamilyMembers(t *testing.T) {
	type fields struct {
		verifyRegCode bool
		maxMembers    int
		repo          userRepository
	}

	tests := []struct {
		name    string
		fields  fields
		input   uint
		wantErr bool
		want    []*User
	}{
		{
			name: "error getting family members",
			fields: fields{
				verifyRegCode: false,
				maxMembers:    0,
				repo: &userRepositoryMock{
					GetFamilyMembersFunc: func(_ uint) ([]*User, error) {
						return nil, errTestError
					},
				},
			},
			input:   0,
			wantErr: true,
		},
		{
			name: "empty family members list",
			fields: fields{
				verifyRegCode: false,
				maxMembers:    0,
				repo: &userRepositoryMock{
					GetFamilyMembersFunc: func(_ uint) ([]*User, error) {
						return nil, nil
					},
				},
			},
			input:   0,
			wantErr: false,
			want:    nil,
		},
		{
			name: "ok",
			fields: fields{
				verifyRegCode: false,
				maxMembers:    0,
				repo: &userRepositoryMock{
					GetFamilyMembersFunc: func(_ uint) ([]*User, error) {
						return []*User{
							{
								ID:        1,
								Phone:     "1",
								FirstName: "1",
								LastName:  "1",
								Active:    false,
								RegCode:   "1",
							},
						}, nil
					},
				},
			},
			input:   0,
			wantErr: false,
			want: []*User{
				{
					ID:        1,
					FirstName: "1",
					LastName:  "1",
					Phone:     "1",
					RegCode:   "1",
					Active:    false,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(defaultLogger, tt.fields.repo, tt.fields.verifyRegCode, tt.fields.maxMembers)
			got, err := s.ListFamilyMembers(context.Background(), tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListFamilyMembers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListFamilyMembers() got = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestService_AddFamilyMember(t *testing.T) {
	type fields struct {
		verifyRegCode bool
		maxMembers    int
		repo          userRepository
	}

	pid := uint(1)
	testRequest := &User{
		ParentID: &pid,
		Phone:    "1",
	}

	tests := []struct {
		name    string
		fields  fields
		input   *User
		wantErr bool
		want    *User
	}{
		{
			name: "error wrong owner",
			fields: fields{
				verifyRegCode: false,
				maxMembers:    0,
				repo: &userRepositoryMock{
					GetUserByIDFunc: func(_ uint) (*User, error) {
						return nil, errTestError
					},
				},
			},
			input:   testRequest,
			wantErr: true,
		},
		{
			name: "error getting family members",
			fields: fields{
				verifyRegCode: false,
				maxMembers:    1,
				repo: &userRepositoryMock{
					GetUserByIDFunc: func(_ uint) (*User, error) {
						return &User{ID: 1}, nil
					},
					GetFamilyMembersFunc: func(_ uint) ([]*User, error) {
						return nil, errTestError
					},
				},
			},
			input:   testRequest,
			wantErr: true,
		},
		{
			name: "error too much family members",
			fields: fields{
				verifyRegCode: false,
				maxMembers:    1,
				repo: &userRepositoryMock{
					GetUserByIDFunc: func(_ uint) (*User, error) {
						return &User{ID: 1}, nil
					},
					GetFamilyMembersFunc: func(_ uint) ([]*User, error) {
						return []*User{
							{},
							{},
						}, nil
					},
				},
			},
			input:   testRequest,
			wantErr: true,
		},
		{
			name: "error getting user by phone",
			fields: fields{
				verifyRegCode: false,
				maxMembers:    1,
				repo: &userRepositoryMock{
					GetUserByIDFunc: func(_ uint) (*User, error) {
						return &User{ID: 1}, nil
					},
					GetFamilyMembersFunc: func(_ uint) ([]*User, error) {
						return nil, nil
					},
					GetUserByPhoneFunc: func(_ string) (*User, error) {
						return nil, errTestError
					},
				},
			},
			input:   testRequest,
			wantErr: true,
		},
		{
			name: "error user exists",
			fields: fields{
				verifyRegCode: false,
				maxMembers:    1,
				repo: &userRepositoryMock{
					GetUserByIDFunc: func(_ uint) (*User, error) {
						return &User{ID: 1}, nil
					},
					GetFamilyMembersFunc: func(_ uint) ([]*User, error) {
						return nil, nil
					},
					GetUserByPhoneFunc: func(_ string) (*User, error) {
						return &User{}, nil
					},
				},
			},
			input:   testRequest,
			wantErr: true,
		},
		{
			name: "error getting reg code",
			fields: fields{
				verifyRegCode: false,
				maxMembers:    1,
				repo: &userRepositoryMock{
					GetUserByIDFunc: func(_ uint) (*User, error) {
						return &User{ID: 1}, nil
					},
					GetFamilyMembersFunc: func(_ uint) ([]*User, error) {
						return nil, nil
					},
					GetUserByPhoneFunc: func(_ string) (*User, error) {
						return nil, nil
					},
					GetRegCodeFunc: func() (string, error) {
						return "", errTestError
					},
				},
			},
			input:   testRequest,
			wantErr: true,
		},
		{
			name: "error creating user",
			fields: fields{
				verifyRegCode: false,
				maxMembers:    1,
				repo: &userRepositoryMock{
					GetUserByIDFunc: func(_ uint) (*User, error) {
						return &User{ID: 1}, nil
					},
					GetFamilyMembersFunc: func(_ uint) ([]*User, error) {
						return nil, nil
					},
					GetUserByPhoneFunc: func(_ string) (*User, error) {
						return nil, nil
					},
					GetRegCodeFunc: func() (string, error) {
						return "123", nil
					},
					CreateUserFunc: func(_ *User) error {
						return errTestError
					},
				},
			},
			input:   testRequest,
			wantErr: true,
		},
		{
			name: "error using reg code",
			fields: fields{
				verifyRegCode: false,
				maxMembers:    1,
				repo: &userRepositoryMock{
					GetUserByIDFunc: func(_ uint) (*User, error) {
						return &User{ID: 1}, nil
					},
					GetFamilyMembersFunc: func(_ uint) ([]*User, error) {
						return nil, nil
					},
					GetUserByPhoneFunc: func(_ string) (*User, error) {
						return nil, nil
					},
					GetRegCodeFunc: func() (string, error) {
						return "123", nil
					},
					CreateUserFunc: func(_ *User) error {
						return nil
					},
					UseRegCodeFunc: func(_ string) error {
						return errTestError
					},
				},
			},
			input:   testRequest,
			wantErr: true,
		},
		{
			name: "ok",
			fields: fields{
				verifyRegCode: false,
				maxMembers:    1,
				repo: &userRepositoryMock{
					GetUserByIDFunc: func(_ uint) (*User, error) {
						return &User{ID: 1}, nil
					},
					GetFamilyMembersFunc: func(_ uint) ([]*User, error) {
						return nil, nil
					},
					GetUserByPhoneFunc: func(_ string) (*User, error) {
						return nil, nil
					},
					GetRegCodeFunc: func() (string, error) {
						return "123", nil
					},
					CreateUserFunc: func(_ *User) error {
						return nil
					},
					UseRegCodeFunc: func(_ string) error {
						return nil
					},
				},
			},
			input:   testRequest,
			wantErr: false,
			want:    &User{RegCode: "123", Phone: "1", Role: defaultUserRole, ParentID: &pid},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(defaultLogger, tt.fields.repo, tt.fields.verifyRegCode, tt.fields.maxMembers)
			got, err := s.AddFamilyMember(context.Background(), tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddFamilyMember() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil {
				got.Password = ""
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddFamilyMember() got = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestService_familyMemberRegister(t *testing.T) {
	type fields struct {
		verifyRegCode bool
		maxMembers    int
		repo          userRepository
	}

	var pid uint = 2

	tests := []struct {
		name    string
		fields  fields
		input   *User
		wantErr bool
		want    *User
	}{
		{
			name: "error wrong reg codes",
			fields: fields{
				verifyRegCode: false,
				maxMembers:    0,
				repo: &userRepositoryMock{
					GetUserByPhoneFunc: func(_ string) (*User, error) {

						return &User{ParentID: &pid, RegCode: "111"}, nil
					},
				},
			},
			input: &User{
				RegCode:  "123",
				Phone:    "123",
				ParentID: &pid,
			},
			wantErr: true,
		},
		{
			name: "error user already registered",
			fields: fields{
				verifyRegCode: false,
				maxMembers:    0,
				repo: &userRepositoryMock{
					GetUserByPhoneFunc: func(_ string) (*User, error) {
						return &User{ParentID: &pid, Active: true, RegCode: "123"}, nil
					},
				},
			},
			input: &User{
				RegCode:  "123",
				Phone:    "123",
				ParentID: &pid,
			},
			wantErr: true,
		},
		{
			name: "error getting user by id",
			fields: fields{
				verifyRegCode: false,
				maxMembers:    0,
				repo: &userRepositoryMock{
					GetUserByPhoneFunc: func(_ string) (*User, error) {
						return &User{ParentID: &pid, Active: false, RegCode: "123"}, nil
					},
					GetUserByIDFunc: func(_ uint) (*User, error) {
						return nil, errTestError
					},
				},
			},
			input: &User{
				RegCode:  "123",
				Phone:    "123",
				ParentID: &pid,
			},
			wantErr: true,
		},
		{
			name: "error member entered wrong address",
			fields: fields{
				verifyRegCode: false,
				maxMembers:    0,
				repo: &userRepositoryMock{
					GetUserByPhoneFunc: func(_ string) (*User, error) {
						return &User{ParentID: &pid, Active: false, RegCode: "123"}, nil
					},
					GetUserByIDFunc: func(_ uint) (*User, error) {
						return &User{Apartment: 5}, nil
					},
				},
			},
			input: &User{
				RegCode:   "123",
				Phone:     "123",
				ParentID:  &pid,
				Apartment: 1,
			},
			wantErr: true,
		},
		{
			name: "error updating user",
			fields: fields{
				verifyRegCode: false,
				maxMembers:    0,
				repo: &userRepositoryMock{
					GetUserByPhoneFunc: func(_ string) (*User, error) {
						return &User{ParentID: &pid, Active: false, RegCode: "123"}, nil
					},
					GetUserByIDFunc: func(_ uint) (*User, error) {
						return &User{}, nil
					},
					UpdateUserFunc: func(_ *UserUpdate) error {
						return errTestError
					},
				},
			},
			input: &User{
				RegCode:  "123",
				Phone:    "123",
				ParentID: &pid,
			},
			wantErr: true,
		},
		{
			name: "ok",
			fields: fields{
				verifyRegCode: false,
				maxMembers:    0,
				repo: &userRepositoryMock{
					GetUserByPhoneFunc: func(_ string) (*User, error) {
						return &User{ID: 5, ParentID: &pid, Active: false, Phone: "123", RegCode: "123"}, nil
					},
					GetUserByIDFunc: func(_ uint) (*User, error) {
						return &User{}, nil
					},
					UpdateUserFunc: func(u *UserUpdate) error {
						return nil
					},
				},
			},
			input: &User{
				RegCode:  "123",
				Phone:    "123",
				ParentID: &pid,
			},
			wantErr: false,
			want: &User{
				ID:       5,
				Phone:    "123",
				RegCode:  "123",
				ParentID: &pid,
				Active:   true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(defaultLogger, tt.fields.repo, tt.fields.verifyRegCode, tt.fields.maxMembers)
			got, err := s.Register(context.Background(), tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("RegisterFamilyMember() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil {
				got.Password = ""
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RegisterFamilyMember() got = %#v, want %#v", got, tt.want)
			}
		})
	}
}
