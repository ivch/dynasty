package dictionaries

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ivch/dynasty/models/dto"
)

func TestHTTP_EntriesList(t *testing.T) {
	tests := []struct {
		name     string
		req      string
		svc      Service
		wantErr  bool
		want     string
		wantCode int
	}{
		{
			name:     "error bad id",
			req:      "asd",
			wantErr:  true,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "error bad id #2",
			req:      "0",
			wantErr:  true,
			wantCode: http.StatusBadRequest,
		},
		{
			name: "error service error",
			req:  "1",
			svc: &ServiceMock{
				EntriesListFunc: func(_ context.Context, _ uint) (*dto.EntriesDictionaryResponse, error) {
					return nil, errTestError
				},
			},
			wantErr:  true,
			wantCode: http.StatusInternalServerError,
		},
		{
			name: "ok",
			req:  "1",
			svc: &ServiceMock{
				EntriesListFunc: func(_ context.Context, _ uint) (*dto.EntriesDictionaryResponse, error) {
					return &dto.EntriesDictionaryResponse{
						Data: []*dto.Entry{
							{
								ID:   1,
								Name: "1",
							},
						},
					}, nil
				},
			},
			want:     `{"data":[{"id":1,"name":"1"}]}`,
			wantCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := tt.svc
			h := newHTTPHandler(defaultLogger, svc)
			rr := httptest.NewRecorder()
			rq, _ := http.NewRequest("GET", "/v1/building/"+tt.req+"/entries", nil)
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

func TestHTTP_BuildingsList(t *testing.T) {
	tests := []struct {
		name     string
		svc      Service
		wantErr  bool
		want     string
		wantCode int
	}{
		{
			name: "error service error",
			svc: &ServiceMock{
				BuildingsListFunc: func(_ context.Context) (*dto.BuildingsDictionaryResposnse, error) {
					return nil, errTestError
				},
			},
			wantErr:  true,
			wantCode: http.StatusInternalServerError,
		},
		{
			name: "ok",
			svc: &ServiceMock{
				BuildingsListFunc: func(_ context.Context) (*dto.BuildingsDictionaryResposnse, error) {
					return &dto.BuildingsDictionaryResposnse{
						Data: []*dto.Building{
							{
								ID:      1,
								Name:    "1",
								Address: "1",
							},
						},
					}, nil
				},
			},
			want:     `{"data":[{"id":1,"name":"1","address":"1"}]}`,
			wantCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := tt.svc
			h := newHTTPHandler(defaultLogger, svc)
			rr := httptest.NewRecorder()
			rq, _ := http.NewRequest("GET", "/v1/buildings", nil)
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
