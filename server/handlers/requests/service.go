package requests

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/ivch/dynasty/common/logger"
)

const (
	defaultRequestStatus = "new"
	allowedFileType      = "image/jpeg"
	filesPerRequest      = 3
	imgPathPrefix        = "req/i/"
	thumbPathPrefix      = "req/t/"
	defaultS3ACL         = "public-read"
)

type requestsRepository interface {
	Create(req *Request) error
	GetRequestByIDAndUser(id, userId uint) (*Request, error)
	Update(req *Request) error
	Delete(id, userID uint) error
	ListByUser(r *RequestListFilter) ([]*Request, error)
	ListForGuard(req *RequestListFilter) ([]*Request, error)
	UpdateForGuard(id uint, status string) error
	CountForGuard(req *RequestListFilter) (int, error)
	AddImage(userID, requestID uint, filename string) error
	DeleteImage(userID, requestID uint, filename string) error
}

type s3Client interface {
	PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error)
	DeleteObject(input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error)
}

type Service struct {
	repo     requestsRepository
	s3Client s3Client
	s3Space  string
	cdnHost  string
	log      logger.Logger
}

func New(log logger.Logger, repo requestsRepository, s3Client s3Client, s3Space, cdnHost string) *Service {
	s := Service{repo: repo, s3Space: s3Space, s3Client: s3Client, cdnHost: cdnHost, log: log}

	return &s
}

func (s *Service) Get(_ context.Context, r *Request) (*Request, error) {
	req, err := s.repo.GetRequestByIDAndUser(r.ID, r.UserID)
	if err != nil {
		s.log.Error("error finding request: %w", err)
		return nil, err
	}
	// todo: make separate requests for user data
	req.ImagesURL = make([]map[string]string, len(req.Images))

	for i := range req.Images {
		req.ImagesURL[i] = s.buildImageURL(req.Images[i])
	}

	return req, nil
}

func (s *Service) Delete(_ context.Context, r *Request) error {
	req, err := s.repo.GetRequestByIDAndUser(r.ID, r.UserID)
	if err != nil {
		s.log.Error("failed to delete request %d: %w", r.ID, err)
		return err
	}

	for i := range req.Images {
		if err := s.deleteImageFromS3(req.Images[i]); err != nil {
			s.log.Error("error deleting image for request %d: %w", r.ID, err)
		}
	}

	return s.repo.Delete(r.ID, r.UserID)
}

func (s *Service) Update(_ context.Context, r *Request) error {
	req, err := s.repo.GetRequestByIDAndUser(r.ID, r.UserID)
	if err != nil {
		s.log.Error("error finding request: %w", err)
		return err
	}

	if err := s.repo.Update(req); err != nil {
		s.log.Error("error updating request: %w", err)
		return err
	}

	return nil
}

func (s *Service) My(_ context.Context, r *RequestListFilter) ([]*Request, error) {
	reqs, err := s.repo.ListByUser(r)
	if err != nil {
		return nil, err
	}

	for i := range reqs {
		reqs[i].ImagesURL = make([]map[string]string, len(reqs[i].Images))
		for j := range reqs[i].Images {
			reqs[i].ImagesURL[j] = s.buildImageURL(reqs[i].Images[j])
		}
	}

	return reqs, nil
}

func (s *Service) Create(_ context.Context, r *Request) (*Request, error) {
	r.Status = defaultRequestStatus

	if err := s.repo.Create(r); err != nil {
		s.log.Error("error creating request: %w", err)
		return nil, errors.New("failed to create request")
	}

	return r, nil
}
