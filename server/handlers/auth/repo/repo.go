package repository

import (
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"

	"github.com/ivch/dynasty/server/handlers/auth"
)

type Auth struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Auth {
	return &Auth{db: db}
}

func (a *Auth) CreateSession(userID uint) (string, error) {
	rt := uuid.NewV4()
	sess := auth.Session{
		UserID:       userID,
		RefreshToken: rt,
		ExpiresIn:    time.Now().Add(30 * 24 * time.Hour).Unix(), // 30 days
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := a.db.Create(&sess).Error; err != nil {
		return "", err
	}

	return rt.String(), nil
}

func (a *Auth) FindSessionByAccessToken(token string) (*auth.Session, error) {
	var sess auth.Session
	if err := a.db.Where("refresh_token = ?", token).First(&sess).Error; err != nil {
		return nil, err
	}
	return &sess, nil
}

func (a *Auth) DeleteSessionByID(id string) error {
	return a.db.Where("id = ?", id).Delete(auth.Session{}).Error
}

func (a *Auth) DeleteSessionByUserID(id uint) error {
	return a.db.Where("user_id = ?", id).Delete(auth.Session{}).Error
}
