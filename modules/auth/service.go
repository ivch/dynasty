package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/rs/zerolog"

	"github.com/ivch/dynasty/models/dto"
	"github.com/ivch/dynasty/models/entities"
)

type Service interface {
	Login(ctx context.Context, request *dto.AuthLoginRequest) (*dto.AuthLoginResponse, error)
	Refresh(ctx context.Context, r *dto.AuthRefreshTokenRequest) (*dto.AuthLoginResponse, error)
	Gwfa(token string) (uint, error)
}

type userService interface {
	UserByPhoneAndPassword(ctx context.Context, phone, password string) (*entities.User, error)
	UserByID(ctx context.Context, id uint) (*entities.User, error)
}

type authRepository interface {
	CreateSession(userID uint) (string, error)
	FindSessionByAccessToken(token string) (*entities.Session, error)
	DeleteSessionByID(id string) error
}

type service struct {
	log       *zerolog.Logger
	uSrv      userService
	repo      authRepository
	jwtSecret string
}

func newService(log *zerolog.Logger, repo authRepository, uSrv userService, jwtSecret string) Service {
	s := &service{
		log:       log,
		repo:      repo,
		uSrv:      uSrv,
		jwtSecret: jwtSecret,
	}
	svc := newLoggingMiddleware(log, s)

	return svc
}

func (s *service) Gwfa(token string) (uint, error) {
	var myClaims entities.Token
	t, err := jwt.ParseWithClaims(token, &myClaims, func(token *jwt.Token) (i interface{}, e error) {
		return []byte(s.jwtSecret), nil
	})
	if err != nil {
		return 0, fmt.Errorf("%s: %w", entities.ErrParsingToken, err)
	}

	claims, ok := t.Claims.(*entities.Token)
	if !ok {
		return 0, entities.ErrParsingTokenClaims
	}

	// if !t.Valid {
	// 	return 0, models.ErrTokenIsInvalid
	// }
	//
	// if time.Now().After(time.Unix(claims.ExpiresAt, 0)) {
	// 	return 0, models.ErrTokenExpired
	// }

	return claims.ID, nil
}

func (s *service) Refresh(ctx context.Context, r *dto.AuthRefreshTokenRequest) (*dto.AuthLoginResponse, error) {
	sess, err := s.repo.FindSessionByAccessToken(r.Token)
	if err != nil {
		return nil, err
	}

	if err := s.repo.DeleteSessionByID(sess.ID); err != nil {
		return nil, err
	}

	if time.Now().Unix() > sess.ExpiresIn {
		return nil, errors.New("token expired")
	}

	u, err := s.uSrv.UserByID(ctx, sess.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	rt, err := s.repo.CreateSession(sess.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	at, err := s.generateAccessToken(u)
	if err != nil {
		return nil, fmt.Errorf("failed to create access token: %w", err)
	}

	return &dto.AuthLoginResponse{
		AccessToken:  at,
		RefreshToken: rt,
	}, nil
}

func (s *service) Login(ctx context.Context, r *dto.AuthLoginRequest) (*dto.AuthLoginResponse, error) {
	u, err := s.uSrv.UserByPhoneAndPassword(ctx, r.Phone, r.Password)
	if err != nil {
		s.log.Error().Err(err).Msg("failed to find session")
		return nil, entities.ErrSessionNotFound
	}

	if !u.Active {
		return nil, errUserIsInactive
	}

	// todo do not create session if user already has one
	rt, err := s.repo.CreateSession(u.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	at, err := s.generateAccessToken(u)
	if err != nil {
		return nil, err
	}

	return &dto.AuthLoginResponse{
		AccessToken:  at,
		RefreshToken: rt,
	}, nil
}

func (s *service) generateAccessToken(u *entities.User) (string, error) {
	claims := entities.Token{
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
