package repository

import (
	"github.com/jinzhu/gorm"

	"github.com/ivch/dynasty/models/entities"
)

type Dictionaries struct {
	db *gorm.DB
}

func NewDictionaries(db *gorm.DB) *Dictionaries {
	return &Dictionaries{db: db}
}

func (r *Dictionaries) EntriesByBuilding(id uint) ([]*entities.Entry, error) {
	var reqs []*entities.Entry
	if err := r.db.Where("building_id = ?", id).Find(&reqs).Error; err != nil {
		return nil, err
	}
	return reqs, nil
}

func (r *Dictionaries) BuildingsList() ([]*entities.Building, error) {
	var reqs []*entities.Building
	if err := r.db.Find(&reqs).Error; err != nil {
		return nil, err
	}
	return reqs, nil
}
