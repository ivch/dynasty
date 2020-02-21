package dto

import (
	"time"

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
	Type        string `json:"type" validate:"required,oneof=taxi delivery guest"`
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

type GuardListRequest struct {
	Status string `json:"status" validate:"oneof=all new closed"`
	Offset uint   `json:"offset" validate:"min=0"`
	Limit  uint   `json:"limit" validate:"required"`
}

type RequestListFilter struct {
	DateFrom   *time.Time `json:"date_from,omitepmpty"`
	DateTo     *time.Time `json:"date_to,omitempty"`
	Offset     uint       `json:"offset" validate:"min=0"`
	Limit      uint       `json:"limit" validate:"required"`
	UserID     *uint      `json:"user_id,omitempty"`
	Appartment *uint      `json:"appartment,omitempty"`
	Status     string     `json:"status" validate:"oneof=all new closed"`
	Guard      bool
}

// type RequestForGuard struct {
// 	ID          uint   `json:"id"`
// 	Type        string `json:"type"`
// 	Time        int64  `json:"time"`
// 	Description string `json:"description"`
// 	Status      string `json:"status"`
// 	// Appartment
// 	// Building   Building
// 	// Apartment  uint   `json:"apartment,omitempty"`
// 	// Email      string `json:"email,omitempty"`
// 	// Password   string `json:"password,omitempty"`
// 	// Phone      string `json:"phone,omitempty"`
// 	// FirstName  string `json:"first_name,omitempty"`
// 	// LastName   string `json:"last_name,omitempty"`
// 	// Role       uint   `json:"role,omitempty"`
//
// }

type GuardUpdateRequest struct {
	ID     uint   `json:"id" validate:"required"`
	Status string `json:"status" validate:"oneof=new closed"`
}
