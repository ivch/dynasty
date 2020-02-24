package dto

import (
	"github.com/ivch/dynasty/models/entities"
)

type RequestByID struct {
	UserID uint `validate:"required"`
	ID     uint `validate:"required"`
}

type RequestByIDResponse struct {
	ID          uint   `json:"id"`
	Type        string `json:"type"`
	UserID      uint   `json:"user_id" gorm:"user_id"`
	Time        int64  `json:"time"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type RequestCreateRequest struct {
	Type        string `json:"type" validate:"oneof=taxi guest delivery noise complain"`
	Time        int64  `json:"time" validate:"required"`
	UserID      uint   `json:"user_id" validate:"required"`
	Description string `json:"description"`
}

type RequestCreateResponse struct {
	ID uint `json:"id"`
}

type RequestMyRequest struct {
	UserID uint `json:"user_id"`
	Offset uint `json:"offset"`
	Limit  uint `json:"limit"`
}

type RequestMyResponse struct {
	Data []*entities.Request `json:"data"`
}

type RequestUpdateRequest struct {
	ID          uint
	UserID      uint    `gorm:"-"`
	Type        *string `json:"type,omitempty"`
	Time        *int64  `json:"time,omitempty"`
	Description *string `json:"description,omitempty"`
	Status      *string `json:"status,omitempty"`
}

type RequestListFilterRequest struct {
	// DateFrom   *time.Time `json:"date_from,omitepmpty"`
	// DateTo     *time.Time `json:"date_to,omitempty"`
	Type      string `json:"type,omitempty" validate:"oneof=all taxi guest delivery noise complain"`
	Offset    uint   `json:"offset" validate:"min=0"`
	Limit     uint   `json:"limit" validate:"required,min=1"`
	UserID    uint   `json:"user_id,omitempty"`
	Apartment string `json:"appartment,omitempty" validate:"omitempty,numeric"`
	Status    string `json:"status,omitempty" validate:"oneof=all new closed"`
}

type RequestForGuard struct {
	ID          uint   `json:"id"`
	UserID      uint   `json:"user_id" gorm:"-"`
	Type        string `json:"type"`
	Time        int64  `json:"time"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status"`
	UserName    string `json:"user_name"`
	Phone       string `json:"phone"`
	Address     string `json:"address"`
	Apartment   uint   `json:"apartment"`
}

type GuardUpdateRequest struct {
	ID     uint   `json:"id" validate:"required"`
	Status string `json:"status" validate:"oneof=new closed"`
}
