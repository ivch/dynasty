package requests

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/ivch/dynasty/models/dto"
	"github.com/ivch/dynasty/models/entities"
)

func TestService_GuardRequestList(t *testing.T) {
	tests := []struct {
		name    string
		repo    requestsRepository
		req     *dto.RequestListFilterRequest
		wantErr bool
		want    *dto.RequestGuardListResponse
	}{
		{
			name: "error from db on search",
			repo: &requestsRepositoryMock{
				ListForGuardFunc: func(_ *dto.RequestListFilterRequest) ([]*entities.Request, error) {
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
			name: "error from db on count",
			repo: &requestsRepositoryMock{
				ListForGuardFunc: func(_ *dto.RequestListFilterRequest) ([]*entities.Request, error) {
					return nil, nil
				},
				CountForGuardFunc: func(_ *dto.RequestListFilterRequest) (int, error) {
					return 0, errTestError
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
				ListForGuardFunc: func(_ *dto.RequestListFilterRequest) ([]*entities.Request, error) {
					return []*entities.Request{
						{
							ID:          1,
							Type:        "1",
							UserID:      1,
							Time:        1,
							Description: "1",
							Status:      "1",
							User: &entities.User{
								Building:  entities.Building{Name: "1"},
								Entry:     entities.Entry{Name: "2"},
								Apartment: 1,
								Phone:     "1",
								FirstName: "1",
								LastName:  "1",
							},
							Images: []string{"a"},
						},
					}, nil
				},
				CountForGuardFunc: func(req *dto.RequestListFilterRequest) (int, error) {
					return 1, nil
				},
			},
			req: &dto.RequestListFilterRequest{
				UserID: 1,
				Offset: 0,
				Limit:  1,
			},
			wantErr: false,
			want: &dto.RequestGuardListResponse{
				Data: []*dto.RequestForGuard{
					{
						ID:          1,
						UserID:      1,
						Type:        "1",
						Time:        1,
						Description: "1",
						Status:      "1",
						UserName:    "1 1",
						Phone:       "1",
						Address:     "1, 2",
						Apartment:   1,
						Images: []map[string]string{
							{
								"img":   fmt.Sprintf("cdnHost/%s%s", imgPathPrefix, "a"),
								"thumb": fmt.Sprintf("cdnHost/%s%s", thumbPathPrefix, "a"),
							},
						},
					},
				},
				Count: 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newService(defaultLogger, tt.repo, nil, "", "cdnHost")
			got, err := s.GuardRequestList(context.Background(), tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GuardRequestList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GuardRequestList() got = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestService_GuardUpdateRequest(t *testing.T) {
	tests := []struct {
		name    string
		repo    requestsRepository
		req     *dto.GuardUpdateRequest
		wantErr bool
	}{
		{
			name: "error from db",
			repo: &requestsRepositoryMock{
				UpdateForGuardFunc: func(_ uint, _ string) error {
					return errTestError
				},
			},
			req: &dto.GuardUpdateRequest{
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
			req: &dto.GuardUpdateRequest{
				ID:     1,
				Status: "1",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newService(defaultLogger, tt.repo, nil, "", "")
			err := s.GuardUpdateRequest(context.Background(), tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GuardUpdateRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
