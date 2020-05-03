package requests

import (
	"context"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/rs/zerolog"

	"github.com/ivch/dynasty/models/dto"
	"github.com/ivch/dynasty/models/entities"
)

type Service interface {
	Create(ctx context.Context, r *dto.RequestCreateRequest) (*dto.RequestCreateResponse, error)
	Update(ctx context.Context, r *dto.RequestUpdateRequest) error
	Delete(ctx context.Context, r *dto.RequestByID) error
	Get(ctx context.Context, r *dto.RequestByID) (*dto.RequestByIDResponse, error)
	My(ctx context.Context, r *dto.RequestListFilterRequest) (*dto.RequestMyResponse, error)

	UploadImage(ctx context.Context, r *dto.UploadImageRequest) (*dto.UploadImageResponse, error)
	DeleteImage(ctx context.Context, r *dto.DeleteImageRequest) error

	GuardRequestList(ctx context.Context, r *dto.RequestListFilterRequest) (*dto.RequestGuardListResponse, error)
	GuardUpdateRequest(ctx context.Context, r *dto.GuardUpdateRequest) error
}

type requestsRepository interface {
	Create(req *entities.Request) (uint, error)
	GetRequestByIDAndUser(id, userId uint) (*entities.Request, error)
	Update(req *entities.Request) error
	Delete(id, userID uint) error
	ListByUser(r *dto.RequestListFilterRequest) ([]*entities.Request, error)
	ListForGuard(req *dto.RequestListFilterRequest) ([]*entities.Request, error)
	UpdateForGuard(id uint, status string) error
	CountForGuard(req *dto.RequestListFilterRequest) (int, error)
	AddImage(userID, requestID uint, filename string) error
	DeleteImage(userID, requestID uint, filename string) error
}

type s3Client interface {
	PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error)
	DeleteObject(input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error)
}

type service struct {
	repo     requestsRepository
	s3Client s3Client
	s3Space  string
	cdnHost  string
}

func (s *service) Get(_ context.Context, r *dto.RequestByID) (*dto.RequestByIDResponse, error) {
	req, err := s.repo.GetRequestByIDAndUser(r.ID, r.UserID)
	if err != nil {
		return nil, err
	}

	resp := dto.RequestByIDResponse{
		ID:          req.ID,
		UserID:      req.UserID,
		Type:        req.Type,
		Time:        req.Time,
		Description: req.Description,
		Status:      req.Status,
		Images:      make([]map[string]string, len(req.Images)),
	}

	for i := range req.Images {
		resp.Images[i] = s.buildImageURL(req.Images[i])
	}
	return &resp, nil
}

func (s *service) Delete(_ context.Context, r *dto.RequestByID) error {
	req, err := s.repo.GetRequestByIDAndUser(r.ID, r.UserID)
	if err != nil {
		return err
	}

	if len(req.Images) == 0 {
		return s.repo.Delete(r.ID, r.UserID)
	}

	for i := range req.Images {
		// todo rework this code to handle errors correctly
		// nolint: errcheck
		s.deleteImageFromS3(req.Images[i])
	}

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

func (s *service) My(_ context.Context, r *dto.RequestListFilterRequest) (*dto.RequestMyResponse, error) {
	reqs, err := s.repo.ListByUser(r)
	if err != nil {
		return nil, err
	}

	data := make([]*dto.RequestByIDResponse, len(reqs))
	for i := range reqs {
		data[i] = &dto.RequestByIDResponse{
			ID:          reqs[i].ID,
			Type:        reqs[i].Type,
			UserID:      reqs[i].UserID,
			Time:        reqs[i].Time,
			Description: reqs[i].Description,
			Status:      reqs[i].Status,
			Images:      make([]map[string]string, len(reqs[i].Images)),
		}
		for j := range reqs[i].Images {
			data[i].Images[j] = s.buildImageURL(reqs[i].Images[j])
		}
	}

	return &dto.RequestMyResponse{Data: data}, nil
}

func (s *service) Create(_ context.Context, r *dto.RequestCreateRequest) (*dto.RequestCreateResponse, error) {
	req := entities.Request{
		Type:        r.Type,
		UserID:      r.UserID,
		Time:        r.Time,
		Description: r.Description,
		Status:      "new",
	}

	id, err := s.repo.Create(&req)

	return &dto.RequestCreateResponse{ID: id}, err
}

func newService(log *zerolog.Logger, repo requestsRepository, s3Client s3Client, s3Space, cdnHost string) Service {
	s := service{repo: repo, s3Space: s3Space, s3Client: s3Client, cdnHost: cdnHost}
	srv := newLoggingMiddleware(log, &s)
	return srv
}
