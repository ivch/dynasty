package entities

const DefaultUserRole = 4

type User struct {
	ID         uint `gorm:"primary_key"`
	Building   Building
	Entry      Entry
	Apartment  uint   `json:"apartment,omitempty"`
	Email      string `json:"email,omitempty"`
	Password   string `json:"password,omitempty"`
	Phone      string `json:"phone,omitempty"`
	FirstName  string `json:"first_name,omitempty"`
	LastName   string `json:"last_name,omitempty"`
	Role       uint   `json:"role,omitempty"`
	BuildingID uint   `gorm:"building_id"`
	EntryID    uint   `gorm:"entry_id"`
	Active     bool   `json:"active" gorm:"active"`
	RegCode    string `json:"-"`
	ParentID   *uint  `json:"-"`
}

func (User) TableName() string { return "users" }
