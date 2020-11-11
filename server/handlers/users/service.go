package users

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"github.com/ivch/dynasty/common/logger"
	"github.com/ivch/dynasty/server/handlers/users/errs"
)

type Service interface {
	Register(ctx context.Context, u *User) (*User, error)
	// UserByPhoneAndPassword(ctx context.Context, phone, password string) (*User, error)
	UserByID(ctx context.Context, id uint) (*User, error)
	// AddFamilyMember(ctx context.Context, r *dto.AddFamilyMemberRequest) (*dto.AddFamilyMemberResponse, error)
	// ListFamilyMembers(ctx context.Context, id uint) (*dto.ListFamilyMembersResponse, error)
	// DeleteFamilyMember(ctx context.Context, r *dto.DeleteFamilyMemberRequest) error
}

type userRepository interface {
	GetUserByID(id uint) (*User, error)
	GetUserByPhone(phone string) (*User, error)
	CreateUser(user *User) error
	UpdateUser(u *User) error
	DeleteUser(u *User) error
	ValidateRegCode(code string) error
	UseRegCode(code string) error
	// GetRegCode() (string, error)
	// GetFamilyMembers(ownerID uint) ([]*entities.User, error)
	FindUserByApartment(building uint, apt uint) (*User, error)
}

type service struct {
	repo          userRepository
	membersLimit  int
	verifyRegCode bool
	log           logger.Logger
}

func New(log logger.Logger, repo userRepository, verifyRegCode bool, membersLimit int) Service {
	s := &service{
		repo:          repo,
		membersLimit:  membersLimit,
		verifyRegCode: verifyRegCode,
		log:           log,
	}

	return s
}

func (s *service) UserByID(_ context.Context, id uint) (*User, error) {
	u, err := s.repo.GetUserByID(id)
	if err != nil {
		s.log.Error("error getting user from db: %w", err)
		return nil, err
	}

	return u, nil
}

// func (s *service) UserByPhoneAndPassword(_ context.Context, phone, password string) (*entities.User, error) {
// 	u, err := s.repo.GetUserByPhone(phone)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	if u == nil {
// 		return nil, errs.UserNotFound
// 	}
//
// 	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
// 		if err == bcrypt.ErrMismatchedHashAndPassword {
// 			return nil, entities.ErrInvalidCredentials
// 		}
// 		return nil, err
// 	}
//
// 	return u, nil
// }

func (s *service) Register(ctx context.Context, r *User) (*User, error) {
	user, err := s.repo.GetUserByPhone(r.Phone)
	if err != nil && user != nil {
		s.log.Error("error getting user by phone: %w", err)
		return nil, err
	}

	if user != nil {
		if r.ParentID != nil {
			return s.registerFamilyMember(ctx, r, user)
		}
		// return nil, entities.ErrUserPhoneExists
	}

	if s.verifyRegCode {
		if err := s.repo.ValidateRegCode(r.RegCode); err != nil {
			s.log.Debug("validated reg code: %w", err)
			return nil, err
		}
	}

	m, err := s.repo.FindUserByApartment(r.BuildingID, r.Apartment)
	if err != nil {
		s.log.Error("error getting user by apt: %w", err)
		return nil, err
	}

	if m != nil {
		return nil, errs.MasterAccountExists
	}

	pwd, err := hashAndSalt(r.Password)
	if err != nil {
		s.log.Error("error hashing password: %w", err)
		return nil, err
	}

	u := User{
		Apartment:  r.Apartment,
		BuildingID: r.BuildingID,
		EntryID:    r.EntryID,
		Email:      r.Email,
		Phone:      r.Phone,
		FirstName:  r.FirstName,
		LastName:   r.LastName,
		Active:     true,
		Password:   pwd,
		Role:       defaultUserRole,
	}

	if err := s.repo.CreateUser(&u); err != nil {
		s.log.Error("error creating user: %w", err)
		return nil, err
	}

	if s.verifyRegCode {
		if err := s.repo.UseRegCode(r.RegCode); err != nil {
			if err := s.repo.DeleteUser(&u); err != nil {
				return nil, err
			}
			return nil, err
		}
	}

	return &u, nil
}

func hashAndSalt(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
