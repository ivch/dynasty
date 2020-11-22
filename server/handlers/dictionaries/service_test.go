package dictionaries

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/ivch/dynasty/common/logger"
)

var (
	defaultLogger *logger.StdLog
	errTestError  = errors.New("some err")
)

func TestMain(m *testing.M) {
	defaultLogger = logger.NewStdLog(logger.WithWriter(ioutil.Discard))
	os.Exit(m.Run())
}

func TestService_EntriesList(t *testing.T) {
	tests := []struct {
		name    string
		repo    dictRepository
		req     uint
		wantErr bool
		want    []*Entry
	}{
		{
			name: "error from repo",
			repo: &dictRepositoryMock{
				EntriesByBuildingFunc: func(_ uint) ([]*Entry, error) {
					return nil, errTestError
				},
			},
			wantErr: true,
		},
		{
			name: "ok",
			repo: &dictRepositoryMock{
				EntriesByBuildingFunc: func(id uint) ([]*Entry, error) {
					return []*Entry{
						{
							ID:   1,
							Name: "1",
						},
					}, nil
				},
			},
			req:     10,
			wantErr: false,
			want: []*Entry{
				{
					ID:   1,
					Name: "1",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(defaultLogger, tt.repo)
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
		want    []*Building
	}{
		{
			name: "error from repo",
			repo: &dictRepositoryMock{
				BuildingsListFunc: func() ([]*Building, error) {
					return nil, errTestError
				},
			},
			wantErr: true,
		},
		{
			name: "ok",
			repo: &dictRepositoryMock{
				BuildingsListFunc: func() ([]*Building, error) {
					return []*Building{
						{
							ID:      1,
							Name:    "1",
							Address: "1",
						},
					}, nil
				},
			},
			wantErr: false,
			want: []*Building{
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
			s := New(defaultLogger, tt.repo)
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
