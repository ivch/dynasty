package repository

import (
	"time"

	"github.com/jinzhu/gorm"

	"github.com/ivch/dynasty/models/dto"
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

func (r *Requests) AddImage(userId, requestId uint, filename string) error {
	return r.db.Table(entities.Request{}.TableName()).
		Where("id = ? and user_id = ?", requestId, userId).
		Update("images", gorm.Expr("array_append(images, ?)", filename)).Error
}

func (r *Requests) DeleteImage(userId, requestId uint, filename string) error {
	return r.db.Table(entities.Request{}.TableName()).
		Where("id = ? and user_id = ?", requestId, userId).
		Update("images", gorm.Expr("array_remove(images, ?)", filename)).Error
}

func (r *Requests) Create(req *entities.Request) (uint, error) {
	if err := r.db.Create(&req).Error; err != nil {
		return 0, err
	}
	return req.ID, nil
}

func (r *Requests) ListForGuard(req *dto.RequestListFilterRequest) ([]*entities.Request, error) {
	q := buildGuardFilterQuery(r.db, req).Limit(req.Limit).Offset(req.Offset).Order("time desc")

	var reqs []*entities.Request
	if err := q.Find(&reqs).Error; err != nil {
		return nil, err
	}
	return reqs, nil
}

func (r *Requests) CountForGuard(req *dto.RequestListFilterRequest) (int, error) {
	q := buildGuardFilterQuery(r.db, req)

	var count int
	if err := q.Model(&entities.Request{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *Requests) UpdateForGuard(id uint, status string) error {
	// todo save who modified the request
	return r.db.Model(&entities.Request{}).Where("id = ?", id).Update("status", status).Error
}

func (r *Requests) ListByUser(req *dto.RequestListFilterRequest) ([]*entities.Request, error) {
	var reqs []*entities.Request
	if err := r.db.Order("time desc").Limit(req.Limit).Offset(req.Offset).Where("user_id = ?", req.UserID).
		Find(&reqs).Error; err != nil {
		return nil, err
	}
	return reqs, nil
}

func buildGuardFilterQuery(db *gorm.DB, req *dto.RequestListFilterRequest) *gorm.DB {
	from := time.Now().Add(-12 * time.Hour).Unix()
	to := time.Now().Add(12 * time.Hour).Unix()

	q := db.Preload("User.Building").Preload("User.Entry").
		Where("time >= ? AND time <= ?", from, to)

	if req.Type != "all" {
		q = q.Where("type = ?", req.Type)
	}

	if req.Place == "kpp" {
		q = q.Where("type IN (?)", []string{"taxi", "guest", "delivery"})
	}

	if req.Status != "all" {
		q = q.Where("status = ?", req.Status)
	}

	if req.Apartment != "" {
		q = q.Joins("left join users on users.id = requests.user_id").
			Where("users.apartment = ?", req.Apartment)
	}

	return q
}
