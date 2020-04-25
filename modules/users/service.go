package users

import (
	"context"
	"errors"

	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"

	"github.com/ivch/dynasty/models/dto"
	"github.com/ivch/dynasty/models/entities"
)

type Service interface {
	Register(ctx context.Context, req *dto.UserRegisterRequest) (*dto.UserRegisterResponse, error)
	UserByPhoneAndPassword(ctx context.Context, phone, password string) (*dto.UserAuthResponse, error)
	UserByID(ctx context.Context, id uint) (*dto.UserByIDResponse, error)
	AddFamilyMember(ctx context.Context, r *dto.AddFamilyMemberRequest) (*dto.AddFamilyMemberResponse, error)
	ListFamilyMembers(ctx context.Context, id uint) (*dto.ListFamilyMembersResponse, error)
	DeleteFamilyMember(ctx context.Context, r *dto.DeleteFamilyMemberRequest) error
}

type userRepository interface {
	GetUserByID(id uint) (*entities.User, error)
	GetUserByPhone(phone string) (*entities.User, error)
	CreateUser(user *entities.User) error
	UpdateUser(u *entities.User) error
	DeleteUser(u *entities.User) error
	ValidateRegCode(code string) error
	UseRegCode(code string) error
	GetRegCode() (string, error)
	GetFamilyMembers(ownerID uint) ([]*entities.User, error)
	FindUserByApartment(building uint, apt uint) (*entities.User, error)
}

type service struct {
	repo          userRepository
	membersLimit  int
	verifyRegCode bool
}

func newService(log *zerolog.Logger, repo userRepository, verifyRegCode bool, membersLimit int) Service {
	s := &service{
		repo:          repo,
		membersLimit:  membersLimit,
		verifyRegCode: verifyRegCode,
	}
	svc := newLoggingMiddleware(log, s)

	return svc
}

func (s *service) UserByID(_ context.Context, id uint) (*dto.UserByIDResponse, error) {
	u, err := s.repo.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	return &dto.UserByIDResponse{
		ID:        u.ID,
		Apartment: u.Apartment,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Phone:     u.Phone,
		Email:     u.Email,
		Building:  &u.Building,
		Entry:     &u.Entry,
	}, nil
}

func (s *service) UserByPhoneAndPassword(_ context.Context, phone, password string) (*dto.UserAuthResponse, error) {
	u, err := s.repo.GetUserByPhone(phone)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return nil, entities.ErrInvalidCredentials
		}
		return nil, err
	}

	return &dto.UserAuthResponse{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Role:      u.Role,
		Active:    u.Active,
	}, nil
}

func (s *service) Register(ctx context.Context, r *dto.UserRegisterRequest) (*dto.UserRegisterResponse, error) {
	u, err := s.repo.GetUserByPhone(r.Phone)
	if err != nil && !errors.Is(err, entities.ErrUserNotFound) {
		return nil, err
	}

	if u != nil {
		if u.ParentID != nil {
			return s.registerFamilyMember(ctx, r, u)
		}
		return nil, entities.ErrUserPhoneExists
	}

	if s.verifyRegCode {
		if err := s.repo.ValidateRegCode(r.Code); err != nil {
			return nil, err
		}
	}

	m, err := s.repo.FindUserByApartment(r.BuildingID, r.Apartment)
	if err != nil {
		return nil, err
	}

	if m != nil {
		return nil, errMasterAccountExists
	}

	pwd, err := hashAndSalt(r.Password)
	if err != nil {
		return nil, err
	}

	usr := entities.User{
		Apartment:  r.Apartment,
		BuildingID: r.BuildingID,
		EntryID:    r.EntryID,
		Email:      r.Email,
		Phone:      r.Phone,
		FirstName:  r.FirstName,
		LastName:   r.LastName,
		Active:     true,
		Password:   pwd,
		Role:       entities.DefaultUserRole,
	}

	if err := s.repo.CreateUser(&usr); err != nil {
		return nil, err
	}

	if s.verifyRegCode {
		if err := s.repo.UseRegCode(r.Code); err != nil {
			if err := s.repo.DeleteUser(&usr); err != nil {
				return nil, err
			}
			return nil, err
		}
	}

	return &dto.UserRegisterResponse{
		ID:    usr.ID,
		Phone: usr.Phone,
	}, nil
}

func hashAndSalt(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
