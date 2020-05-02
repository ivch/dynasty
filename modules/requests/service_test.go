package requests

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/microcosm-cc/bluemonday"
	"github.com/rs/zerolog"

	"github.com/ivch/dynasty/models/dto"
	"github.com/ivch/dynasty/models/entities"
)

var (
	defaultLogger *zerolog.Logger
	errTestError  = errors.New("some err")
	defaultPolicy = bluemonday.StrictPolicy()
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
		req     *dto.RequestByID
		wantErr bool
		want    *dto.RequestByIDResponse
	}{
		{
			name: "error no request",
			repo: &requestsRepositoryMock{
				GetRequestByIDAndUserFunc: func(_ uint, _ uint) (*entities.Request, error) {
					return nil, errTestError
				},
			},
			req: &dto.RequestByID{
				UserID: 1,
				ID:     1,
			},
			wantErr: true,
		},
		{
			name: "ok",
			repo: &requestsRepositoryMock{
				GetRequestByIDAndUserFunc: func(_ uint, _ uint) (*entities.Request, error) {
					return &entities.Request{
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
			req: &dto.RequestByID{
				UserID: 1,
				ID:     1,
			},
			wantErr: false,
			want: &dto.RequestByIDResponse{
				ID:          1,
				Type:        "1",
				UserID:      1,
				Time:        1,
				Description: "1",
				Status:      "1",
				Images:      []string{"/1/a"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newService(defaultLogger, tt.repo, nil, "", "")
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
		req     *dto.RequestUpdateRequest
		wantErr bool
	}{
		{
			name: "error no request",
			repo: &requestsRepositoryMock{
				GetRequestByIDAndUserFunc: func(_ uint, _ uint) (*entities.Request, error) {
					return nil, errTestError
				},
			},
			req: &dto.RequestUpdateRequest{
				ID:     1,
				UserID: 1,
			},
			wantErr: true,
		},
		{
			name: "error type not updated",
			repo: &requestsRepositoryMock{
				GetRequestByIDAndUserFunc: func(_ uint, _ uint) (*entities.Request, error) {
					return &entities.Request{
						Type: "1",
					}, nil
				},
				UpdateFunc: func(req *entities.Request) error {
					if req.Type != "2" {
						return errTestError
					}
					return nil
				},
			},
			req: &dto.RequestUpdateRequest{
				ID:     1,
				UserID: 1,
				Type:   nil,
			},
			wantErr: true,
		},
		{
			name: "error description not updated",
			repo: &requestsRepositoryMock{
				GetRequestByIDAndUserFunc: func(_ uint, _ uint) (*entities.Request, error) {
					return &entities.Request{
						Description: "1",
					}, nil
				},
				UpdateFunc: func(req *entities.Request) error {
					if req.Description != "2" {
						return errTestError
					}
					return nil
				},
			},
			req: &dto.RequestUpdateRequest{
				ID:          1,
				UserID:      1,
				Description: nil,
			},
			wantErr: true,
		},
		{
			name: "error status not updated",
			repo: &requestsRepositoryMock{
				GetRequestByIDAndUserFunc: func(_ uint, _ uint) (*entities.Request, error) {
					return &entities.Request{
						Status: "1",
					}, nil
				},
				UpdateFunc: func(req *entities.Request) error {
					if req.Status != "2" {
						return errTestError
					}
					return nil
				},
			},
			req: &dto.RequestUpdateRequest{
				ID:     1,
				UserID: 1,
				Status: nil,
			},
			wantErr: true,
		},
		{
			name: "error time not updated",
			repo: &requestsRepositoryMock{
				GetRequestByIDAndUserFunc: func(_ uint, _ uint) (*entities.Request, error) {
					return &entities.Request{
						Time: 1,
					}, nil
				},
				UpdateFunc: func(req *entities.Request) error {
					if req.Time != 2 {
						return errTestError
					}
					return nil
				},
			},
			req: &dto.RequestUpdateRequest{
				ID:     1,
				UserID: 1,
				Time:   nil,
			},
			wantErr: true,
		},
		{
			name: "ok",
			repo: &requestsRepositoryMock{
				GetRequestByIDAndUserFunc: func(_ uint, _ uint) (*entities.Request, error) {
					return &entities.Request{
						Type: "1",
					}, nil
				},
				UpdateFunc: func(req *entities.Request) error {
					return nil
				},
			},
			req: &dto.RequestUpdateRequest{
				ID:     1,
				UserID: 1,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newService(defaultLogger, tt.repo, nil, "", "")
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
		req     *dto.RequestListFilterRequest
		wantErr bool
		want    *dto.RequestMyResponse
	}{
		{
			name: "error from db",
			repo: &requestsRepositoryMock{
				ListByUserFunc: func(_ *dto.RequestListFilterRequest) ([]*entities.Request, error) {
					return nil, errTestError
				},
			},
			req: &dto.RequestListFilterRequest{
				UserID: 1,
				Offset: 0,
				Limit:  1,
			},
			wantErr: true,
		},
		{
			name: "ok",
			repo: &requestsRepositoryMock{
				ListByUserFunc: func(_ *dto.RequestListFilterRequest) ([]*entities.Request, error) {
					return []*entities.Request{
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
			req: &dto.RequestListFilterRequest{
				UserID: 1,
				Offset: 0,
				Limit:  1,
			},
			wantErr: false,
			want: &dto.RequestMyResponse{Data: []*dto.RequestByIDResponse{
				{
					ID:          1,
					Type:        "1",
					UserID:      1,
					Time:        1,
					Description: "1",
					Status:      "1",
					Images:      []string{"/1/a"},
				},
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newService(defaultLogger, tt.repo, nil, "", "")
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
		req     *dto.RequestCreateRequest
		wantErr bool
		want    *dto.RequestCreateResponse
	}{
		{
			name: "error from db",
			repo: &requestsRepositoryMock{
				CreateFunc: func(_ *entities.Request) (uint, error) {
					return 0, errTestError
				},
			},
			req: &dto.RequestCreateRequest{
				Type:        "",
				Time:        0,
				UserID:      0,
				Description: "",
			},
			wantErr: true,
			want:    &dto.RequestCreateResponse{ID: 0},
		},
		{
			name: "ok",
			repo: &requestsRepositoryMock{
				CreateFunc: func(_ *entities.Request) (u uint, err error) {
					return 1, nil
				},
			},
			req: &dto.RequestCreateRequest{
				Type:        "",
				Time:        0,
				UserID:      0,
				Description: "",
			},
			wantErr: false,
			want:    &dto.RequestCreateResponse{ID: 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newService(defaultLogger, tt.repo, nil, "", "")
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
		req     *dto.RequestByID
		wantErr bool
	}{
		{
			name: "error finding request",
			repo: &requestsRepositoryMock{
				GetRequestByIDAndUserFunc: func(_ uint, _ uint) (*entities.Request, error) {
					return nil, errTestError
				},
			},
			req: &dto.RequestByID{
				UserID: 1,
				ID:     1,
			},
			wantErr: true,
		},
		{
			name: "error deleting request",
			repo: &requestsRepositoryMock{
				GetRequestByIDAndUserFunc: func(_ uint, _ uint) (*entities.Request, error) {
					return &entities.Request{ID: 1}, nil
				},
				DeleteFunc: func(_ uint, _ uint) error {
					return errTestError
				},
			},
			req: &dto.RequestByID{
				UserID: 1,
				ID:     1,
			},
			wantErr: true,
		},
		{
			name: "ok w/o images",
			repo: &requestsRepositoryMock{
				GetRequestByIDAndUserFunc: func(_ uint, _ uint) (*entities.Request, error) {
					return &entities.Request{ID: 1}, nil
				},
				DeleteFunc: func(_ uint, _ uint) error {
					return nil
				},
			},
			req: &dto.RequestByID{
				UserID: 1,
				ID:     1,
			},
			wantErr: false,
		},
		{
			name: "error deleting with files",
			repo: &requestsRepositoryMock{
				GetRequestByIDAndUserFunc: func(_ uint, _ uint) (*entities.Request, error) {
					return &entities.Request{ID: 1, Images: []string{"a"}}, nil
				},
				DeleteFunc: func(_ uint, _ uint) error {
					return errTestError
				},
			},
			req: &dto.RequestByID{
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
				GetRequestByIDAndUserFunc: func(_ uint, _ uint) (*entities.Request, error) {
					return &entities.Request{ID: 1, Images: []string{"a"}}, nil
				},
				DeleteFunc: func(_ uint, _ uint) error {
					return nil
				},
			},
			req: &dto.RequestByID{
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
			s := newService(defaultLogger, tt.repo, tt.s3cli, "", "")
			err := s.Delete(context.Background(), tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestService_UploadImage(t *testing.T) {
	loadFile := func(filename string) []byte {
		f, err := os.Open(filename)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		fileBytes, err := ioutil.ReadAll(f)
		if err != nil {
			log.Fatal(err)
		}
		return fileBytes
	}

	tests := []struct {
		name    string
		repo    requestsRepository
		s3cli   s3Client
		req     *dto.UploadImageRequest
		wantErr bool
	}{
		{
			name: "error wrong file type",
			req: &dto.UploadImageRequest{
				UserID:    1,
				RequestID: 1,
				File:      loadFile("../../test_image.png"),
			},
			wantErr: true,
		},
		{
			name: "error upload to s3",
			req: &dto.UploadImageRequest{
				UserID:    1,
				RequestID: 1,
				File:      loadFile("../../test_image.jpeg"),
			},
			s3cli: &s3ClientMock{
				PutObjectFunc: func(_ *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
					return nil, errTestError
				},
			},
			wantErr: true,
		},
		{
			name: "error add image to db + err delete from s3",
			req: &dto.UploadImageRequest{
				UserID:    1,
				RequestID: 1,
				File:      loadFile("../../test_image.jpeg"),
			},
			s3cli: &s3ClientMock{
				PutObjectFunc: func(_ *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
					return nil, nil
				},
				DeleteObjectFunc: func(_ *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) {
					return nil, errTestError
				},
			},
			repo: &requestsRepositoryMock{
				AddImageFunc: func(_ uint, _ uint, _ string) error {
					return errTestError
				},
			},
			wantErr: true,
		},
		{
			name: "error add image to db",
			req: &dto.UploadImageRequest{
				UserID:    1,
				RequestID: 1,
				File:      loadFile("../../test_image.jpeg"),
			},
			s3cli: &s3ClientMock{
				PutObjectFunc: func(_ *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
					return nil, nil
				},
				DeleteObjectFunc: func(_ *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) {
					return nil, nil
				},
			},
			repo: &requestsRepositoryMock{
				AddImageFunc: func(_ uint, _ uint, _ string) error {
					return errTestError
				},
			},
			wantErr: true,
		},
		{
			name: "ok",
			req: &dto.UploadImageRequest{
				UserID:    1,
				RequestID: 1,
				File:      loadFile("../../test_image.jpeg"),
			},
			s3cli: &s3ClientMock{
				PutObjectFunc: func(_ *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
					return nil, nil
				},
			},
			repo: &requestsRepositoryMock{
				AddImageFunc: func(_ uint, _ uint, _ string) error {
					return nil
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newService(defaultLogger, tt.repo, tt.s3cli, "", "")
			_, err := s.UploadImage(context.Background(), tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("UploadImage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestService_DeleteImage(t *testing.T) {
	tests := []struct {
		name    string
		repo    requestsRepository
		s3cli   s3Client
		req     *dto.DeleteImageRequest
		wantErr bool
	}{
		{
			name: "error deleting from db",
			req: &dto.DeleteImageRequest{
				UserID:    1,
				RequestID: 1,
				Filepath:  "1",
			},
			repo: &requestsRepositoryMock{
				DeleteImageFunc: func(_ uint, _ uint, _ string) error {
					return errTestError
				},
			},
			wantErr: true,
		},
		{
			name: "error deleting from s3 + err add to db",
			req: &dto.DeleteImageRequest{
				UserID:    1,
				RequestID: 1,
				Filepath:  "1",
			},
			repo: &requestsRepositoryMock{
				DeleteImageFunc: func(_ uint, _ uint, _ string) error {
					return nil
				},
				AddImageFunc: func(_ uint, _ uint, _ string) error {
					return errTestError
				},
			},
			s3cli: &s3ClientMock{
				DeleteObjectFunc: func(_ *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) {
					return nil, errTestError
				},
			},
			wantErr: true,
		},
		{
			name: "error deleting from s3",
			req: &dto.DeleteImageRequest{
				UserID:    1,
				RequestID: 1,
				Filepath:  "1",
			},
			repo: &requestsRepositoryMock{
				DeleteImageFunc: func(_ uint, _ uint, _ string) error {
					return nil
				},
				AddImageFunc: func(_ uint, _ uint, _ string) error {
					return nil
				},
			},
			s3cli: &s3ClientMock{
				DeleteObjectFunc: func(_ *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) {
					return nil, errTestError
				},
			},
			wantErr: true,
		},
		{
			name: "ok",
			req: &dto.DeleteImageRequest{
				UserID:    1,
				RequestID: 1,
				Filepath:  "1",
			},
			repo: &requestsRepositoryMock{
				DeleteImageFunc: func(_ uint, _ uint, _ string) error {
					return nil
				},
			},
			s3cli: &s3ClientMock{
				DeleteObjectFunc: func(_ *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) {
					return nil, nil
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newService(defaultLogger, tt.repo, tt.s3cli, "", "")
			err := s.DeleteImage(context.Background(), tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteImage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
