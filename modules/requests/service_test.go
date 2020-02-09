package requests

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/rs/zerolog"

	"github.com/ivch/dynasty/models"
)

var (
	defaultLogger *zerolog.Logger
	errTestError  = errors.New("some err")
)

func TestMain(m *testing.M) {
	logger := zerolog.New(ioutil.Discard)
	defaultLogger = &logger
	os.Exit(m.Run())
}

func TestService_Get(t *testing.T) {
	tests := []struct {
		name    string
		repo    requestsRepository
		req     *byIDRequest
		wantErr bool
		want    *getResponse
	}{
		{
			name: "error no request",
			repo: &requestsRepositoryMock{
				GetRequestByIDAndUserFunc: func(_ uint, _ uint) (*models.Request, error) {
					return nil, errTestError
				},
			},
			req: &byIDRequest{
				UserID: 1,
				ID:     1,
			},
			wantErr: true,
		},
		{
			name: "ok",
			repo: &requestsRepositoryMock{
				GetRequestByIDAndUserFunc: func(_ uint, _ uint) (*models.Request, error) {
					return &models.Request{
						ID:          1,
						Type:        "1",
						UserID:      1,
						Time:        1,
						Description: "1",
						Status:      "1",
					}, nil
				},
			},
			req: &byIDRequest{
				UserID: 1,
				ID:     1,
			},
			wantErr: false,
			want: &getResponse{
				Request: &models.Request{
					ID:          1,
					Type:        "1",
					UserID:      1,
					Time:        1,
					Description: "1",
					Status:      "1",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newService(defaultLogger, tt.repo)
			got, err := s.Get(context.Background(), tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestService_Update(t *testing.T) {
	tests := []struct {
		name    string
		repo    requestsRepository
		req     *updateRequest
		wantErr bool
	}{
		{
			name: "error no request",
			repo: &requestsRepositoryMock{
				GetRequestByIDAndUserFunc: func(_ uint, _ uint) (*models.Request, error) {
					return nil, errTestError
				},
			},
			req: &updateRequest{
				ID:     1,
				UserID: 1,
			},
			wantErr: true,
		},
		{
			name: "error type not updated",
			repo: &requestsRepositoryMock{
				GetRequestByIDAndUserFunc: func(_ uint, _ uint) (*models.Request, error) {
					return &models.Request{
						Type: "1",
					}, nil
				},
				UpdateFunc: func(req *models.Request) error {
					if req.Type != "2" {
						return errTestError
					}
					return nil
				},
			},
			req: &updateRequest{
				ID:     1,
				UserID: 1,
				Type:   nil,
			},
			wantErr: true,
		},
		{
			name: "error description not updated",
			repo: &requestsRepositoryMock{
				GetRequestByIDAndUserFunc: func(_ uint, _ uint) (*models.Request, error) {
					return &models.Request{
						Description: "1",
					}, nil
				},
				UpdateFunc: func(req *models.Request) error {
					if req.Description != "2" {
						return errTestError
					}
					return nil
				},
			},
			req: &updateRequest{
				ID:          1,
				UserID:      1,
				Description: nil,
			},
			wantErr: true,
		},
		{
			name: "error status not updated",
			repo: &requestsRepositoryMock{
				GetRequestByIDAndUserFunc: func(_ uint, _ uint) (*models.Request, error) {
					return &models.Request{
						Status: "1",
					}, nil
				},
				UpdateFunc: func(req *models.Request) error {
					if req.Status != "2" {
						return errTestError
					}
					return nil
				},
			},
			req: &updateRequest{
				ID:     1,
				UserID: 1,
				Status: nil,
			},
			wantErr: true,
		},
		{
			name: "error time not updated",
			repo: &requestsRepositoryMock{
				GetRequestByIDAndUserFunc: func(_ uint, _ uint) (*models.Request, error) {
					return &models.Request{
						Time: 1,
					}, nil
				},
				UpdateFunc: func(req *models.Request) error {
					if req.Time != 2 {
						return errTestError
					}
					return nil
				},
			},
			req: &updateRequest{
				ID:     1,
				UserID: 1,
				Time:   nil,
			},
			wantErr: true,
		},
		{
			name: "ok",
			repo: &requestsRepositoryMock{
				GetRequestByIDAndUserFunc: func(_ uint, _ uint) (*models.Request, error) {
					return &models.Request{
						Type: "1",
					}, nil
				},
				UpdateFunc: func(req *models.Request) error {
					return nil
				},
			},
			req: &updateRequest{
				ID:     1,
				UserID: 1,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newService(defaultLogger, tt.repo)
			err := s.Update(context.Background(), tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestService_My(t *testing.T) {
	tests := []struct {
		name    string
		repo    requestsRepository
		req     *myRequest
		wantErr bool
		want    *myResponse
	}{
		{
			name: "error from db",
			repo: &requestsRepositoryMock{
				ListByUserFunc: func(_ uint, _ uint, _ uint) ([]*models.Request, error) {
					return nil, errTestError
				},
			},
			req: &myRequest{
				UserID: 1,
				Offset: 0,
				Limit:  1,
			},
			wantErr: true,
		},
		{
			name: "ok",
			repo: &requestsRepositoryMock{
				ListByUserFunc: func(_ uint, _ uint, _ uint) ([]*models.Request, error) {
					return []*models.Request{
						{
							ID:          1,
							Type:        "1",
							UserID:      1,
							Time:        1,
							Description: "1",
							Status:      "1",
						},
					}, nil
				},
			},
			req: &myRequest{
				UserID: 1,
				Offset: 0,
				Limit:  1,
			},
			wantErr: false,
			want: &myResponse{Data: []*models.Request{
				{
					ID:          1,
					Type:        "1",
					UserID:      1,
					Time:        1,
					Description: "1",
					Status:      "1",
				},
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newService(defaultLogger, tt.repo)
			got, err := s.My(context.Background(), tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("My() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("My() got = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestService_Create(t *testing.T) {
	tests := []struct {
		name    string
		repo    requestsRepository
		req     *createRequest
		wantErr bool
		want    *createResponse
	}{
		{
			name: "error from db",
			repo: &requestsRepositoryMock{
				CreateFunc: func(_ *models.Request) (uint, error) {
					return 0, errTestError
				},
			},
			req: &createRequest{
				Type:        "",
				Time:        0,
				UserID:      0,
				Description: "",
			},
			wantErr: true,
			want:    &createResponse{ID: 0},
		},
		{
			name: "ok",
			repo: &requestsRepositoryMock{
				CreateFunc: func(_ *models.Request) (u uint, err error) {
					return 1, nil
				},
			},
			req: &createRequest{
				Type:        "",
				Time:        0,
				UserID:      0,
				Description: "",
			},
			wantErr: false,
			want:    &createResponse{ID: 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newService(defaultLogger, tt.repo)
			got, err := s.Create(context.Background(), tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Create() got = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestService_Delete(t *testing.T) {
	tests := []struct {
		name    string
		repo    requestsRepository
		req     *byIDRequest
		wantErr bool
	}{
		{
			name: "error no request",
			repo: &requestsRepositoryMock{
				DeleteFunc: func(_ uint, _ uint) error {
					return errTestError
				},
			},
			req: &byIDRequest{
				UserID: 1,
				ID:     1,
			},
			wantErr: true,
		},
		{
			name: "ok",
			repo: &requestsRepositoryMock{
				DeleteFunc: func(_ uint, _ uint) error {
					return nil
				},
			},
			req: &byIDRequest{
				UserID: 1,
				ID:     1,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newService(defaultLogger, tt.repo)
			err := s.Delete(context.Background(), tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
