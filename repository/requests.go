package repository

import (
	"github.com/jinzhu/gorm"

	"github.com/ivch/dynasty/models/entities"
)

type Requests struct {
	db *gorm.DB
}

func NewRequests(db *gorm.DB) *Requests {
	return &Requests{db: db}
}

func (r *Requests) Delete(id, userID uint) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&entities.Request{}).Error
}

func (r *Requests) GetRequestByIDAndUser(id, userId uint) (*entities.Request, error) {
	var req entities.Request

	if err := r.db.Where("id = ? AND user_id = ?", id, userId).First(&req).Error; err != nil {
		return nil, err
	}

	return &req, nil
}

func (r *Requests) Update(req *entities.Request) error {
	return r.db.Table(entities.Request{}.TableName()).Save(req).Error
}

func (r *Requests) ListByUser(userID, limit, offset uint) ([]*entities.Request, error) {
	var reqs []*entities.Request
	if err := r.db.Limit(limit).Offset(offset).Where("user_id = ?", userID).Find(&reqs).Error; err != nil {
		return nil, err
	}
	return reqs, nil
}

func (r *Requests) Create(req *entities.Request) (uint, error) {
	if err := r.db.Create(&req).Error; err != nil {
		return 0, err
	}
	return req.ID, nil
}
