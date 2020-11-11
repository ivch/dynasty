package dictionaries

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/rs/zerolog"

	"github.com/ivch/dynasty/models/dto"
	"github.com/ivch/dynasty/models/entities"
)

var (
	defaultLogger *zerolog.Logger
	errTestError  = errors.New("some err")
)

func TestMain(m *testing.M) {
	logger := zerolog.New(ioutil.Discard)
	defaultLogger = &logger
	os.Exit(m.Run())
}

func TestService_EntriesList(t *testing.T) {
	tests := []struct {
		name    string
		repo    dictRepository
		req     uint
		wantErr bool
		want    *dto.EntriesDictionaryResponse
	}{
		{
			name: "error from repo",
			repo: &dictRepositoryMock{
				EntriesByBuildingFunc: func(_ uint) ([]*entities.Entry, error) {
					return nil, errTestError
				},
			},
			wantErr: true,
		},
		{
			name: "ok",
			repo: &dictRepositoryMock{
				EntriesByBuildingFunc: func(id uint) ([]*entities.Entry, error) {
					return []*entities.Entry{
						{
							ID:   1,
							Name: "1",
						},
					}, nil
				},
			},
			req:     10,
			wantErr: false,
			want: &dto.EntriesDictionaryResponse{
				Data: []*dto.Entry{
					{
						ID:   1,
						Name: "1",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newService(tt.repo, defaultLogger)
			got, err := s.EntriesList(context.Background(), tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("EntriesList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EntriesList() got = %#v, want %#v", got, tt.want)
			}
		})
	}
}
func TestService_BuildingsList(t *testing.T) {
	tests := []struct {
		name    string
		repo    dictRepository
		wantErr bool
		want    *dto.BuildingsDictionaryResposnse
	}{
		{
			name: "error from repo",
			repo: &dictRepositoryMock{
				BuildingsListFunc: func() ([]*entities.Building, error) {
					return nil, errTestError
				},
			},
			wantErr: true,
		},
		{
			name: "ok",
			repo: &dictRepositoryMock{
				BuildingsListFunc: func() ([]*entities.Building, error) {
					return []*entities.Building{
						{
							ID:      1,
							Name:    "1",
							Address: "1",
						},
					}, nil
				},
			},
			wantErr: false,
			want: &dto.BuildingsDictionaryResposnse{
				Data: []*dto.Building{
					{
						ID:      1,
						Name:    "1",
						Address: "1",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newService(tt.repo, defaultLogger)
			got, err := s.BuildingsList(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("BuildingsList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BuildingsList() got = %#v, want %#v", got, tt.want)
			}
		})
	}
}
