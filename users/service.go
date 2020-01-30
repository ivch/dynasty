package users

import (
	"context"
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

const defaultUserRole = 4

var (
	errUserNotFound       = errors.New("user not found")
	errUserPhoneExists    = errors.New("user with this phone number already exists")
	errInvalidCredentials = errors.New("invalid login credentials")
	errInvalidRegCode     = errors.New("registration code is invalid or used")
)

type Service interface {
	Register(ctx context.Context, req *userRegisterRequest) (*userRegisterResponse, error)
	UserByPhoneAndPassword(ctx context.Context, r *userByPhoneAndPasswordRequest) (*User, error)
	UserByID(ctx context.Context, id int) (*userByIDResponse, error)
}

type userRegisterRequest struct {
	Password   string `json:"password,omitempty" validate:"required"`
	Phone      string `json:"phone,omitempty" validate:"required"`
	FirstName  string `json:"first_name,omitempty" validate:"required"`
	LastName   string `json:"last_name,omitempty" validate:"required"`
	BuildingID int    `json:"building_id,omitempty" validate:"required"`
	Apartment  uint   `json:"apartment,omitempty" validate:"required"`
	Email      string `json:"email,omitempty" validate:"email"`
	Code       string `json:"code" validate:"required"`
}

type userByPhoneAndPasswordRequest struct {
	Phone    string `validate:"required"`
	Password string `validate:"required"`
}

type userByIDResponse struct {
	ID        uint     `json:"id"`
	Apartment uint     `json:"apartment"`
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Phone     string   `json:"phone"`
	Email     string   `json:"email"`
	Role      uint     `json:"role"`
	Building  Building `json:"building"`
}

type userRegisterResponse struct {
	ID    uint   `json:"id"`
	Phone string `json:"phone"`
}

type User struct {
	ID         uint `gorm:"primary_key"`
	Building   Building
	Apartment  uint   `json:"apartment,omitempty"`
	Email      string `json:"email,omitempty"`
	Password   string `json:"password,omitempty"`
	Phone      string `json:"phone,omitempty"`
	FirstName  string `json:"first_name,omitempty"`
	LastName   string `json:"last_name,omitempty"`
	Role       int    `json:"role,omitempty"`
	BuildingID int    `gorm:"building_id"`
}

func (User) TableName() string { return "users" }

type Building struct {
	ID      uint   `gorm:"primary_key" json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

func (Building) TableName() string { return "buildings" }

type service struct {
	db            *gorm.DB
	verifyRegCode bool
}

func New(log *zerolog.Logger, db *gorm.DB, verifyRegCode bool) Service {
	s := &service{
		db:            db,
		verifyRegCode: verifyRegCode,
	}
	svc := newLoggingMiddleware(log, s)

	return svc
}

func (s *service) UserByID(_ context.Context, id int) (*userByIDResponse, error) {
	var u User
	if err := s.db.Debug().Preload("Building").Where("id = ?", id).First(&u).Error; err != nil {
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

func (s *service) UserByPhoneAndPassword(ctx context.Context, r *userByPhoneAndPasswordRequest) (*User, error) {
	u, err := s.userByPhone(ctx, r.Phone)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(r.Password)); err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		// Password does not match!
		return nil, errInvalidCredentials
	}

	return u, nil
}

func (s *service) userByPhone(_ context.Context, phone string) (*User, error) {
	var u User
	if err := s.db.Where("phone = ?", phone).First(&u).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errUserNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (s *service) Register(ctx context.Context, r *userRegisterRequest) (*userRegisterResponse, error) {
	u, err := s.userByPhone(ctx, r.Phone)
	if err != nil && err != errUserNotFound {
		return nil, err
	}

	if u != nil {
		return nil, errUserPhoneExists
	}

	if s.verifyRegCode {
		var code struct{ Exists bool }
		if err := s.db.Raw("select exists(select id from reg_codes where code = ? and not used)", r.Code).Scan(&code).Error; err != nil {
			return nil, err
		}

		if !code.Exists {
			return nil, errInvalidRegCode
		}
	}

	pwd, err := s.hashAndSalt(r.Password)
	if err != nil {
		return nil, err
	}

	usr := User{
		Apartment:  r.Apartment,
		BuildingID: r.BuildingID,
		Email:      r.Email,
		Phone:      r.Phone,
		FirstName:  r.FirstName,
		LastName:   r.LastName,
		Password:   pwd,
		Role:       defaultUserRole,
	}

	if err := s.db.Debug().Create(&usr).Error; err != nil {
		return nil, err
	}

	if s.verifyRegCode {
		if err := s.db.Exec("update reg_codes set used = true where code = ?", r.Code).Error; err != nil {
			if err := s.db.Delete(&usr).Error; err != nil {
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

func (s *service) hashAndSalt(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
