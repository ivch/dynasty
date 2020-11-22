package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/ivch/dynasty/common/errs"
	"github.com/ivch/dynasty/common/logger"
	users "github.com/ivch/dynasty/server/handlers/users/transport"
)

type userService interface {
	UserByPhoneAndPassword(ctx context.Context, phone, password string) (*users.UserByIDResponse, error)
	UserByID(ctx context.Context, id uint) (*users.UserByIDResponse, error)
}

type authRepository interface {
	CreateSession(userID uint) (string, error)
	FindSessionByAccessToken(token string) (*Session, error)
	DeleteSessionByID(id string) error
}

type Service struct {
	log       logger.Logger
	uSrv      userService
	repo      authRepository
	jwtSecret string
}

func New(log logger.Logger, repo authRepository, uSrv userService, jwtSecret string) *Service {
	s := Service{
		log:       log,
		repo:      repo,
		uSrv:      uSrv,
		jwtSecret: jwtSecret,
	}

	return &s
}

func (s *Service) Gwfa(token string) (uint, error) {
	var myClaims Token
	t, err := jwt.ParseWithClaims(token, &myClaims, func(token *jwt.Token) (i interface{}, e error) {
		return []byte(s.jwtSecret), nil
	})
	if err != nil {
		return 0, fmt.Errorf("%s: %w", errs.FailedParsingToken, err)
	}

	claims, ok := t.Claims.(*Token)
	if !ok {
		return 0, errs.FailedParsingTokenClaims
	}

	return claims.ID, nil
}

func (s *Service) Refresh(ctx context.Context, token string) (*Tokens, error) {
	sess, err := s.repo.FindSessionByAccessToken(token)
	if err != nil {
		return nil, errs.NoSessionToRefresh
	}

	if err := s.repo.DeleteSessionByID(sess.ID); err != nil {
		s.log.Error("error deleting sessions: %w", err)
		return nil, err
	}

	if time.Now().Unix() > sess.ExpiresIn {
		return nil, errs.TokenExpired
	}

	u, err := s.uSrv.UserByID(ctx, sess.UserID)
	if err != nil {
		s.log.Error("failed to get user: %w", err)
		return nil, err
	}

	rt, err := s.repo.CreateSession(sess.UserID)
	if err != nil {
		s.log.Error("failed to create session: %w", err)
		return nil, err
	}

	at, err := s.generateAccessToken(u)
	if err != nil {
		s.log.Error("failed to create access token: %w", err)
		return nil, err
	}

	return &Tokens{
		AccessToken:  at,
		RefreshToken: rt,
	}, nil
}

func (s *Service) Login(ctx context.Context, phone, password string) (*Tokens, error) {
	u, err := s.uSrv.UserByPhoneAndPassword(ctx, phone, password)
	if err != nil {
		s.log.Error("error finding users: %w", err.Error())
		return nil, errs.InvalidCredentials
	}

	if !u.Active {
		return nil, errs.UserIsInactive
	}

	// todo do not create session if user already has one
	rt, err := s.repo.CreateSession(u.ID)
	if err != nil {
		s.log.Error("failed to create session: %w", err)
		return nil, err
	}

	at, err := s.generateAccessToken(u)
	if err != nil {
		return nil, err
	}

	return &Tokens{
		AccessToken:  at,
		RefreshToken: rt,
	}, nil
}

func (s *Service) generateAccessToken(u *users.UserByIDResponse) (string, error) {
	claims := Token{
		ID:   u.ID,
		Name: fmt.Sprintf("%s %s", u.FirstName, u.LastName),
		Role: u.Role,
		StandardClaims: jwt.StandardClaims{
			Audience:  "dynapp",
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			Issuer:    "auth.dynapp",
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claims)

	return token.SignedString([]byte(s.jwtSecret))
}
