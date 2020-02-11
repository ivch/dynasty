package entities

const DefaultUserRole = 4

type User struct {
	ID         uint `gorm:"primary_key"`
	Building   Building
	Apartment  uint   `json:"apartment,omitempty"`
	Email      string `json:"email,omitempty"`
	Password   string `json:"password,omitempty"`
	Phone      string `json:"phone,omitempty"`
	FirstName  string `json:"first_name,omitempty"`
	LastName   string `json:"last_name,omitempty"`
	Role       uint   `json:"role,omitempty"`
	BuildingID int    `gorm:"building_id"`
}

func (User) TableName() string { return "users" }
