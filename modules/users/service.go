package users

import (
	"context"

	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"

	"github.com/ivch/dynasty/models/dto"
	"github.com/ivch/dynasty/models/entities"
)

type Service interface {
	Register(ctx context.Context, req *dto.UserRegisterRequest) (*dto.UserRegisterResponse, error)
	UserByPhoneAndPassword(ctx context.Context, phone, password string) (*dto.UserAuthResponse, error)
	UserByID(ctx context.Context, id uint) (*dto.UserByIDResponse, error)
}

type userRepository interface {
	GetUserByID(id uint) (*entities.User, error)
	GetUserByPhone(phone string) (*entities.User, error)
	CreateUser(user *entities.User) error
	DeleteUser(u *entities.User) error
	ValidateRegCode(code string) error
	UseRegCode(code string) error
}

type service struct {
	repo          userRepository
	verifyRegCode bool
}

func newService(log *zerolog.Logger, repo userRepository, verifyRegCode bool) Service {
	s := &service{
		repo:          repo,
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
		Building:  u.Building,
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
	}, nil
}

func (s *service) Register(_ context.Context, r *dto.UserRegisterRequest) (*dto.UserRegisterResponse, error) {
	u, err := s.repo.GetUserByPhone(r.Phone)
	if err != nil && err != entities.ErrUserNotFound {
		return nil, err
	}

	if u != nil {
		return nil, entities.ErrUserPhoneExists
	}

	if s.verifyRegCode {
		if err := s.repo.ValidateRegCode(r.Code); err != nil {
			return nil, err
		}
	}

	pwd, err := hashAndSalt(r.Password)
	if err != nil {
		return nil, err
	}

	usr := entities.User{
		Apartment:  r.Apartment,
		BuildingID: r.BuildingID,
		Email:      r.Email,
		Phone:      r.Phone,
		FirstName:  r.FirstName,
		LastName:   r.LastName,
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
