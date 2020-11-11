package users

import (
	"context"
	"reflect"
	"testing"

	"github.com/ivch/dynasty/models/dto"
	"github.com/ivch/dynasty/models/entities"
)

func TestService_DeleteFamilyMember(t *testing.T) {
	type fields struct {
		verifyRegCode bool
		maxMembers    int
		repo          userRepository
	}

	tests := []struct {
		name    string
		fields  fields
		args    *dto.DeleteFamilyMemberRequest
		wantErr bool
		want    *dto.ListFamilyMembersResponse
	}{
		{
			name: "error getting user",
			fields: fields{
				verifyRegCode: false,
				maxMembers:    0,
				repo: &userRepositoryMock{
					GetUserByIDFunc: func(_ uint) (*entities.User, error) {
						return nil, errTestError
					},
				},
			},
			args:    &dto.DeleteFamilyMemberRequest{MemberID: 1},
			wantErr: true,
		},
		{
			name: "error wrong owner or member id",
			fields: fields{
				verifyRegCode: false,
				maxMembers:    0,
				repo: &userRepositoryMock{
					GetUserByIDFunc: func(_ uint) (*entities.User, error) {
						var id uint = 2
						return &entities.User{ParentID: &id}, nil
					},
				},
			},
			args:    &dto.DeleteFamilyMemberRequest{MemberID: 1, OwnerID: 3},
			wantErr: true,
		},
		{
			name: "error delete user",
			fields: fields{
				verifyRegCode: false,
				maxMembers:    0,
				repo: &userRepositoryMock{
					GetUserByIDFunc: func(_ uint) (*entities.User, error) {
						var id uint = 2
						return &entities.User{ParentID: &id}, nil
					},
					DeleteUserFunc: func(_ *entities.User) error {
						return errTestError
					},
				},
			},
			args:    &dto.DeleteFamilyMemberRequest{MemberID: 1, OwnerID: 2},
			wantErr: true,
		},
		{
			name: "ok",
			fields: fields{
				verifyRegCode: false,
				maxMembers:    0,
				repo: &userRepositoryMock{
					GetUserByIDFunc: func(_ uint) (*entities.User, error) {
						var id uint = 2
						return &entities.User{ParentID: &id}, nil
					},
					DeleteUserFunc: func(_ *entities.User) error {
						return nil
					},
				},
			},
			args:    &dto.DeleteFamilyMemberRequest{MemberID: 1, OwnerID: 2},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newService(defaultLogger, tt.fields.repo, tt.fields.verifyRegCode, tt.fields.maxMembers)
			err := s.DeleteFamilyMember(context.Background(), tt.args)
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
		args    uint
		wantErr bool
		want    *dto.ListFamilyMembersResponse
	}{
		{
			name: "error getting family members",
			fields: fields{
				verifyRegCode: false,
				maxMembers:    0,
				repo: &userRepositoryMock{
					GetFamilyMembersFunc: func(_ uint) ([]*entities.User, error) {
						return nil, errTestError
					},
				},
			},
			args:    0,
			wantErr: true,
		},
		{
			name: "empty family members list",
			fields: fields{
				verifyRegCode: false,
				maxMembers:    0,
				repo: &userRepositoryMock{
					GetFamilyMembersFunc: func(_ uint) ([]*entities.User, error) {
						return nil, nil
					},
				},
			},
			args:    0,
			wantErr: false,
			want:    nil,
		},
		{
			name: "ok",
			fields: fields{
				verifyRegCode: false,
				maxMembers:    0,
				repo: &userRepositoryMock{
					GetFamilyMembersFunc: func(_ uint) ([]*entities.User, error) {
						return []*entities.User{
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
			args:    0,
			wantErr: false,
			want: &dto.ListFamilyMembersResponse{Data: []*dto.FamilyMember{
				{
					ID:     1,
					Name:   "1 1",
					Phone:  "1",
					Code:   "1",
					Active: false,
				},
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newService(defaultLogger, tt.fields.repo, tt.fields.verifyRegCode, tt.fields.maxMembers)
			got, err := s.ListFamilyMembers(context.Background(), tt.args)
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

	testRequest := &dto.AddFamilyMemberRequest{
		OwnerID: 1,
		Phone:   "1",
	}

	tests := []struct {
		name    string
		fields  fields
		args    *dto.AddFamilyMemberRequest
		wantErr bool
		want    *dto.AddFamilyMemberResponse
	}{
		{
			name: "error wrong owner",
			fields: fields{
				verifyRegCode: false,
				maxMembers:    0,
				repo: &userRepositoryMock{
					GetUserByIDFunc: func(_ uint) (*entities.User, error) {
						return nil, errTestError
					},
				},
			},
			args:    testRequest,
			wantErr: true,
		},
		{
			name: "error getting family members",
			fields: fields{
				verifyRegCode: false,
				maxMembers:    1,
				repo: &userRepositoryMock{
					GetUserByIDFunc: func(_ uint) (*entities.User, error) {
						return &entities.User{ID: 1}, nil
					},
					GetFamilyMembersFunc: func(_ uint) ([]*entities.User, error) {
						return nil, errTestError
					},
				},
			},
			args:    testRequest,
			wantErr: true,
		},
		{
			name: "error too much family members",
			fields: fields{
				verifyRegCode: false,
				maxMembers:    1,
				repo: &userRepositoryMock{
					GetUserByIDFunc: func(_ uint) (*entities.User, error) {
						return &entities.User{ID: 1}, nil
					},
					GetFamilyMembersFunc: func(_ uint) ([]*entities.User, error) {
						return []*entities.User{
							{},
							{},
						}, nil
					},
				},
			},
			args:    testRequest,
			wantErr: true,
		},
		{
			name: "error getting user by phone",
			fields: fields{
				verifyRegCode: false,
				maxMembers:    1,
				repo: &userRepositoryMock{
					GetUserByIDFunc: func(_ uint) (*entities.User, error) {
						return &entities.User{ID: 1}, nil
					},
					GetFamilyMembersFunc: func(_ uint) ([]*entities.User, error) {
						return nil, nil
					},
					GetUserByPhoneFunc: func(_ string) (*entities.User, error) {
						return nil, errTestError
					},
				},
			},
			args:    testRequest,
			wantErr: true,
		},
		{
			name: "error user exists",
			fields: fields{
				verifyRegCode: false,
				maxMembers:    1,
				repo: &userRepositoryMock{
					GetUserByIDFunc: func(_ uint) (*entities.User, error) {
						return &entities.User{ID: 1}, nil
					},
					GetFamilyMembersFunc: func(_ uint) ([]*entities.User, error) {
						return nil, nil
					},
					GetUserByPhoneFunc: func(_ string) (*entities.User, error) {
						return &entities.User{}, nil
					},
				},
			},
			args:    testRequest,
			wantErr: true,
		},
		{
			name: "error getting reg code",
			fields: fields{
				verifyRegCode: false,
				maxMembers:    1,
				repo: &userRepositoryMock{
					GetUserByIDFunc: func(_ uint) (*entities.User, error) {
						return &entities.User{ID: 1}, nil
					},
					GetFamilyMembersFunc: func(_ uint) ([]*entities.User, error) {
						return nil, nil
					},
					GetUserByPhoneFunc: func(_ string) (*entities.User, error) {
						return nil, nil
					},
					GetRegCodeFunc: func() (string, error) {
						return "", errTestError
					},
				},
			},
			args:    testRequest,
			wantErr: true,
		},
		{
			name: "error creating user",
			fields: fields{
				verifyRegCode: false,
				maxMembers:    1,
				repo: &userRepositoryMock{
					GetUserByIDFunc: func(_ uint) (*entities.User, error) {
						return &entities.User{ID: 1}, nil
					},
					GetFamilyMembersFunc: func(_ uint) ([]*entities.User, error) {
						return nil, nil
					},
					GetUserByPhoneFunc: func(_ string) (*entities.User, error) {
						return nil, nil
					},
					GetRegCodeFunc: func() (string, error) {
						return "123", nil
					},
					CreateUserFunc: func(_ *entities.User) error {
						return errTestError
					},
				},
			},
			args:    testRequest,
			wantErr: true,
		},
		{
			name: "error using reg code",
			fields: fields{
				verifyRegCode: false,
				maxMembers:    1,
				repo: &userRepositoryMock{
					GetUserByIDFunc: func(_ uint) (*entities.User, error) {
						return &entities.User{ID: 1}, nil
					},
					GetFamilyMembersFunc: func(_ uint) ([]*entities.User, error) {
						return nil, nil
					},
					GetUserByPhoneFunc: func(_ string) (*entities.User, error) {
						return nil, nil
					},
					GetRegCodeFunc: func() (string, error) {
						return "123", nil
					},
					CreateUserFunc: func(_ *entities.User) error {
						return nil
					},
					UseRegCodeFunc: func(_ string) error {
						return errTestError
					},
				},
			},
			args:    testRequest,
			wantErr: true,
		},
		{
			name: "ok",
			fields: fields{
				verifyRegCode: false,
				maxMembers:    1,
				repo: &userRepositoryMock{
					GetUserByIDFunc: func(_ uint) (*entities.User, error) {
						return &entities.User{ID: 1}, nil
					},
					GetFamilyMembersFunc: func(_ uint) ([]*entities.User, error) {
						return nil, nil
					},
					GetUserByPhoneFunc: func(_ string) (*entities.User, error) {
						return nil, nil
					},
					GetRegCodeFunc: func() (string, error) {
						return "123", nil
					},
					CreateUserFunc: func(_ *entities.User) error {
						return nil
					},
					UseRegCodeFunc: func(_ string) error {
						return nil
					},
				},
			},
			args:    testRequest,
			wantErr: false,
			want:    &dto.AddFamilyMemberResponse{Code: "123"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newService(defaultLogger, tt.fields.repo, tt.fields.verifyRegCode, tt.fields.maxMembers)
			got, err := s.AddFamilyMember(context.Background(), tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddFamilyMember() error = %v, wantErr %v", err, tt.wantErr)
				return
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

	tests := []struct {
		name    string
		fields  fields
		args    *dto.UserRegisterRequest
		wantErr bool
		want    *dto.UserRegisterResponse
	}{
		{
			name: "error wrong reg codes",
			fields: fields{
				verifyRegCode: false,
				maxMembers:    0,
				repo: &userRepositoryMock{
					GetUserByPhoneFunc: func(_ string) (*entities.User, error) {
						var id uint = 2
						return &entities.User{ParentID: &id, RegCode: "111"}, nil
					},
				},
			},
			args: &dto.UserRegisterRequest{
				Code:  "123",
				Phone: "123",
			},
			wantErr: true,
		},
		{
			name: "error user already registered",
			fields: fields{
				verifyRegCode: false,
				maxMembers:    0,
				repo: &userRepositoryMock{
					GetUserByPhoneFunc: func(_ string) (*entities.User, error) {
						var id uint = 2
						return &entities.User{ParentID: &id, Active: true, RegCode: "123"}, nil
					},
				},
			},
			args: &dto.UserRegisterRequest{
				Code:  "123",
				Phone: "123",
			},
			wantErr: true,
		},
		{
			name: "error getting user by id",
			fields: fields{
				verifyRegCode: false,
				maxMembers:    0,
				repo: &userRepositoryMock{
					GetUserByPhoneFunc: func(_ string) (*entities.User, error) {
						var id uint = 2
						return &entities.User{ParentID: &id, Active: false, RegCode: "123"}, nil
					},
					GetUserByIDFunc: func(_ uint) (*entities.User, error) {
						return nil, errTestError
					},
				},
			},
			args: &dto.UserRegisterRequest{
				Code:  "123",
				Phone: "123",
			},
			wantErr: true,
		},
		{
			name: "error member entered wrong address",
			fields: fields{
				verifyRegCode: false,
				maxMembers:    0,
				repo: &userRepositoryMock{
					GetUserByPhoneFunc: func(_ string) (*entities.User, error) {
						var id uint = 2
						return &entities.User{ParentID: &id, Active: false, RegCode: "123"}, nil
					},
					GetUserByIDFunc: func(_ uint) (*entities.User, error) {
						return &entities.User{Apartment: 5}, nil
					},
				},
			},
			args: &dto.UserRegisterRequest{
				Code:      "123",
				Phone:     "123",
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
					GetUserByPhoneFunc: func(_ string) (*entities.User, error) {
						var id uint = 2
						return &entities.User{ParentID: &id, Active: false, RegCode: "123"}, nil
					},
					GetUserByIDFunc: func(_ uint) (*entities.User, error) {
						return &entities.User{}, nil
					},
					UpdateUserFunc: func(_ *entities.User) error {
						return errTestError
					},
				},
			},
			args: &dto.UserRegisterRequest{
				Code:  "123",
				Phone: "123",
			},
			wantErr: true,
		},
		{
			name: "ok",
			fields: fields{
				verifyRegCode: false,
				maxMembers:    0,
				repo: &userRepositoryMock{
					GetUserByPhoneFunc: func(_ string) (*entities.User, error) {
						var id uint = 2
						return &entities.User{ID: 5, ParentID: &id, Active: false, Phone: "123", RegCode: "123"}, nil
					},
					GetUserByIDFunc: func(_ uint) (*entities.User, error) {
						return &entities.User{}, nil
					},
					UpdateUserFunc: func(u *entities.User) error {
						return nil
					},
				},
			},
			args: &dto.UserRegisterRequest{
				Code:  "123",
				Phone: "123",
			},
			wantErr: false,
			want: &dto.UserRegisterResponse{
				ID:    5,
				Phone: "123",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newService(defaultLogger, tt.fields.repo, tt.fields.verifyRegCode, tt.fields.maxMembers)
			got, err := s.Register(context.Background(), tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("RegisterFamilyMember() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RegisterFamilyMember() got = %#v, want %#v", got, tt.want)
			}
		})
	}
}
