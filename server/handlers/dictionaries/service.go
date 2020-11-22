package dictionaries

import (
	"context"

	"github.com/ivch/dynasty/common/logger"
)

type dictRepository interface {
	BuildingsList() ([]*Building, error)
	EntriesByBuilding(id uint) ([]*Entry, error)
}

type Service struct {
	log  logger.Logger
	repo dictRepository
}

func New(log logger.Logger, repo dictRepository) *Service {
	s := Service{
		log:  log,
		repo: repo,
	}

	return &s
}

func (s *Service) EntriesList(_ context.Context, buildingID uint) ([]*Entry, error) {
	list, err := s.repo.EntriesByBuilding(buildingID)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (s *Service) BuildingsList(_ context.Context) ([]*Building, error) {
	list, err := s.repo.BuildingsList()
	if err != nil {
		return nil, err
	}

	return list, nil
}
