package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog"
	uuid "github.com/satori/go.uuid"

	userClient "github.com/dynastiateam/backend/bkp/auth/client/user"
)

type Service interface {
	Login(ctx context.Context, request *loginRequest) (*loginResponse, error)
	Refresh(ctx context.Context, r *refreshTokenRequest) (*loginResponse, error)
	Gwfa(token string) (*uint, error)
}

type UserService interface {
	UserByPhoneAndPassword(ctx context.Context, phone, password string) (*userClient.User, error)
	UserByID(ctx context.Context, id uint) (*userClient.User, error)
}

type service struct {
	log       *zerolog.Logger
	uSrv      UserService
	db        *gorm.DB
	jwtSecret string
}

type session struct {
	ID           string
	UserID       uint
	RefreshToken uuid.UUID
	ExpiresIn    int64
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type loginRequest struct {
	Phone    string `json:"phone" validate:"required"`
	Password string `json:"password" validate:"required"`
	// IP       net.IP `json:"ip" validate:"required"`
	// Ua       string `json:"ua" validate:"required"`
}

type loginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type refreshTokenRequest struct {
	Token string `validate:"required"`
}

type Token struct {
	ID   uint
	Name string
	Role int
	jwt.StandardClaims
}

func NewService(log *zerolog.Logger, db *gorm.DB, uSrv UserService, jwtSecret string) Service {
	s := &service{
		log:       log,
		uSrv:      uSrv,
		db:        db,
		jwtSecret: jwtSecret,
	}
	svc := newLoggingMiddleware(log, s)

	return svc
}

var (
	errParsingToken       = errors.New("failed to parse token")
	errParsingTokenClaims = errors.New("failed to parse token claims")
	errTokenIsInvalid     = errors.New("token is invalid")
	errTokenExpired       = errors.New("token expired")
)

func (s *service) Gwfa(token string) (*uint, error) {
	var myClaims Token
	t, err := jwt.ParseWithClaims(token, &myClaims, func(token *jwt.Token) (i interface{}, e error) {
		return []byte(s.jwtSecret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errParsingToken, err)
	}

	claims, ok := t.Claims.(*Token)
	if !ok {
		return nil, errParsingTokenClaims
	}

	if !t.Valid {
		return nil, errTokenIsInvalid
	}

	if time.Now().After(time.Unix(claims.ExpiresAt, 0)) {
		return nil, errTokenExpired
	}

	return &claims.ID, nil
}

func (s *service) Refresh(ctx context.Context, r *refreshTokenRequest) (*loginResponse, error) {
	var sess session
	if err := s.db.Where("refresh_token = ?", r.Token).First(&sess).Error; err != nil {
		return nil, err
	}

	if err := s.db.Where("id = ?", sess.ID).Delete(session{}).Error; err != nil {
		return nil, err
	}

	if time.Now().Unix() > sess.ExpiresIn {
		return nil, errors.New("token expired")
	}

	u, err := s.uSrv.UserByID(ctx, sess.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	rt, err := s.createSession(sess.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	at, err := s.generateAccessToken(u)
	if err != nil {
		return nil, fmt.Errorf("failed to create access token: %w", err)
	}

	return &loginResponse{
		AccessToken:  at,
		RefreshToken: rt,
	}, nil
}

func (s *service) Login(ctx context.Context, r *loginRequest) (*loginResponse, error) {
	u, err := s.uSrv.UserByPhoneAndPassword(ctx, r.Phone, r.Password)
	if err != nil {
		return nil, err
	}

	// todo do not create session if user already has one
	rt, err := s.createSession(u.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	at, err := s.generateAccessToken(u)
	if err != nil {
		return nil, err
	}

	return &loginResponse{
		AccessToken:  at,
		RefreshToken: rt,
	}, nil
}

func (s *service) generateAccessToken(u *userClient.User) (string, error) {
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

func (s *service) createSession(userID uint) (string, error) {
	rt := uuid.NewV4()
	sess := session{
		UserID:       userID,
		RefreshToken: rt,
		ExpiresIn:    time.Now().Add(100 * 365 * 24 * time.Hour).Unix(),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.db.Create(&sess).Error; err != nil {
		return "", err
	}

	return rt.String(), nil
}
