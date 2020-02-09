package repository

import (
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"

	"github.com/ivch/dynasty/models"
)

type Auth struct {
	db *gorm.DB
}

func NewAuth(db *gorm.DB) *Auth {
	return &Auth{db: db}
}

func (a *Auth) CreateSession(userID uint) (string, error) {
	rt := uuid.NewV4()
	sess := models.Session{
		UserID:       userID,
		RefreshToken: rt,
		ExpiresIn:    time.Now().Add(100 * 365 * 24 * time.Hour).Unix(),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := a.db.Create(&sess).Error; err != nil {
		return "", err
	}

	return rt.String(), nil
}

func (a *Auth) FindSessionByAccessToken(token string) (*models.Session, error) {
	var sess models.Session
	if err := a.db.Where("refresh_token = ?", token).First(&sess).Error; err != nil {
		return nil, err
	}
	return &sess, nil
}

func (a *Auth) DeleteSessionByID(id string) error {
	return a.db.Where("id = ?", id).Delete(models.Session{}).Error
}
