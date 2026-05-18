package dictionaries_test

import (
	"context"
	"errors"
	"io"
	"os"
	"reflect"
	"testing"

	"github.com/ivch/dynasty/common/logger"
	"github.com/ivch/dynasty/server/handlers/dictionaries"
)

var (
	defaultLogger *logger.StdLog
	errTestError  = errors.New("some err")
)

func TestMain(m *testing.M) {
	defaultLogger = logger.NewStdLog(logger.WithWriter(io.Discard))
	os.Exit(m.Run())
}

func TestService_EntriesList(t *testing.T) {
	tests := []struct {
		name    string
		repo    dictionaries.DictRepository
		req     uint
		wantErr bool
		want    []*dictionaries.Entry
	}{
		{
			name: "error from repo",
			repo: &dictionaries.DictRepositoryMock{
				EntriesByBuildingFunc: func(_ uint) ([]*dictionaries.Entry, error) {
					return nil, errTestError
				},
			},
			wantErr: true,
		},
		{
			name: "ok",
			repo: &dictionaries.DictRepositoryMock{
				EntriesByBuildingFunc: func(id uint) ([]*dictionaries.Entry, error) {
					return []*dictionaries.Entry{
						{
							ID:   1,
							Name: "1",
						},
					}, nil
				},
			},
			req:     10,
			wantErr: false,
			want: []*dictionaries.Entry{
				{
					ID:   1,
					Name: "1",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := dictionaries.New(defaultLogger, tt.repo)
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
		repo    dictionaries.DictRepository
		wantErr bool
		want    []*dictionaries.Building
	}{
		{
			name: "error from repo",
			repo: &dictionaries.DictRepositoryMock{
				BuildingsListFunc: func() ([]*dictionaries.Building, error) {
					return nil, errTestError
				},
			},
			wantErr: true,
		},
		{
			name: "ok",
			repo: &dictionaries.DictRepositoryMock{
				BuildingsListFunc: func() ([]*dictionaries.Building, error) {
					return []*dictionaries.Building{
						{
							ID:      1,
							Name:    "1",
							Address: "1",
						},
					}, nil
				},
			},
			wantErr: false,
			want: []*dictionaries.Building{
				{
					ID:      1,
					Name:    "1",
					Address: "1",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := dictionaries.New(defaultLogger, tt.repo)
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
