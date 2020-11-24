package requests

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/microcosm-cc/bluemonday"

	"github.com/ivch/dynasty/common/logger"
)

var (
	defaultLogger *logger.StdLog
	errTestError  = errors.New("some err")
	defaultPolicy = bluemonday.StrictPolicy()
)

func TestMain(m *testing.M) {
	defaultLogger = logger.NewStdLog(logger.WithWriter(ioutil.Discard))
	os.Exit(m.Run())
}

func TestService_Get(t *testing.T) {
	tests := []struct {
		name    string
		repo    requestsRepository
		req     *Request
		wantErr bool
		want    *Request
	}{
		{
			name: "error no request",
			repo: &requestsRepositoryMock{
				GetRequestByIDAndUserFunc: func(_ uint, _ uint) (*Request, error) {
					return nil, errTestError
				},
			},
			req: &Request{
				UserID: 1,
				ID:     1,
			},
			wantErr: true,
		},
		{
			name: "ok",
			repo: &requestsRepositoryMock{
				GetRequestByIDAndUserFunc: func(_ uint, _ uint) (*Request, error) {
					return &Request{
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
			req: &Request{
				UserID: 1,
				ID:     1,
			},
			wantErr: false,
			want: &Request{
				ID:          1,
				Type:        "1",
				UserID:      1,
				Time:        1,
				Description: "1",
				Status:      "1",
				Images:      []string{"a"},
				ImagesURL: []map[string]string{
					{
						"img":   fmt.Sprintf("cdnHost/%s%s", imgPathPrefix, "a"),
						"thumb": fmt.Sprintf("cdnHost/%s%s", thumbPathPrefix, "a"),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(defaultLogger, tt.repo, nil, "", "cdnHost")
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
		req     *UpdateRequest
		wantErr bool
	}{
		{
			name: "error no request",
			repo: &requestsRepositoryMock{
				GetRequestByIDAndUserFunc: func(_ uint, _ uint) (*Request, error) {
					return nil, errTestError
				},
			},
			req: &UpdateRequest{
				ID:     1,
				UserID: 1,
			},
			wantErr: true,
		},
		{
			name: "error type not updated",
			repo: &requestsRepositoryMock{
				GetRequestByIDAndUserFunc: func(_ uint, _ uint) (*Request, error) {
					return &Request{
						Type: "1",
					}, nil
				},
				UpdateFunc: func(req *UpdateRequest) error {
					if *req.Type != "2" {
						return errTestError
					}
					return nil
				},
			},
			req: &UpdateRequest{
				ID:     1,
				UserID: 1,
				Type:   func(s string) *string { return &s }("1"),
			},
			wantErr: true,
		},
		{
			name: "error description not updated",
			repo: &requestsRepositoryMock{
				GetRequestByIDAndUserFunc: func(_ uint, _ uint) (*Request, error) {
					return &Request{
						Description: "1",
					}, nil
				},
				UpdateFunc: func(req *UpdateRequest) error {
					if *req.Description != "2" {
						return errTestError
					}
					return nil
				},
			},
			req: &UpdateRequest{
				ID:          1,
				UserID:      1,
				Description: func(s string) *string { return &s }("1"),
			},
			wantErr: true,
		},
		{
			name: "error status not updated",
			repo: &requestsRepositoryMock{
				GetRequestByIDAndUserFunc: func(_ uint, _ uint) (*Request, error) {
					return &Request{
						Status: "1",
					}, nil
				},
				UpdateFunc: func(req *UpdateRequest) error {
					if *req.Status != "2" {
						return errTestError
					}
					return nil
				},
			},
			req: &UpdateRequest{
				ID:     1,
				UserID: 1,
				Status: func(s string) *string { return &s }("1"),
			},
			wantErr: true,
		},
		{
			name: "error time not updated",
			repo: &requestsRepositoryMock{
				GetRequestByIDAndUserFunc: func(_ uint, _ uint) (*Request, error) {
					return &Request{
						Time: 1,
					}, nil
				},
				UpdateFunc: func(req *UpdateRequest) error {
					if *req.Time != 2 {
						return errTestError
					}
					return nil
				},
			},
			req: &UpdateRequest{
				ID:     1,
				UserID: 1,
				Time:   func(s int64) *int64 { return &s }(1),
			},
			wantErr: true,
		},
		{
			name: "ok",
			repo: &requestsRepositoryMock{
				GetRequestByIDAndUserFunc: func(_ uint, _ uint) (*Request, error) {
					return &Request{
						Type: "1",
					}, nil
				},
				UpdateFunc: func(req *UpdateRequest) error {
					return nil
				},
			},
			req: &UpdateRequest{
				ID:     1,
				UserID: 1,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(defaultLogger, tt.repo, nil, "", "")
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
		req     *RequestListFilter
		wantErr bool
		want    []*Request
	}{
		{
			name: "error from db",
			repo: &requestsRepositoryMock{
				ListByUserFunc: func(_ *RequestListFilter) ([]*Request, error) {
					return nil, errTestError
				},
			},
			req: &RequestListFilter{
				UserID: 1,
				Offset: 0,
				Limit:  1,
			},
			wantErr: true,
		},
		{
			name: "ok",
			repo: &requestsRepositoryMock{
				ListByUserFunc: func(_ *RequestListFilter) ([]*Request, error) {
					return []*Request{
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
			req: &RequestListFilter{
				UserID: 1,
				Offset: 0,
				Limit:  1,
			},
			wantErr: false,
			want: []*Request{
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
							"img":   fmt.Sprintf("cdnHost/%s%s", imgPathPrefix, "a"),
							"thumb": fmt.Sprintf("cdnHost/%s%s", thumbPathPrefix, "a"),
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(defaultLogger, tt.repo, nil, "", "cdnHost")
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
		req     *Request
		wantErr bool
		want    *Request
	}{
		{
			name: "error from db",
			repo: &requestsRepositoryMock{
				CreateFunc: func(_ *Request) error {
					return errTestError
				},
			},
			req: &Request{
				Type:        "",
				Time:        0,
				UserID:      0,
				Description: "",
			},
			wantErr: true,
		},
		{
			name: "ok",
			repo: &requestsRepositoryMock{
				CreateFunc: func(req *Request) error {
					req.ID = 1
					req.Status = "new"
					return nil
				},
			},
			req: &Request{
				Type:        "",
				Time:        0,
				UserID:      0,
				Description: "",
			},
			wantErr: false,
			want:    &Request{ID: 1, Status: "new"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(defaultLogger, tt.repo, nil, "", "")
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
		s3cli   s3Client
		req     *Request
		wantErr bool
	}{
		{
			name: "error finding request",
			repo: &requestsRepositoryMock{
				GetRequestByIDAndUserFunc: func(_ uint, _ uint) (*Request, error) {
					return nil, errTestError
				},
			},
			req: &Request{
				UserID: 1,
				ID:     1,
			},
			wantErr: true,
		},
		{
			name: "error deleting request",
			repo: &requestsRepositoryMock{
				GetRequestByIDAndUserFunc: func(_ uint, _ uint) (*Request, error) {
					return &Request{ID: 1}, nil
				},
				DeleteFunc: func(_ uint, _ uint) error {
					return errTestError
				},
			},
			req: &Request{
				UserID: 1,
				ID:     1,
			},
			wantErr: true,
		},
		{
			name: "ok w/o images",
			repo: &requestsRepositoryMock{
				GetRequestByIDAndUserFunc: func(_ uint, _ uint) (*Request, error) {
					return &Request{ID: 1}, nil
				},
				DeleteFunc: func(_ uint, _ uint) error {
					return nil
				},
			},
			req: &Request{
				UserID: 1,
				ID:     1,
			},
			wantErr: false,
		},
		{
			name: "error deleting with files",
			repo: &requestsRepositoryMock{
				GetRequestByIDAndUserFunc: func(_ uint, _ uint) (*Request, error) {
					return &Request{ID: 1, Images: []string{"a"}}, nil
				},
				DeleteFunc: func(_ uint, _ uint) error {
					return errTestError
				},
			},
			req: &Request{
				UserID: 1,
				ID:     1,
			},
			s3cli: &s3ClientMock{
				DeleteObjectFunc: func(input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) {
					return nil, nil
				},
			},
			wantErr: true,
		},
		{
			name: "ok with files",
			repo: &requestsRepositoryMock{
				GetRequestByIDAndUserFunc: func(_ uint, _ uint) (*Request, error) {
					return &Request{ID: 1, Images: []string{"a"}}, nil
				},
				DeleteFunc: func(_ uint, _ uint) error {
					return nil
				},
			},
			req: &Request{
				UserID: 1,
				ID:     1,
			},
			s3cli: &s3ClientMock{
				DeleteObjectFunc: func(input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) {
					return nil, nil
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(defaultLogger, tt.repo, tt.s3cli, "", "")
			err := s.Delete(context.Background(), tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
