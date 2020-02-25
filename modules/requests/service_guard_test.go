package requests

import (
	"context"
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
		want    []*dto.RequestForGuard
	}{
		{
			name: "error from db",
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
								Building:  entities.Building{Address: "1"},
								Apartment: 1,
								Phone:     "1",
								FirstName: "1",
								LastName:  "1",
							},
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
			want: []*dto.RequestForGuard{
				{
					ID:          1,
					UserID:      1,
					Type:        "1",
					Time:        1,
					Description: "1",
					Status:      "1",
					UserName:    "1 1",
					Phone:       "1",
					Address:     "1",
					Apartment:   1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newService(defaultLogger, tt.repo)
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
			s := newService(defaultLogger, tt.repo)
			err := s.GuardUpdateRequest(context.Background(), tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GuardUpdateRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
