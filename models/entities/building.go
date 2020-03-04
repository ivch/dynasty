package entities

type Building struct {
	ID      uint   `gorm:"primary_key" json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

func (Building) TableName() string { return "buildings" }

type Entry struct {
	ID         uint `gorm:"primary_key" json:"id"`
	Name       string
	BuildingID uint
}

func (Entry) TableName() string { return "entries" }
