package requests

import (
	"context"

	"github.com/rs/zerolog"

	"github.com/ivch/dynasty/models/dto"
	"github.com/ivch/dynasty/models/entities"
)

const (
	requestStatusNew = "new"
)

type Service interface {
	Create(ctx context.Context, r *dto.RequestCreateRequest) (*dto.RequestCreateResponse, error)
	Update(ctx context.Context, r *dto.RequestUpdateRequest) error
	Delete(ctx context.Context, r *dto.RequestByID) error
	Get(ctx context.Context, r *dto.RequestByID) (*dto.RequestByIDResponse, error)
	My(ctx context.Context, r *dto.RequestMyRequest) (*dto.RequestMyResponse, error)
}

type requestsRepository interface {
	Create(req *entities.Request) (uint, error)
	GetRequestByIDAndUser(id, userId uint) (*entities.Request, error)
	Update(req *entities.Request) error
	Delete(id, userID uint) error
	ListByUser(userID, limit, offset uint) ([]*entities.Request, error)
}

type service struct {
	repo requestsRepository
}

func (s *service) Get(_ context.Context, r *dto.RequestByID) (*dto.RequestByIDResponse, error) {
	req, err := s.repo.GetRequestByIDAndUser(r.ID, r.UserID)
	if err != nil {
		return nil, err
	}
	return &dto.RequestByIDResponse{Data: req}, nil
}

func (s *service) Delete(_ context.Context, r *dto.RequestByID) error {
	return s.repo.Delete(r.ID, r.UserID)
}

func (s *service) Update(_ context.Context, r *dto.RequestUpdateRequest) error {
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

func (s *service) My(_ context.Context, r *dto.RequestMyRequest) (*dto.RequestMyResponse, error) {
	reqs, err := s.repo.ListByUser(r.UserID, r.Limit, r.Offset)
	if err != nil {
		return nil, err
	}
	return &dto.RequestMyResponse{Data: reqs}, nil
}

func (s *service) Create(_ context.Context, r *dto.RequestCreateRequest) (*dto.RequestCreateResponse, error) {
	req := entities.Request{
		Type:        r.Type,
		UserID:      r.UserID,
		Time:        r.Time,
		Description: r.Description,
		Status:      requestStatusNew,
	}

	id, err := s.repo.Create(&req)

	return &dto.RequestCreateResponse{ID: id}, err
}

func newService(log *zerolog.Logger, repo requestsRepository) Service {
	s := service{repo: repo}
	srv := newLoggingMiddleware(log, &s)
	return srv
}
