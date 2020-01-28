package requests

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog"
)

const (
	requestTaxi     = "taxi"
	requestDelivery = "delivery"
	requestGuest    = "guest"

	requestStatusNew = "new"
)

type Service interface {
	Create(ctx context.Context, r *createRequest) (*createResponse, error)
}

type service struct {
	db *gorm.DB
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

type request struct {
	ID          uint   `json:"id"`
	Type        string `json:"type"`
	UserID      uint   `json:"user_id"`
	Time        int64  `json:"time"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

func New(log *zerolog.Logger, db *gorm.DB) Service {
	s := service{db: db}
	srv := newLoggingMiddleware(log, &s)
	return srv
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
