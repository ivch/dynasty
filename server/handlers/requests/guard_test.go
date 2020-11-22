package requests

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/ivch/dynasty/server/handlers/users"
)

func TestService_GuardRequestList(t *testing.T) {
	tests := []struct {
		name    string
		repo    requestsRepository
		req     *RequestListFilter
		wantErr bool
		want    []*Request
		wantCnt int
	}{
		{
			name: "error from db on search",
			repo: &requestsRepositoryMock{
				ListForGuardFunc: func(_ *RequestListFilter) ([]*Request, error) {
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
			name: "error from db on count",
			repo: &requestsRepositoryMock{
				ListForGuardFunc: func(_ *RequestListFilter) ([]*Request, error) {
					return nil, nil
				},
				CountForGuardFunc: func(_ *RequestListFilter) (int, error) {
					return 0, errTestError
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
				ListForGuardFunc: func(_ *RequestListFilter) ([]*Request, error) {
					return []*Request{
						{
							ID:          1,
							Type:        "1",
							UserID:      1,
							Time:        1,
							Description: "1",
							Status:      "1",
							User: &users.User{
								Building:  users.Building{Name: "1"},
								Apartment: 1,
								Phone:     "1",
								FirstName: "1",
								LastName:  "1",
							},
							Images: []string{"a"},
						},
					}, nil
				},
				CountForGuardFunc: func(req *RequestListFilter) (int, error) {
					return 1, nil
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
					UserID:      1,
					Type:        "1",
					Time:        1,
					Description: "1",
					Status:      "1",
					User: &users.User{
						Building:  users.Building{Name: "1"},
						Apartment: 1,
						Phone:     "1",
						FirstName: "1",
						LastName:  "1",
					},
					Images: []string{"a"},
					ImagesURL: []map[string]string{
						{
							"img":   fmt.Sprintf("cdnHost/%s%s", imgPathPrefix, "a"),
							"thumb": fmt.Sprintf("cdnHost/%s%s", thumbPathPrefix, "a"),
						},
					},
				},
			},
			wantCnt: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(defaultLogger, tt.repo, nil, "", "cdnHost")
			got, cnt, err := s.GuardRequestList(context.Background(), tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GuardRequestList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GuardRequestList() got = %#v, want %#v", got, tt.want)
			}
			if cnt != tt.wantCnt {
				t.Errorf("GuardRequestList() got = %#v, want %#v", cnt, tt.wantCnt)
			}
		})
	}
}

func TestService_GuardUpdateRequest(t *testing.T) {
	tests := []struct {
		name    string
		repo    requestsRepository
		req     *Request
		wantErr bool
	}{
		{
			name: "error from db",
			repo: &requestsRepositoryMock{
				UpdateForGuardFunc: func(_ uint, _ string) error {
					return errTestError
				},
			},
			req: &Request{
				ID:     1,
				Status: "1",
			},
			wantErr: true,
		},
		{
			name: "ok",
			repo: &requestsRepositoryMock{
				UpdateForGuardFunc: func(_ uint, _ string) error {
					return nil
				},
			},
			req: &Request{
				ID:     1,
				Status: "1",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(defaultLogger, tt.repo, nil, "", "")
			err := s.GuardUpdateRequest(context.Background(), tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GuardUpdateRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
