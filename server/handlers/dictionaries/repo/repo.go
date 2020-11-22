package repository

import (
	"github.com/jinzhu/gorm"

	"github.com/ivch/dynasty/server/handlers/dictionaries"
)

type Dictionaries struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Dictionaries {
	return &Dictionaries{db: db}
}

func (r *Dictionaries) EntriesByBuilding(id uint) ([]*dictionaries.Entry, error) {
	var reqs []*dictionaries.Entry
	if err := r.db.Where("building_id = ?", id).Find(&reqs).Error; err != nil {
		return nil, err
	}
	return reqs, nil
}

func (r *Dictionaries) BuildingsList() ([]*dictionaries.Building, error) {
	var reqs []*dictionaries.Building
	if err := r.db.Find(&reqs).Error; err != nil {
		return nil, err
	}
	return reqs, nil
}
