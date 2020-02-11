package requests

import (
	"context"

	"github.com/rs/zerolog"

	"github.com/ivch/dynasty/models"
)

const (
	requestStatusNew = "new"
)

type Service interface {
	Create(ctx context.Context, r *createRequest) (*createResponse, error)
	Update(ctx context.Context, r *updateRequest) error
	Delete(ctx context.Context, r *byIDRequest) error
	Get(ctx context.Context, r *byIDRequest) (*getResponse, error)
	My(ctx context.Context, r *myRequest) (*myResponse, error)
}

type requestsRepository interface {
	Create(req *models.Request) (uint, error)
	GetRequestByIDAndUser(id, userId uint) (*models.Request, error)
	Update(req *models.Request) error
	Delete(id, userID uint) error
	ListByUser(userID, limit, offset uint) ([]*models.Request, error)
}

type service struct {
	repo requestsRepository
}

type byIDRequest struct {
	UserID uint `validate:"required"`
	ID     uint `validate:"required"`
}

type getResponse struct {
	*models.Request
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

type myResponse struct {
	Data []*models.Request `json:"data"`
}

type updateRequest struct {
	ID          uint
	UserID      uint    `gorm:"-"`
	Type        *string `json:"type,omitempty"`
	Time        *int64  `json:"time,omitempty"`
	Description *string `json:"description,omitempty"`
	Status      *string `json:"status,omitempty"`
}

func (s *service) Get(_ context.Context, r *byIDRequest) (*getResponse, error) {
	req, err := s.repo.GetRequestByIDAndUser(r.ID, r.UserID)
	if err != nil {
		return nil, err
	}
	return &getResponse{req}, nil
}

func (s *service) Delete(_ context.Context, r *byIDRequest) error {
	return s.repo.Delete(r.ID, r.UserID)
}

func (s *service) Update(_ context.Context, r *updateRequest) error {
	req, err := s.repo.GetRequestByIDAndUser(r.ID, r.UserID)
	if err != nil {
		return err
	}

	if r.Type != nil {
		req.Type = *r.Type
	}

	if r.Description != nil {
		req.Description = *r.Description
	}

	if r.Status != nil {
		req.Status = *r.Status
	}

	if r.Time != nil {
		req.Time = *r.Time
	}

	return s.repo.Update(req)
}

func (s *service) My(_ context.Context, r *myRequest) (*myResponse, error) {
	reqs, err := s.repo.ListByUser(r.UserID, r.Limit, r.Offset)
	if err != nil {
		return nil, err
	}
	return &myResponse{Data: reqs}, nil
}

func (s *service) Create(_ context.Context, r *createRequest) (*createResponse, error) {
	req := models.Request{
		Type:        r.Type,
		UserID:      r.UserID,
		Time:        r.Time,
		Description: r.Description,
		Status:      requestStatusNew,
	}

	id, err := s.repo.Create(&req)

	return &createResponse{ID: id}, err
}

func newService(log *zerolog.Logger, repo requestsRepository) Service {
	s := service{repo: repo}
	srv := newLoggingMiddleware(log, &s)
	return srv
}
