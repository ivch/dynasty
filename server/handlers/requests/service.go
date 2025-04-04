package requests

import (
	"context"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/ivch/dynasty/common/errs"
	"github.com/ivch/dynasty/common/logger"
)

type RequestType int

const (
	_ RequestType = iota
	Guest
	Taxi
	Delivery
	Cargo

	defaultRequestStatus = "new"
	allowedFileType      = "image/jpeg"
	requestsPerDay       = 20
	filesPerRequest      = 3
	imgPathPrefix        = "req/i/"
	thumbPathPrefix      = "req/t/"
	defaultS3ACL         = "public-read"
)

var (
	oldRequestTypes = map[string]RequestType{
		"guest":    Guest,
		"taxi":     Taxi,
		"delivery": Delivery,
		"cargo":    Cargo,
	}
	newRequestTypes = map[RequestType]map[string]string{
		Guest: {
			"key": "guest",
			"en":  "Guest",
			"ru":  "Гость",
			"ua":  "Гість",
		},
		Taxi: {
			"key": "taxi",
			"en":  "Taxi",
			"ru":  "Такси",
			"ua":  "Таксі",
		},
		Delivery: {
			"key": "delivery",
			"en":  "Delivery",
			"ru":  "Доставка",
			"ua":  "Доставка",
		},
		Cargo: {
			"key": "cargo",
			"en":  "37-b Unload Area",
			"ru":  "37-Б Разгрузка",
			"ua":  "37-Б Розвантаження",
		},
	}
)

type requestsRepository interface {
	Create(req *Request) error
	GetRequestByIDAndUser(id, userID uint) (*Request, error)
	Update(update *UpdateRequest) error
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

func (s *Service) Update(_ context.Context, r *UpdateRequest) error {
	_, err := s.repo.GetRequestByIDAndUser(r.ID, r.UserID)
	if err != nil {
		s.log.Error("error finding request: %w", err)
		return err
	}

	// backward compatibility
	if r.Type != nil {
		if _, ok := oldRequestTypes[*r.Type]; ok {
			newType := oldRequestTypes[*r.Type]
			r.Rtype = &newType
		}
	}

	if r.Rtype != nil {
		if _, ok := newRequestTypes[*r.Rtype]; ok {
			oldType := newRequestTypes[*r.Rtype]["key"]
			r.Type = &oldType
		}
	}
	// end backward compatibility

	return s.repo.Update(r)
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
	dateFrom := time.Now().Add(-24 * time.Hour)
	list, err := s.repo.ListByUser(&RequestListFilter{
		DateFrom: &dateFrom,
		Offset:   0,
		Limit:    25,
		UserID:   r.UserID,
	})

	// backward compatibility
	if r.Type != "" {
		if _, ok := oldRequestTypes[r.Type]; ok {
			r.Rtype = oldRequestTypes[r.Type]
		}
	}

	if r.Rtype != 0 {
		if _, ok := newRequestTypes[r.Rtype]; ok {
			r.Type = newRequestTypes[r.Rtype]["key"]
		}
	}
	// end backward compatibility

	if err != nil {
		return nil, err
	}

	if len(list) >= requestsPerDay {
		return nil, errs.RequestPerDayLimitExceeded
	}

	r.Status = defaultRequestStatus

	if err := s.repo.Create(r); err != nil {
		s.log.Error("error creating request: %w", err)
		return nil, errors.New("failed to create request")
	}

	return r, nil
}

func GetRequestTypes() map[RequestType]map[string]string {
	return newRequestTypes
}
