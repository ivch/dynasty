package requests

import (
	"github.com/lib/pq"

	"github.com/ivch/dynasty/server/handlers/users"
)

type Request struct {
	ID          uint                `json:"id"`
	Type        string              `json:"type"`
	UserID      uint                `json:"user_id" gorm:"user_id"`
	Time        int64               `json:"time"`
	Description string              `json:"description"`
	Status      string              `json:"status"`
	Images      pq.StringArray      `json:"-" gorm:"type:text[]"`
	ImagesURL   []map[string]string `json:"images" gorm:"-"`
	User        *users.User         `json:"user,omitempty"`
}

func (Request) TableName() string { return "requests" }

type RequestListFilter struct {
	// DateFrom   *time.Time `json:"date_from,omitepmpty"`
	// DateTo     *time.Time `json:"date_to,omitempty"`
	Type      string `json:"type,omitempty" validate:"oneof=all taxi guest delivery noise complain"`
	Place     string `json:"place,omitempty" validate:"oneof=all kpp"`
	Offset    uint   `json:"offset" validate:"min=0"`
	Limit     uint   `json:"limit" validate:"required,min=1"`
	UserID    uint   `json:"user_id,omitempty"`
	Apartment string `json:"apartment,omitempty" validate:"omitempty,numeric"`
	Status    string `json:"status,omitempty" validate:"oneof=all new closed"`
}

type Image struct {
	UserID    uint
	RequestID uint
	File      []byte
	URL       string
	Thumb     string
}
