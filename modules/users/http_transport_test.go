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
			wantCode: http.StatusInternalServerError,
		},
		{
			name:     "error wrong id",
			svc:      nil,
			header:   "a",
			wantErr:  true,
			wantCode: http.StatusInternalServerError,
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
						Building: entities.Building{
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
			h := newHTTPHandler(defaultLogger, svc)
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
			request: `{"email":"test@test.com","first_name":"John","last_name":"Doe","apartment":1,"password":"1213123", "phone":"+380671234567","building_id": 2, "code":"1231"}`,
			wantErr: false,
			want:    `{"id":1,"phone":"380671234567"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := tt.svc
			h := newHTTPHandler(defaultLogger, svc)
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
