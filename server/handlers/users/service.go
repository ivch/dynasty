package users

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"github.com/ivch/dynasty/common/errs"
	"github.com/ivch/dynasty/common/logger"
)

type userRepository interface {
	GetUserByID(id uint) (*User, error)
	GetUserByPhone(phone string) (*User, error)
	CreateUser(user *User) error
	UpdateUser(u *UserUpdate) error
	DeleteUser(u *User) error
	ValidateRegCode(code string) error
	UseRegCode(code string) error
	GetRegCode() (string, error)
	GetFamilyMembers(ownerID uint) ([]*User, error)
	FindUserByApartment(building uint, apt uint) (*User, error)
}

type Service struct {
	repo          userRepository
	membersLimit  int
	verifyRegCode bool
	log           logger.Logger
}

func New(log logger.Logger, repo userRepository, verifyRegCode bool, membersLimit int) *Service {
	s := Service{
		repo:          repo,
		membersLimit:  membersLimit,
		verifyRegCode: verifyRegCode,
		log:           log,
	}

	return &s
}

func (s *Service) UserByID(_ context.Context, id uint) (*User, error) {
	u, err := s.repo.GetUserByID(id)
	if err != nil {
		s.log.Error("error getting user from db: %w", err)
		return nil, err
	}
	return u, nil
}

func (s *Service) UserByPhoneAndPassword(_ context.Context, phone, password string) (*User, error) {
	u, err := s.repo.GetUserByPhone(phone)
	if err != nil {
		s.log.Error("error getting user from db: %w", err)
		return nil, err
	}

	if u == nil {
		return nil, errs.UserNotFound
	}

	if err := comparePasswords(u.Password, password); err != nil {
		s.log.Error("error comparing hash: %w", err)
		return nil, err
	}

	return u, nil
}

func (s *Service) Register(ctx context.Context, r *User) (*User, error) {
	user, err := s.repo.GetUserByPhone(r.Phone)
	if err != nil {
		s.log.Error("error getting user by phone: %w", err)
		return nil, err
	}

	if user != nil {
		if user.ParentID != nil {
			return s.registerFamilyMember(ctx, r, user)
		}
		return nil, errs.FamilyMemberPhoneExists
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

func (s *Service) Update(ctx context.Context, r *UserUpdate) error {
	u, err := s.UserByID(ctx, r.ID)
	if err != nil {
		return err
	}

	if r.NewPassword == nil {
		return s.repo.UpdateUser(r)
	}

	if r.Password == nil {
		return errs.InvalidCredentials
	}

	if err := comparePasswords(u.Password, *r.Password); err != nil {
		return errs.InvalidCredentials
	}

	// todo in case of password change delete current user session and invalidate refresh token
	pwd, err := hashAndSalt(*r.NewPassword)
	if err != nil {
		s.log.Error("error hashing password: %w", err)
		return err
	}

	r.Password = &pwd

	return s.repo.UpdateUser(r)
}

func hashAndSalt(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func comparePasswords(p1, p2 string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(p1), []byte(p2)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return errs.InvalidCredentials
		}
		return err
	}
	return nil
}
