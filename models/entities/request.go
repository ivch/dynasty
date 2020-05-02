package entities

import "github.com/lib/pq"

type Request struct {
	ID          uint           `json:"id"`
	Type        string         `json:"type"`
	UserID      uint           `json:"user_id" gorm:"user_id"`
	Time        int64          `json:"time"`
	Description string         `json:"description"`
	Status      string         `json:"status"`
	Images      pq.StringArray `json:"-" gorm:"type:text[]"`
	User        *User          `json:"user,omitempty"`
}

func (Request) TableName() string { return "requests" }
