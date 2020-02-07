package users

import (
	"context"

	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"

	"github.com/ivch/dynasty/models"
)

type Service interface {
	Register(ctx context.Context, req *userRegisterRequest) (*userRegisterResponse, error)
	UserByPhoneAndPassword(ctx context.Context, phone, password string) (*userAuthResponse, error)
	UserByID(ctx context.Context, id uint) (*userByIDResponse, error)
}

type userRepository interface {
	GetUserByID(id uint) (*models.User, error)
	GetUserByPhone(phone string) (*models.User, error)
	CreateUser(user *models.User) error
	DeleteUser(u *models.User) error
	ValidateRegCode(code string) error
	UseRegCode(code string) error
}

type userRegisterRequest struct {
	Password   string `json:"password,omitempty" validate:"required,min=6"`
	Phone      string `json:"phone,omitempty" validate:"required"`
	FirstName  string `json:"first_name,omitempty" validate:"required"`
	LastName   string `json:"last_name,omitempty" validate:"required"`
	BuildingID int    `json:"building_id,omitempty" validate:"required"`
	Apartment  uint   `json:"apartment,omitempty" validate:"required"`
	Email      string `json:"email,omitempty" validate:"email"`
	Code       string `json:"code"`
}

type userAuthResponse struct {
	ID        uint   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      uint   `json:"role"`
}

type userByIDResponse struct {
	ID        uint            `json:"id"`
	Apartment uint            `json:"apartment"`
	FirstName string          `json:"first_name"`
	LastName  string          `json:"last_name"`
	Phone     string          `json:"phone"`
	Email     string          `json:"email"`
	Role      uint            `json:"role"`
	Building  models.Building `json:"building"`
}

type userRegisterResponse struct {
	ID    uint   `json:"id"`
	Phone string `json:"phone"`
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

func (s *service) UserByID(_ context.Context, id uint) (*userByIDResponse, error) {
	u, err := s.repo.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	return &userByIDResponse{
		ID:        u.ID,
		Apartment: u.Apartment,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Phone:     u.Phone,
		Email:     u.Email,
		Building:  u.Building,
	}, nil
}

func (s *service) UserByPhoneAndPassword(_ context.Context, phone, password string) (*userAuthResponse, error) {
	u, err := s.repo.GetUserByPhone(phone)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return nil, models.ErrInvalidCredentials
		}
		return nil, err
	}

	return &userAuthResponse{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Role:      u.Role,
	}, nil
}

func (s *service) Register(_ context.Context, r *userRegisterRequest) (*userRegisterResponse, error) {
	u, err := s.repo.GetUserByPhone(r.Phone)
	if err != nil && err != models.ErrUserNotFound {
		return nil, err
	}

	if u != nil {
		return nil, models.ErrUserPhoneExists
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

	usr := models.User{
		Apartment:  r.Apartment,
		BuildingID: r.BuildingID,
		Email:      r.Email,
		Phone:      r.Phone,
		FirstName:  r.FirstName,
		LastName:   r.LastName,
		Password:   pwd,
		Role:       models.DefaultUserRole,
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

	return &userRegisterResponse{
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
