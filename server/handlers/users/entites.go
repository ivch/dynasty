package users

const defaultUserRole = 4

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
