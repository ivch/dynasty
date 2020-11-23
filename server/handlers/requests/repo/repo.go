package repository

import (
	"time"

	"github.com/jinzhu/gorm"

	"github.com/ivch/dynasty/server/handlers/requests"
)

type Requests struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Requests {
	return &Requests{db: db}
}

func (r *Requests) Delete(id, userID uint) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&requests.Request{}).Error
}

func (r *Requests) GetRequestByIDAndUser(id, userId uint) (*requests.Request, error) {
	var req requests.Request
	if err := r.db.Where("id = ? AND user_id = ?", id, userId).First(&req).Error; err != nil {
		return nil, err
	}
	return &req, nil
}

func (r *Requests) Update(req *requests.Request) error {
	return r.db.Table(requests.Request{}.TableName()).Save(req).Error
}

func (r *Requests) Create(req *requests.Request) error {
	return r.db.Create(req).Error
}

func (r *Requests) ListForGuard(req *requests.RequestListFilter) ([]*requests.Request, error) {
	q := buildGuardFilterQuery(r.db, req).Limit(req.Limit).Offset(req.Offset).Order("time desc")

	var reqs []*requests.Request
	if err := q.Find(&reqs).Error; err != nil {
		return nil, err
	}
	return reqs, nil
}

func (r *Requests) CountForGuard(req *requests.RequestListFilter) (int, error) {
	q := buildGuardFilterQuery(r.db, req)

	var count int
	if err := q.Model(&requests.Request{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *Requests) UpdateForGuard(id uint, status string) error {
	// todo save who modified the request
	return r.db.Model(&requests.Request{}).Where("id = ?", id).Update("status", status).Error
}

func (r *Requests) ListByUser(req *requests.RequestListFilter) ([]*requests.Request, error) {
	var reqs []*requests.Request
	// to get soft deleted records db.Unscoped().Where().Find()
	if err := r.db.Order("time desc").Limit(req.Limit).Offset(req.Offset).Where("user_id = ?", req.UserID).
		Find(&reqs).Error; err != nil {
		return nil, err
	}
	return reqs, nil
}

func buildGuardFilterQuery(db *gorm.DB, req *requests.RequestListFilter) *gorm.DB {
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

func (r *Requests) AddImage(userId, requestId uint, filename string) error {
	return r.db.Table(requests.Request{}.TableName()).
		Where("id = ? and user_id = ?", requestId, userId).
		Update("images", gorm.Expr("array_append(images, ?)", filename)).Error
}

func (r *Requests) DeleteImage(userId, requestId uint, filename string) error {
	return r.db.Table(requests.Request{}.TableName()).
		Where("id = ? and user_id = ?", requestId, userId).
		Update("images", gorm.Expr("array_remove(images, ?)", filename)).Error
}
