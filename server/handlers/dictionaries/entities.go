package dictionaries

type Building struct {
	ID      uint   `gorm:"primary_key" json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

func (Building) TableName() string { return "buildings" }

type Entry struct {
	ID         uint   `json:"id" gorm:"primary_key"`
	Name       string `json:"name"`
	BuildingID uint   `json:"building_id"`
}

func (Entry) TableName() string { return "entries" }
