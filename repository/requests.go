package repository

import (
	"time"

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

func (r *Requests) ListForGuard(limit, offset uint, status string) ([]*entities.Request, error) {
	from := time.Now().Add(-12 * time.Hour).Unix()
	to := time.Now().Add(24 * time.Hour).Unix()

	q := r.db.Preload("User").
		Preload("User.Building").
		Limit(limit).Offset(offset).
		Where("time >= ? AND time <= ?", from, to)
	if status != "all" {
		q = q.Where("status = ?", status)
	}

	var reqs []*entities.Request
	if err := q.Find(&reqs).Error; err != nil {
		return nil, err
	}
	return reqs, nil
}

func (r *Requests) UpdateForGuard(id uint, status string) error {
	// todo save who modified the request
	return r.db.Model(&entities.Request{}).Where("id = ?", id).Update("status", status).Error
}
