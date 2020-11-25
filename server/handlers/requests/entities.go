package requests

import (
	"fmt"
	"time"

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
	History     pq.StringArray      `json:"-" gorm:"type:text[]"`
	ImagesURL   []map[string]string `json:"images" gorm:"-"`
	User        *users.User         `json:"user,omitempty"`
	CreatedAt   *time.Time
	DeletedAt   *time.Time
}

type UpdateRequest struct {
	ID          uint
	UserID      uint `gorm:"user_id"`
	Type        *string
	Time        *int64
	Description *string
	Status      *string
}

func (Request) TableName() string { return "requests" }

type RequestListFilter struct {
	DateFrom  *time.Time `json:"date_from,omitempty"`
	DateTo    *time.Time `json:"date_to,omitempty"`
	Type      string     `json:"type,omitempty" validate:"oneof=all taxi guest delivery noise complain"`
	Place     string     `json:"place,omitempty" validate:"oneof=all kpp"`
	Offset    uint       `json:"offset" validate:"min=0"`
	Limit     uint       `json:"limit" validate:"required,min=1"`
	UserID    uint       `json:"user_id,omitempty"`
	Apartment string     `json:"apartment,omitempty" validate:"omitempty,numeric"`
	Status    string     `json:"status,omitempty" validate:"oneof=all new closed"`
}

type Image struct {
	UserID    uint
	RequestID uint
	File      []byte
	URL       string
	Thumb     string
}

type HistoryRecord struct {
	Time   time.Time
	UserID uint
	Action string
}

func (h *HistoryRecord) String() string {
	return fmt.Sprintf("%s@%s@%d", h.Time.Format("2006-01-02 15:04"), h.Action, h.UserID)
}

// alter table requests
// add history text[] default '{}'::text[];
//
// alter table requests
// add created_at timestamp default CURRENT_TIMESTAMP not null;
//
// alter table requests
// add deleted_at timestamp default null;
