package users

import (
	"context"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/ivch/dynasty/common"
	"github.com/ivch/dynasty/common/errs"
	"github.com/ivch/dynasty/common/logger"
)

type userRepository interface {
	GetUserByID(id uint) (*User, error)
	GetUserByPhone(phone string) (*User, error)
	GetUserByEmail(email string) (*User, error)
	CreateUser(user *User) error
	UpdateUser(u *UserUpdate) error
	DeleteUser(u *User) error
	ValidateRegCode(code string) error
	UseRegCode(code string) error
	GetRegCode() (string, error)
	GetFamilyMembers(ownerID uint) ([]*User, error)
	FindUserByApartment(building uint, apt uint) (*User, error)
	CreateRecoverCode(c *PasswordRecovery) error
	CountRecoveryCodesByUserIn24h(userID uint) (int, error)
	GetRecoveryCode(c *PasswordRecovery) (*PasswordRecovery, error)
	ResetPassword(codeID uint, req *UserUpdate) error
}

type mailSender interface {
	SendRecoveryCodeEmail(to, username, code string) error
}

type Service struct {
	repo          userRepository
	membersLimit  int
	verifyRegCode bool
	email         mailSender
	log           logger.Logger
}

func New(log logger.Logger, repo userRepository, verifyRegCode bool, membersLimit int, email mailSender) *Service {
	s := Service{
		repo:          repo,
		membersLimit:  membersLimit,
		verifyRegCode: verifyRegCode,
		log:           log,
		email:         email,
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

// nolint: gocyclo,funlen
func (s *Service) Register(ctx context.Context, r *User) (*User, error) {
	user, err := s.repo.GetUserByEmail(r.Email)
	if err != nil {
		s.log.Error("error getting user by email: %w", err)
		return nil, err
	}

	if user != nil {
		return nil, errs.EmailAlreadyExists
	}

	user, err = s.repo.GetUserByPhone(r.Phone)
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

	m, err := s.repo.FindUserByApartment(r.BuildingID, r.Apartment)
	if err != nil {
		s.log.Error("error getting user by apt: %w", err)
		return nil, err
	}

	if m != nil && m.Role != predefinedUserRole {
		return nil, errs.MasterAccountExists
	}

	if s.verifyRegCode && (m != nil && m.Role != predefinedUserRole) {
		if err := s.repo.ValidateRegCode(r.RegCode); err != nil {
			s.log.Debug("validated reg code: %w", err)
			return nil, err
		}
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

	if m != nil && m.Role == predefinedUserRole {
		if m.RegCode != r.RegCode {
			return nil, errs.RegCodeInvalid
		}
		return s.registerPredefinedUser(ctx, m, &u)
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

func (s *Service) registerPredefinedUser(_ context.Context, predefinedUser, newUser *User) (*User, error) {
	upd := UserUpdate{
		ID:        predefinedUser.ID,
		Email:     &newUser.Email,
		Phone:     &newUser.Phone,
		Password:  &newUser.Password,
		Role:      &newUser.Role,
		FirstName: &newUser.FirstName,
		LastName:  &newUser.LastName,
		Active:    &newUser.Active,
	}

	if err := s.repo.UpdateUser(&upd); err != nil {
		return nil, err
	}

	newUser.ID = predefinedUser.ID

	return newUser, nil
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

func (s *Service) RecoveryCode(_ context.Context, r *User) error {
	u, err := s.repo.GetUserByPhone(r.Phone)
	if err != nil {
		return err
	}

	if u == nil {
		return errs.UserNotFound
	}

	cnt, err := s.repo.CountRecoveryCodesByUserIn24h(u.ID)
	if err != nil {
		s.log.Error("error counting codes: %w", err)
		return err
	}

	if cnt >= 3 {
		return errs.PasswordRecoveryLimit
	}

	if u.Email != r.Email {
		return errs.EmailInvalid
	}

	code := PasswordRecovery{
		UserID: u.ID,
		Code:   strings.ToUpper(common.RandomString(10)),
		Active: true,
	}

	if err := s.repo.CreateRecoverCode(&code); err != nil {
		return err
	}

	if err := s.email.SendRecoveryCodeEmail(u.Email, fmt.Sprintf("%s %s", u.FirstName, u.LastName), code.Code); err != nil {
		s.log.Error("error sending recovery email: %w", err)
		return err
	}

	return nil
}

func (s *Service) ResetPassword(ctx context.Context, code string, r *UserUpdate) error {
	c, err := s.repo.GetRecoveryCode(&PasswordRecovery{
		Code:   code,
		Active: true,
	})
	if err != nil {
		return errs.BadRecoveryCode
	}

	if time.Now().After(c.CreatedAt.Add(3 * time.Hour)) {
		return errs.RecoveryCodeOutdated
	}

	r.ID = c.UserID

	if _, err := s.repo.GetUserByID(r.ID); err != nil {
		return err
	}

	// todo in case of password change delete current user session and invalidate refresh token
	pwd, err := hashAndSalt(*r.NewPassword)
	if err != nil {
		s.log.Error("error hashing password: %w", err)
		return err
	}

	r.Password = &pwd

	return s.repo.ResetPassword(c.ID, r)
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
