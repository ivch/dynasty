package requests_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/ivch/dynasty/common/logger"
	"github.com/ivch/dynasty/server/handlers/requests"
)

var (
	defaultLogger *logger.StdLog
	errTestError  = errors.New("some err")
)

func TestMain(m *testing.M) {
	defaultLogger = logger.NewStdLog(logger.WithWriter(io.Discard))
	os.Exit(m.Run())
}

func TestService_Get(t *testing.T) {
	tests := []struct {
		name    string
		repo    requests.RequestsRepository
		req     *requests.Request
		wantErr bool
		want    *requests.Request
	}{
		{
			name: "error no request",
			repo: &requests.RequestsRepositoryMock{
				GetRequestByIDAndUserFunc: func(_ uint, _ uint) (*requests.Request, error) {
					return nil, errTestError
				},
			},
			req: &requests.Request{
				UserID: 1,
				ID:     1,
			},
			wantErr: true,
		},
		{
			name: "ok",
			repo: &requests.RequestsRepositoryMock{
				GetRequestByIDAndUserFunc: func(_ uint, _ uint) (*requests.Request, error) {
					return &requests.Request{
						ID:          1,
						Type:        "1",
						UserID:      1,
						Time:        1,
						Description: "1",
						Status:      "1",
						Images:      []string{"a"},
					}, nil
				},
			},
			req: &requests.Request{
				UserID: 1,
				ID:     1,
			},
			wantErr: false,
			want: &requests.Request{
				ID:          1,
				Type:        "1",
				UserID:      1,
				Time:        1,
				Description: "1",
				Status:      "1",
				Images:      []string{"a"},
				ImagesURL: []map[string]string{
					{
						"img":   fmt.Sprintf("cdnHost/%s%s", requests.ImgPathPrefix, "a"),
						"thumb": fmt.Sprintf("cdnHost/%s%s", requests.ThumbPathPrefix, "a"),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := requests.New(defaultLogger, tt.repo, nil, "", "cdnHost")
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
		repo    requests.RequestsRepository
		req     *requests.UpdateRequest
		wantErr bool
	}{
		{
			name: "error no request",
			repo: &requests.RequestsRepositoryMock{
				GetRequestByIDAndUserFunc: func(_ uint, _ uint) (*requests.Request, error) {
					return nil, errTestError
				},
			},
			req: &requests.UpdateRequest{
				ID:     1,
				UserID: 1,
			},
			wantErr: true,
		},
		{
			name: "error type not updated",
			repo: &requests.RequestsRepositoryMock{
				GetRequestByIDAndUserFunc: func(_ uint, _ uint) (*requests.Request, error) {
					return &requests.Request{
						Type: "1",
					}, nil
				},
				UpdateFunc: func(req *requests.UpdateRequest) error {
					if *req.Type != "2" {
						return errTestError
					}
					return nil
				},
			},
			req: &requests.UpdateRequest{
				ID:     1,
				UserID: 1,
				Type:   func(s string) *string { return &s }("1"),
			},
			wantErr: true,
		},
		{
			name: "error description not updated",
			repo: &requests.RequestsRepositoryMock{
				GetRequestByIDAndUserFunc: func(_ uint, _ uint) (*requests.Request, error) {
					return &requests.Request{
						Description: "1",
					}, nil
				},
				UpdateFunc: func(req *requests.UpdateRequest) error {
					if *req.Description != "2" {
						return errTestError
					}
					return nil
				},
			},
			req: &requests.UpdateRequest{
				ID:          1,
				UserID:      1,
				Description: func(s string) *string { return &s }("1"),
			},
			wantErr: true,
		},
		{
			name: "error status not updated",
			repo: &requests.RequestsRepositoryMock{
				GetRequestByIDAndUserFunc: func(_ uint, _ uint) (*requests.Request, error) {
					return &requests.Request{
						Status: "1",
					}, nil
				},
				UpdateFunc: func(req *requests.UpdateRequest) error {
					if *req.Status != "2" {
						return errTestError
					}
					return nil
				},
			},
			req: &requests.UpdateRequest{
				ID:     1,
				UserID: 1,
				Status: func(s string) *string { return &s }("1"),
			},
			wantErr: true,
		},
		{
			name: "error time not updated",
			repo: &requests.RequestsRepositoryMock{
				GetRequestByIDAndUserFunc: func(_ uint, _ uint) (*requests.Request, error) {
					return &requests.Request{
						Time: 1,
					}, nil
				},
				UpdateFunc: func(req *requests.UpdateRequest) error {
					if *req.Time != 2 {
						return errTestError
					}
					return nil
				},
			},
			req: &requests.UpdateRequest{
				ID:     1,
				UserID: 1,
				Time:   func(s int64) *int64 { return &s }(1),
			},
			wantErr: true,
		},
		{
			name: "ok",
			repo: &requests.RequestsRepositoryMock{
				GetRequestByIDAndUserFunc: func(_ uint, _ uint) (*requests.Request, error) {
					return &requests.Request{
						Type: "1",
					}, nil
				},
				UpdateFunc: func(req *requests.UpdateRequest) error {
					return nil
				},
			},
			req: &requests.UpdateRequest{
				ID:     1,
				UserID: 1,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := requests.New(defaultLogger, tt.repo, nil, "", "")
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
		repo    requests.RequestsRepository
		req     *requests.RequestListFilter
		wantErr bool
		want    []*requests.Request
	}{
		{
			name: "error from db",
			repo: &requests.RequestsRepositoryMock{
				ListByUserFunc: func(_ *requests.RequestListFilter) ([]*requests.Request, error) {
					return nil, errTestError
				},
			},
			req: &requests.RequestListFilter{
				UserID: 1,
				Offset: 0,
				Limit:  1,
			},
			wantErr: true,
		},
		{
			name: "ok",
			repo: &requests.RequestsRepositoryMock{
				ListByUserFunc: func(_ *requests.RequestListFilter) ([]*requests.Request, error) {
					return []*requests.Request{
						{
							ID:          1,
							Type:        "1",
							UserID:      1,
							Time:        1,
							Description: "1",
							Status:      "1",
							Images:      []string{"a"},
						},
					}, nil
				},
			},
			req: &requests.RequestListFilter{
				UserID: 1,
				Offset: 0,
				Limit:  1,
			},
			wantErr: false,
			want: []*requests.Request{
				{
					ID:          1,
					Type:        "1",
					UserID:      1,
					Time:        1,
					Description: "1",
					Status:      "1",
					Images:      []string{"a"},
					ImagesURL: []map[string]string{
						{
							"img":   fmt.Sprintf("cdnHost/%s%s", requests.ImgPathPrefix, "a"),
							"thumb": fmt.Sprintf("cdnHost/%s%s", requests.ThumbPathPrefix, "a"),
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := requests.New(defaultLogger, tt.repo, nil, "", "cdnHost")
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
		repo    requests.RequestsRepository
		req     *requests.Request
		wantErr bool
		want    *requests.Request
	}{
		{
			name: "error req limit exceeded",
			repo: &requests.RequestsRepositoryMock{
				ListByUserFunc: func(r *requests.RequestListFilter) ([]*requests.Request, error) {
					res := make([]*requests.Request, 22)
					return res, nil
				},
			},
			req: &requests.Request{
				Type:        "",
				Time:        0,
				UserID:      0,
				Description: "",
			},
			wantErr: true,
		},
		{
			name: "error from db",
			repo: &requests.RequestsRepositoryMock{
				ListByUserFunc: func(r *requests.RequestListFilter) ([]*requests.Request, error) {
					res := make([]*requests.Request, 1)
					return res, nil
				},
				CreateFunc: func(_ *requests.Request) error {
					return errTestError
				},
			},
			req: &requests.Request{
				Type:        "",
				Time:        0,
				UserID:      0,
				Description: "",
			},
			wantErr: true,
		},
		{
			name: "ok",
			repo: &requests.RequestsRepositoryMock{
				ListByUserFunc: func(r *requests.RequestListFilter) ([]*requests.Request, error) {
					res := make([]*requests.Request, 1)
					return res, nil
				},
				CreateFunc: func(req *requests.Request) error {
					req.ID = 1
					req.Status = "new"
					return nil
				},
			},
			req: &requests.Request{
				Type:        "",
				Time:        0,
				UserID:      0,
				Description: "",
			},
			wantErr: false,
			want:    &requests.Request{ID: 1, Status: "new"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := requests.New(defaultLogger, tt.repo, nil, "", "")
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
		repo    requests.RequestsRepository
		s3cli   requests.S3Client
		req     *requests.Request
		wantErr bool
	}{
		{
			name: "error finding request",
			repo: &requests.RequestsRepositoryMock{
				GetRequestByIDAndUserFunc: func(_ uint, _ uint) (*requests.Request, error) {
					return nil, errTestError
				},
			},
			req: &requests.Request{
				UserID: 1,
				ID:     1,
			},
			wantErr: true,
		},
		{
			name: "error deleting request",
			repo: &requests.RequestsRepositoryMock{
				GetRequestByIDAndUserFunc: func(_ uint, _ uint) (*requests.Request, error) {
					return &requests.Request{ID: 1}, nil
				},
				DeleteFunc: func(_ uint, _ uint) error {
					return errTestError
				},
			},
			req: &requests.Request{
				UserID: 1,
				ID:     1,
			},
			wantErr: true,
		},
		{
			name: "ok w/o images",
			repo: &requests.RequestsRepositoryMock{
				GetRequestByIDAndUserFunc: func(_ uint, _ uint) (*requests.Request, error) {
					return &requests.Request{ID: 1}, nil
				},
				DeleteFunc: func(_ uint, _ uint) error {
					return nil
				},
			},
			req: &requests.Request{
				UserID: 1,
				ID:     1,
			},
			wantErr: false,
		},
		{
			name: "error deleting with files",
			repo: &requests.RequestsRepositoryMock{
				GetRequestByIDAndUserFunc: func(_ uint, _ uint) (*requests.Request, error) {
					return &requests.Request{ID: 1, Images: []string{"a"}}, nil
				},
				DeleteFunc: func(_ uint, _ uint) error {
					return errTestError
				},
			},
			req: &requests.Request{
				UserID: 1,
				ID:     1,
			},
			s3cli: &requests.S3ClientMock{
				DeleteObjectFunc: func(input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) {
					return nil, nil
				},
			},
			wantErr: true,
		},
		{
			name: "ok with files",
			repo: &requests.RequestsRepositoryMock{
				GetRequestByIDAndUserFunc: func(_ uint, _ uint) (*requests.Request, error) {
					return &requests.Request{ID: 1, Images: []string{"a"}}, nil
				},
				DeleteFunc: func(_ uint, _ uint) error {
					return nil
				},
			},
			req: &requests.Request{
				UserID: 1,
				ID:     1,
			},
			s3cli: &requests.S3ClientMock{
				DeleteObjectFunc: func(input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) {
					return nil, nil
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := requests.New(defaultLogger, tt.repo, tt.s3cli, "", "")
			err := s.Delete(context.Background(), tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
