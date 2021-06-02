package users

import "time"

const (
	defaultUserRole    = 4
	predefinedUserRole = 5
)

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

type UserUpdate struct {
	ID          uint    `gorm:"primary_key"`
	Email       *string `json:"email,omitempty"`
	Phone       *string `json:"phone,omitempty" gorm:"phone"`
	Password    *string `json:"password,omitempty"`
	NewPassword *string `json:"new_password,omitempty"`
	FirstName   *string `json:"first_name,omitempty"`
	LastName    *string `json:"last_name,omitempty"`
	Active      *bool   `json:"active" gorm:"active"`
	Role        *uint   `json:"role,omitempty" gorm:"role"`
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

type PasswordRecovery struct {
	ID        uint `gorm:"primary_key"`
	UserID    uint `gorm:"user_id"`
	Code      string
	CreatedAt *time.Time
	Active    bool
}

func (PasswordRecovery) TableName() string { return "password_recovery" }
