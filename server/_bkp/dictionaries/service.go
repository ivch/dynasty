package dictionaries

import (
	"context"

	"github.com/rs/zerolog"

	"github.com/ivch/dynasty/models/dto"
	"github.com/ivch/dynasty/models/entities"
)

type Service interface {
	BuildingsList(ctx context.Context) (*dto.BuildingsDictionaryResposnse, error)
	EntriesList(ctx context.Context, buildingID uint) (*dto.EntriesDictionaryResponse, error)
}

type dictRepository interface {
	BuildingsList() ([]*entities.Building, error)
	EntriesByBuilding(id uint) ([]*entities.Entry, error)
}

type service struct {
	repo dictRepository
}

func newService(repo dictRepository, logger *zerolog.Logger) Service {
	s := &service{
		repo: repo,
	}
	svc := newLoggingMiddleware(logger, s)
	return svc
}

func (s *service) EntriesList(_ context.Context, buildingID uint) (*dto.EntriesDictionaryResponse, error) {
	list, err := s.repo.EntriesByBuilding(buildingID)
	if err != nil {
		return nil, err
	}

	res := make([]*dto.Entry, len(list))
	for i, e := range list {
		res[i] = &dto.Entry{
			ID:   e.ID,
			Name: e.Name,
		}
	}

	return &dto.EntriesDictionaryResponse{Data: res}, nil
}

func (s *service) BuildingsList(_ context.Context) (*dto.BuildingsDictionaryResposnse, error) {
	list, err := s.repo.BuildingsList()
	if err != nil {
		return nil, err
	}

	res := make([]*dto.Building, len(list))
	for i, b := range list {
		res[i] = &dto.Building{
			ID:      b.ID,
			Name:    b.Name,
			Address: b.Address,
		}
	}

	return &dto.BuildingsDictionaryResposnse{Data: res}, nil
}
