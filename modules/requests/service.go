package requests

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog"
)

const (
	// requestTaxi     = "taxi"
	// requestDelivery = "delivery"
	// requestGuest    = "guest"

	requestStatusNew = "new"
)

type Service interface {
	Create(ctx context.Context, r *createRequest) (*createResponse, error)
	Update(ctx context.Context, r *updateRequest) error
	Delete(ctx context.Context, r *byIDRequest) error
	Get(ctx context.Context, r *byIDRequest) (*getResponse, error)
	My(ctx context.Context, r *myRequest) (*myResponse, error)
}

type service struct {
	db *gorm.DB
}

type byIDRequest struct {
	UserID uint `validate:"required"`
	ID     uint `validate:"required"`
}

type getResponse struct {
	request
}

type createRequest struct {
	Type        string `json:"type" validate:"required"`
	Time        int64  `json:"time" validate:"required"`
	UserID      uint   `json:"user_id" validate:"required"`
	Description string `json:"description"`
}

type createResponse struct {
	ID uint `json:"id"`
}

type myRequest struct {
	UserID uint `json:"user_id"`
	Offset uint `json:"offset"`
	Limit  uint `json:"limit"`
}

type myResponse []request

type updateRequest struct {
	ID          uint
	UserID      uint    `gorm:"-"`
	Type        *string `json:"type,omitempty"`
	Time        *int64  `json:"time,omitempty"`
	Description *string `json:"description,omitempty"`
	Status      *string `json:"status,omitempty"`
}

type request struct {
	ID          uint   `json:"id"`
	Type        string `json:"type"`
	UserID      uint   `json:"user_id"`
	Time        int64  `json:"time"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

func (request) TableName() string { return "requests" }

func (s *service) Get(_ context.Context, r *byIDRequest) (*getResponse, error) {
	var req request

	if err := s.db.Where("id = ? AND user_id = ?", r.ID, r.UserID).First(&req).Error; err != nil {
		return nil, err
	}

	return &getResponse{req}, nil
}
func (s *service) Delete(_ context.Context, r *byIDRequest) error {
	var req request

	if err := s.db.Where("id = ? AND user_id = ?", r.ID, r.UserID).First(&req).Error; err != nil {
		return err
	}

	return s.db.Delete(req).Error
}

func (s *service) Update(_ context.Context, r *updateRequest) error {
	if err := s.db.Where("id = ? AND user_id = ?", r.ID, r.UserID).First(&request{}).Error; err != nil {
		return err
	}

	if err := s.db.Table(request{}.TableName()).Save(&r).Error; err != nil {
		return err
	}

	return nil
}

func (s *service) My(_ context.Context, r *myRequest) (*myResponse, error) {
	var res myResponse

	if err := s.db.Limit(r.Limit).Offset(r.Offset).Where("user_id = ?", r.UserID).Find(&res).Error; err != nil {
		return nil, err
	}

	return &res, nil
}

func (s *service) Create(_ context.Context, r *createRequest) (*createResponse, error) {
	req := request{
		Type:        r.Type,
		UserID:      r.UserID,
		Time:        r.Time,
		Description: r.Description,
		Status:      requestStatusNew,
	}

	if err := s.db.Create(&req).Error; err != nil {
		return nil, err
	}

	return &createResponse{ID: req.ID}, nil
}

func newService(log *zerolog.Logger, db *gorm.DB) Service {
	s := service{db: db}
	srv := newLoggingMiddleware(log, &s)
	return srv
}
