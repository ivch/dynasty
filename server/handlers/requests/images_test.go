package requests

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/service/s3"
)

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
		req     *Image
		wantErr bool
		want    *Image
	}{
		{
			name: "error getting request from db",
			req: &Image{
				UserID:    1,
				RequestID: 1,
				File:      loadFile("../../../test_image.png"),
			},
			repo: &requestsRepositoryMock{
				GetRequestByIDAndUserFunc: func(_ uint, _ uint) (*Request, error) {
					return nil, errTestError
				},
			},
			wantErr: true,
		},
		{
			name: "error too much files for request",
			req: &Image{
				UserID:    1,
				RequestID: 1,
				File:      loadFile("../../../test_image.png"),
			},
			repo: &requestsRepositoryMock{
				GetRequestByIDAndUserFunc: func(_ uint, _ uint) (*Request, error) {
					return &Request{Images: []string{"1", "2", "3"}}, nil
				},
			},
			wantErr: true,
		},
		{
			name: "error wrong file type",
			req: &Image{
				UserID:    1,
				RequestID: 1,
				File:      loadFile("../../../test_image.png"),
			},
			repo: &requestsRepositoryMock{
				GetRequestByIDAndUserFunc: func(_ uint, _ uint) (*Request, error) {
					return &Request{}, nil
				},
			},
			wantErr: true,
		},
		{
			name: "error upload to s3",
			req: &Image{
				UserID:    1,
				RequestID: 1,
				File:      loadFile("../../../test_image.jpeg"),
			},
			repo: &requestsRepositoryMock{
				GetRequestByIDAndUserFunc: func(_ uint, _ uint) (*Request, error) {
					return &Request{}, nil
				},
			},
			s3cli: &s3ClientMock{
				PutObjectFunc: func(_ *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
					return nil, errTestError
				},
			},
			wantErr: true,
		},
		{
			name: "error upload thumb to s3 and delete err",
			req: &Image{
				UserID:    1,
				RequestID: 1,
				File:      loadFile("../../../test_image.jpeg"),
			},
			repo: &requestsRepositoryMock{
				GetRequestByIDAndUserFunc: func(_ uint, _ uint) (*Request, error) {
					return &Request{}, nil
				},
			},
			s3cli: &s3ClientMock{
				PutObjectFunc: func(input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
					if !strings.Contains(*input.Key, thumbPathPrefix) {
						return nil, nil
					}
					return nil, errTestError
				},
				DeleteObjectFunc: func(_ *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) {
					return nil, errTestError
				},
			},
			wantErr: true,
		},
		{
			name: "error upload thumb to s3 and no err on delete",
			req: &Image{
				UserID:    1,
				RequestID: 1,
				File:      loadFile("../../../test_image.jpeg"),
			},
			repo: &requestsRepositoryMock{
				GetRequestByIDAndUserFunc: func(_ uint, _ uint) (*Request, error) {
					return &Request{}, nil
				},
			},
			s3cli: &s3ClientMock{
				PutObjectFunc: func(input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
					if !strings.Contains(*input.Key, thumbPathPrefix) {
						return nil, nil
					}
					return nil, errTestError
				},
				DeleteObjectFunc: func(_ *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) {
					return nil, nil
				},
			},
			wantErr: true,
		},
		{
			name: "error add image to db + err delete from s3",
			req: &Image{
				UserID:    1,
				RequestID: 1,
				File:      loadFile("../../../test_image.jpeg"),
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
				GetRequestByIDAndUserFunc: func(_ uint, _ uint) (*Request, error) {
					return &Request{}, nil
				},
				AddImageFunc: func(_ uint, _ uint, _ string) error {
					return errTestError
				},
			},
			wantErr: true,
		},
		{
			name: "error add image to db + err delete thumb from s3",
			req: &Image{
				UserID:    1,
				RequestID: 1,
				File:      loadFile("../../../test_image.jpeg"),
			},
			s3cli: &s3ClientMock{
				PutObjectFunc: func(_ *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
					return nil, nil
				},
				DeleteObjectFunc: func(input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) {
					if !strings.Contains(*input.Key, thumbPathPrefix) {
						return nil, nil
					}
					return nil, errTestError
				},
			},
			repo: &requestsRepositoryMock{
				GetRequestByIDAndUserFunc: func(_ uint, _ uint) (*Request, error) {
					return &Request{}, nil
				},
				AddImageFunc: func(_ uint, _ uint, _ string) error {
					return errTestError
				},
			},
			wantErr: true,
		},
		{
			name: "error add image to db",
			req: &Image{
				UserID:    1,
				RequestID: 1,
				File:      loadFile("../../../test_image.jpeg"),
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
				GetRequestByIDAndUserFunc: func(_ uint, _ uint) (*Request, error) {
					return &Request{}, nil
				},
				AddImageFunc: func(_ uint, _ uint, _ string) error {
					return errTestError
				},
			},
			wantErr: true,
		},
		{
			name: "ok",
			req: &Image{
				UserID:    1,
				RequestID: 1,
				File:      loadFile("../../../test_image.jpeg"),
			},
			s3cli: &s3ClientMock{
				PutObjectFunc: func(_ *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
					return nil, nil
				},
			},
			repo: &requestsRepositoryMock{
				GetRequestByIDAndUserFunc: func(_ uint, _ uint) (*Request, error) {
					return &Request{}, nil
				},
				AddImageFunc: func(_ uint, _ uint, _ string) error {
					return nil
				},
			},
			wantErr: false,
			want: &Image{
				URL:   "cdnHost/" + imgPathPrefix,
				Thumb: "cdnHost/" + thumbPathPrefix,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(defaultLogger, tt.repo, tt.s3cli, "", "cdnHost")
			res, err := s.UploadImage(context.Background(), tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("UploadImage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if res == nil && !tt.wantErr {
				t.Error("UploadImage() error = empty response")
				return
			}

			if !tt.wantErr && !strings.Contains(res.URL, "cdnHost/"+imgPathPrefix) {
				t.Errorf("UploadImage() error = wrong img path, %s", res.URL)
				return
			}

			if !tt.wantErr && !strings.Contains(res.Thumb, "cdnHost/"+thumbPathPrefix) {
				t.Errorf("UploadImage() error = wrong thumb path: %s", res.Thumb)
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
		req     *Image
		wantErr bool
	}{
		{
			name: "error deleting from db",
			req: &Image{
				UserID:    1,
				RequestID: 1,
				URL:       "1",
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
			req: &Image{
				UserID:    1,
				RequestID: 1,
				URL:       "1",
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
			req: &Image{
				UserID:    1,
				RequestID: 1,
				URL:       "1",
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
			req: &Image{
				UserID:    1,
				RequestID: 1,
				URL:       "1",
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
			s := New(defaultLogger, tt.repo, tt.s3cli, "", "")
			err := s.DeleteImage(context.Background(), tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteImage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
