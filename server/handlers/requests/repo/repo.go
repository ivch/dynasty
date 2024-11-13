package repository

import (
	"fmt"
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
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := r.updateRequestHistory(tx, id, &requests.HistoryRecord{
			Time:   time.Now(),
			UserID: userID,
			Action: "record deleted",
		}); err != nil {
			return err
		}
		return tx.Where("id = ? AND user_id = ?", id, userID).Delete(&requests.Request{}).Error
	})
}

func (r *Requests) GetRequestByIDAndUser(id, userID uint) (*requests.Request, error) {
	var req requests.Request
	if err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&req).Error; err != nil {
		return nil, err
	}
	return &req, nil
}

func (r *Requests) Update(req *requests.UpdateRequest) error {
	update := make(map[string]interface{})
	if req.Type != nil {
		update["type"] = *req.Type
	}
	if req.Rtype != nil {
		update["rtype"] = *req.Rtype
	}

	if req.Description != nil {
		update["description"] = *req.Description
	}

	if req.Status != nil {
		update["status"] = *req.Status
	}

	if req.Time != nil {
		update["time"] = *req.Time
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := r.updateRequestHistory(tx, req.ID, &requests.HistoryRecord{
			Time:   time.Now(),
			UserID: req.UserID,
			Action: fmt.Sprintf("record updated:%v", update),
		}); err != nil {
			return err
		}
		return tx.Table(requests.Request{}.TableName()).Where("id = ?", req.ID).Updates(update).Error
	})
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
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := r.updateRequestHistory(tx, id, &requests.HistoryRecord{
			Time:   time.Now(),
			UserID: 0,
			Action: fmt.Sprintf("guard changed status: %s", status),
		}); err != nil {
			return err
		}
		return tx.Model(&requests.Request{}).Where("id = ?", id).Update("status", status).Error
	})
}

func (r *Requests) ListByUser(req *requests.RequestListFilter) ([]*requests.Request, error) {
	var reqs []*requests.Request
	// to get soft deleted records db.Unscoped().Where().Find()
	q := r.db.Where("user_id = ?", req.UserID)
	if req.DateFrom != nil {
		q = q.Where("time >= ?", req.DateFrom.Unix())
	}

	if req.DateTo != nil {
		q = q.Where("time <= ? ", req.DateTo.Unix())
	}

	if err := q.Order("time desc").Limit(req.Limit).Offset(req.Offset).Find(&reqs).Error; err != nil {
		return nil, err
	}
	return reqs, nil
}

func buildGuardFilterQuery(db *gorm.DB, req *requests.RequestListFilter) *gorm.DB {
	from := time.Now().Add(-2 * 24 * time.Hour).Unix()
	to := time.Now().Add(2 * 24 * time.Hour).Unix()

	q := db.Preload("User.Building").Preload("User.Entry").
		Where("time >= ? AND time <= ?", from, to)

	if req.Type != "all" {
		q = q.Where("type = ?", req.Type)
	}

	if req.Place == "kpp" {
		q = q.Where("type IN (?)", []string{"taxi", "guest", "delivery", "cargo"})
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

func (r *Requests) AddImage(userID, requestID uint, filename string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := r.updateRequestHistory(tx, requestID, &requests.HistoryRecord{
			Time:   time.Now(),
			UserID: userID,
			Action: fmt.Sprintf("uploaded image: %s", filename),
		}); err != nil {
			return err
		}
		return tx.Table(requests.Request{}.TableName()).
			Where("id = ? and user_id = ?", requestID, userID).
			Update("images", gorm.Expr("array_append(images, ?)", filename)).Error
	})
}

func (r *Requests) DeleteImage(userID, requestID uint, filename string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := r.updateRequestHistory(tx, requestID, &requests.HistoryRecord{
			Time:   time.Now(),
			UserID: userID,
			Action: fmt.Sprintf("deleted image: %s", filename),
		}); err != nil {
			return err
		}
		return tx.Table(requests.Request{}.TableName()).
			Where("id = ? and user_id = ?", requestID, userID).
			Update("images", gorm.Expr("array_remove(images, ?)", filename)).Error
	})
}

func (r *Requests) updateRequestHistory(tx *gorm.DB, requestID uint, rec fmt.Stringer) error {
	return tx.Table(requests.Request{}.TableName()).
		Where("id = ?", requestID).
		Update("history", gorm.Expr("array_append(history, ?)", rec.String())).Error
}
