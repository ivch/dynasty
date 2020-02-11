package dto

import "github.com/ivch/dynasty/models/entities"

type RequestByID struct {
	UserID uint `validate:"required"`
	ID     uint `validate:"required"`
}

type RequestByIDResponse struct {
	*entities.Request
}

type RequestCreateRequest struct {
	Type        string `json:"type" validate:"required"`
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
