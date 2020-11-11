package users

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ivch/dynasty/models/dto"
	"github.com/ivch/dynasty/models/entities"
)

func TestHTTP_GetUser(t *testing.T) {
	tests := []struct {
		name     string
		svc      Service
		header   string
		wantErr  bool
		want     string
		wantCode int
	}{
		{
			name:     "error no id",
			svc:      nil,
			header:   "0",
			wantErr:  true,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "error wrong id",
			svc:      nil,
			header:   "a",
			wantErr:  true,
			wantCode: http.StatusBadRequest,
		},
		{
			name: "error service error",
			svc: &ServiceMock{
				UserByIDFunc: func(_ context.Context, _ uint) (*dto.UserByIDResponse, error) {
					return nil, errTestError
				},
			},
			header:   "1",
			wantErr:  true,
			wantCode: http.StatusInternalServerError,
		},
		{
			name: "ok",
			svc: &ServiceMock{
				UserByIDFunc: func(_ context.Context, _ uint) (*dto.UserByIDResponse, error) {
					return &dto.UserByIDResponse{
						ID:        1,
						Apartment: 1,
						FirstName: "1",
						LastName:  "1",
						Phone:     "1",
						Email:     "1",
						Role:      1,
						Building: &entities.Building{
							ID:      1,
							Name:    "1",
							Address: "1",
						},
					}, nil
				},
			},
			header:   "1",
			wantErr:  false,
			want:     `{"id":1,"apartment":1,"first_name":"1","last_name":"1","phone":"1","email":"1","role":1,"building":{"id":1,"name":"1","address":"1"}}`,
			wantCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := tt.svc
			h := newHTTPHandler(defaultLogger, svc, defaultPolicy)
			rr := httptest.NewRecorder()
			rq, _ := http.NewRequest("GET", "/v1/user", nil)
			rq.Header.Add("X-Auth-User", tt.header)
			h.ServeHTTP(rr, rq)
			if (rr.Code != tt.wantCode) && tt.wantErr {
				t.Errorf("Request error. status = %d, wantCode = %d, wantErr %v", rr.Code, tt.wantCode, tt.wantErr)
			}

			if !tt.wantErr && tt.want != strings.TrimSpace(rr.Body.String()) {
				t.Errorf("Response error, got = %s, want = %s", rr.Body.String(), tt.want)
			}
		})
	}
}

func TestHTTP_Register(t *testing.T) {
	tests := []struct {
		name    string
		svc     Service
		request string
		wantErr bool
		want    string
	}{
		{
			name:    "error unmarshal request",
			svc:     nil,
			request: "}",
			wantErr: true,
		},
		{
			name:    "error invalid password",
			svc:     nil,
			request: `{"email":"test@test.com","first_name":"John","last_name":"Doe","apartment":1,"phone":"380671234567","building_id": 2, "code":"1231"}`,
			wantErr: true,
		},
		{
			name:    "error password too short",
			svc:     nil,
			request: `{"email":"test@test.com","first_name":"John","last_name":"Doe","apartment":1,"password":"1213", "phone":"380671234567","building_id": 2, "code":"1231"}`,
			wantErr: true,
		},
		{
			name:    "error invalid phone",
			svc:     nil,
			request: `{"email":"test@test.com","first_name":"John","last_name":"Doe","apartment":1,"password":"1213123","building_id": 2, "code":"1231"}`,
			wantErr: true,
		},
		{
			name:    "error invalid first name",
			svc:     nil,
			request: `{"email":"test@test.com","last_name":"Doe","apartment":1,"password":"1213123", "phone":"380671234567","building_id": 2, "code":"1231"}`,
			wantErr: true,
		},
		{
			name:    "error invalid last name",
			svc:     nil,
			request: `{"email":"test@test.com","first_name":"John","apartment":1,"password":"1213123", "phone":"380671234567","building_id": 2, "code":"1231"}`,
			wantErr: true,
		},
		{
			name:    "error invalid building",
			svc:     nil,
			request: `{"email":"test@test.com","first_name":"John","last_name":"Doe","apartment":1,"password":"1213123", "phone":"380671234567", "code":"1231"}`,
			wantErr: true,
		},
		{
			name:    "error invalid apartment",
			svc:     nil,
			request: `{"email":"test@test.com","first_name":"John","last_name":"Doe","password":"1213123", "phone":"380671234567","building_id": 2, "code":"1231"}`,
			wantErr: true,
		},
		{
			name:    "error invalid email",
			svc:     nil,
			request: `{"first_name":"John","last_name":"Doe","apartment":1,"password":"1213123", "phone":"380671234567","building_id": 2, "code":"1231"}`,
			wantErr: true,
		},
		{
			name:    "error invalid email#2",
			svc:     nil,
			request: `{"email":"testst.com","first_name":"John","last_name":"Doe","apartment":1,"password":"1213123", "phone":"380671234567","building_id": 2, "code":"1231"}`,
			wantErr: true,
		},
		{
			name: "error service",
			svc: &ServiceMock{
				RegisterFunc: func(_ context.Context, _ *dto.UserRegisterRequest) (*dto.UserRegisterResponse, error) {
					return nil, errTestError
				},
			},
			request: `{"email":"test@test.com","first_name":"John","last_name":"Doe","apartment":1,"password":"1213123", "phone":"380671234567","building_id": 2, "code":"1231"}`,
			wantErr: true,
		},
		{
			name: "ok",
			svc: &ServiceMock{
				RegisterFunc: func(_ context.Context, _ *dto.UserRegisterRequest) (*dto.UserRegisterResponse, error) {
					return &dto.UserRegisterResponse{
						ID:    1,
						Phone: "380671234567",
					}, nil
				},
			},
			request: `{"email":"test@test.com","first_name":"John","last_name":"Doe","apartment":1,"password":"1213123", "phone":"+380671234567","building_id": 2, "code":"1231", "entry_id":1}`,
			wantErr: false,
			want:    `{"id":1,"phone":"380671234567"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := tt.svc
			h := newHTTPHandler(defaultLogger, svc, defaultPolicy)
			rr := httptest.NewRecorder()
			rq, _ := http.NewRequest("POST", "/v1/register", strings.NewReader(tt.request))
			h.ServeHTTP(rr, rq)
			if (rr.Code != http.StatusOK) != tt.wantErr {
				t.Errorf("Request error. status = %d, wantErr %v", rr.Code, tt.wantErr)
			}

			if !tt.wantErr && tt.want != strings.TrimSpace(rr.Body.String()) {
				t.Errorf("Response error, got = %s, want = %s", rr.Body.String(), tt.want)
			}
		})
	}
}

func TestHTTP_AddFamilyMember(t *testing.T) {
	tests := []struct {
		name    string
		svc     Service
		header  string
		request string
		wantErr bool
		want    string
	}{
		{
			name:    "error no user",
			svc:     nil,
			header:  "",
			wantErr: true,
		},
		{
			name:    "error unmarshal request",
			svc:     nil,
			header:  "1",
			request: "}",
			wantErr: true,
		},
		{
			name:    "error empty phone",
			svc:     nil,
			header:  "1",
			request: `{"phone":""}`,
			wantErr: true,
		},
		{
			name:   "error service",
			header: "1",
			svc: &ServiceMock{
				AddFamilyMemberFunc: func(_ context.Context, _ *dto.AddFamilyMemberRequest) (*dto.AddFamilyMemberResponse, error) {
					return nil, errTestError
				},
			},
			request: `{"phone":"123"}`,
			wantErr: true,
		},
		{
			name:   "ok",
			header: "1",
			svc: &ServiceMock{
				AddFamilyMemberFunc: func(_ context.Context, r *dto.AddFamilyMemberRequest) (*dto.AddFamilyMemberResponse, error) {
					return &dto.AddFamilyMemberResponse{Code: "123"}, nil
				},
			},
			request: `{"phone":"380671234567"}`,
			wantErr: false,
			want:    `{"code":"123"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := tt.svc
			h := newHTTPHandler(defaultLogger, svc, defaultPolicy)
			rr := httptest.NewRecorder()
			rq, _ := http.NewRequest("POST", "/v1/member", strings.NewReader(tt.request))
			rq.Header.Add("X-Auth-User", tt.header)
			h.ServeHTTP(rr, rq)
			if (rr.Code != http.StatusOK) != tt.wantErr {
				t.Errorf("Request error. status = %d, wantErr %v", rr.Code, tt.wantErr)
			}

			if !tt.wantErr && tt.want != strings.TrimSpace(rr.Body.String()) {
				t.Errorf("Response error, got = %s, want = %s", rr.Body.String(), tt.want)
			}
		})
	}
}

func TestHTTP_ListFamilyMembers(t *testing.T) {
	tests := []struct {
		name    string
		svc     Service
		header  string
		wantErr bool
		want    string
	}{
		{
			name:    "error no user",
			svc:     nil,
			header:  "",
			wantErr: true,
		},
		{
			name:   "error service",
			header: "1",
			svc: &ServiceMock{
				ListFamilyMembersFunc: func(_ context.Context, _ uint) (*dto.ListFamilyMembersResponse, error) {
					return nil, errTestError
				},
			},
			wantErr: true,
		},
		{
			name:   "ok",
			header: "1",
			svc: &ServiceMock{
				ListFamilyMembersFunc: func(_ context.Context, _ uint) (*dto.ListFamilyMembersResponse, error) {
					return &dto.ListFamilyMembersResponse{
						Data: []*dto.FamilyMember{
							{
								ID:     1,
								Name:   "1",
								Phone:  "1",
								Code:   "1",
								Active: false,
							},
						},
					}, nil
				},
			},
			wantErr: false,
			want:    `{"data":[{"id":1,"name":"1","phone":"1","code":"1","active":false}]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := tt.svc
			h := newHTTPHandler(defaultLogger, svc, defaultPolicy)
			rr := httptest.NewRecorder()
			rq, _ := http.NewRequest("GET", "/v1/members", nil)
			rq.Header.Add("X-Auth-User", tt.header)
			h.ServeHTTP(rr, rq)
			if (rr.Code != http.StatusOK) != tt.wantErr {
				t.Errorf("Request error. status = %d, wantErr %v", rr.Code, tt.wantErr)
			}

			if !tt.wantErr && tt.want != strings.TrimSpace(rr.Body.String()) {
				t.Errorf("Response error, got = %s, want = %s", rr.Body.String(), tt.want)
			}
		})
	}
}

func TestHTTP_DeleteFamilyMember(t *testing.T) {
	tests := []struct {
		name    string
		svc     Service
		header  string
		request string
		wantErr bool
	}{
		{
			name:    "error no user",
			svc:     nil,
			header:  "",
			request: "1",
			wantErr: true,
		},
		{
			name:    "error bad request",
			svc:     nil,
			header:  "1",
			request: "asd",
			wantErr: true,
		},
		{
			name:    "error zero member",
			svc:     nil,
			header:  "1",
			request: "0",
			wantErr: true,
		},
		{
			name:    "error same ids",
			svc:     nil,
			header:  "1",
			request: "1",
			wantErr: true,
		},
		{
			name:    "error service",
			header:  "1",
			request: "2",
			svc: &ServiceMock{
				DeleteFamilyMemberFunc: func(_ context.Context, _ *dto.DeleteFamilyMemberRequest) error {
					return errTestError
				},
			},
			wantErr: true,
		},
		{
			name:    "ok",
			header:  "1",
			request: "2",
			svc: &ServiceMock{
				DeleteFamilyMemberFunc: func(_ context.Context, _ *dto.DeleteFamilyMemberRequest) error {
					return nil
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := tt.svc
			h := newHTTPHandler(defaultLogger, svc, defaultPolicy)
			rr := httptest.NewRecorder()
			rq, _ := http.NewRequest("DELETE", "/v1/member/"+tt.request, nil)
			rq.Header.Add("X-Auth-User", tt.header)
			h.ServeHTTP(rr, rq)
			if (rr.Code != http.StatusOK) != tt.wantErr {
				t.Errorf("Request error. status = %d, wantErr %v", rr.Code, tt.wantErr)
			}
		})
	}
}
