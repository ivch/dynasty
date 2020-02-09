package repository

import (
	"github.com/jinzhu/gorm"

	"github.com/ivch/dynasty/models"
)

type Requests struct {
	db *gorm.DB
}

func NewRequests(db *gorm.DB) *Requests {
	return &Requests{db: db}
}

func (r *Requests) Delete(id, userID uint) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Request{}).Error
}

func (r *Requests) GetRequestByIDAndUser(id, userId uint) (*models.Request, error) {
	var req models.Request

	if err := r.db.Where("id = ? AND user_id = ?", id, userId).First(&req).Error; err != nil {
		return nil, err
	}

	return &req, nil
}

func (r *Requests) Update(req *models.Request) error {
	return r.db.Table(models.Request{}.TableName()).Save(req).Error
}

func (r *Requests) ListByUser(userID, limit, offset uint) ([]*models.Request, error) {
	var reqs []*models.Request
	if err := r.db.Limit(limit).Offset(offset).Where("user_id = ?", userID).Find(&reqs).Error; err != nil {
		return nil, err
	}
	return reqs, nil
}

func (r *Requests) Create(req *models.Request) (uint, error) {
	if err := r.db.Create(&req).Error; err != nil {
		return 0, err
	}
	return req.ID, nil
}
