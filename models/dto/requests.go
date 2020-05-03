package dto

import (
	"github.com/microcosm-cc/bluemonday"
)

type RequestByID struct {
	UserID uint `validate:"required"`
	ID     uint `validate:"required"`
}

type RequestByIDResponse struct {
	ID          uint                `json:"id"`
	Type        string              `json:"type"`
	UserID      uint                `json:"user_id" gorm:"user_id"`
	Time        int64               `json:"time"`
	Description string              `json:"description"`
	Status      string              `json:"status"`
	Images      []map[string]string `json:"images,omitempty"`
}

type RequestCreateRequest struct {
	Type        string `json:"type" validate:"oneof=taxi guest delivery noise complain"`
	Time        int64  `json:"time" validate:"required"`
	UserID      uint   `json:"user_id" validate:"required"`
	Description string `json:"description"`
}

func (r *RequestCreateRequest) Sanitize(p *bluemonday.Policy) {
	r.Description = p.Sanitize(r.Description)
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
	Data []*RequestByIDResponse `json:"data"`
}

type RequestUpdateRequest struct {
	ID          uint
	UserID      uint    `gorm:"-"`
	Type        *string `json:"type,omitempty" validate:"omitempty,oneof=all taxi guest delivery noise complain"`
	Time        *int64  `json:"time,omitempty"`
	Description *string `json:"description,omitempty"`
	Status      *string `json:"status,omitempty" validate:"omitempty,oneof=all new closed"`
}

func (r *RequestUpdateRequest) Sanitize(p *bluemonday.Policy) {
	if r.Description == nil {
		return
	}
	desc := p.Sanitize(*r.Description)
	r.Description = &desc
}

type RequestListFilterRequest struct {
	// DateFrom   *time.Time `json:"date_from,omitepmpty"`
	// DateTo     *time.Time `json:"date_to,omitempty"`
	Type      string `json:"type,omitempty" validate:"oneof=all taxi guest delivery noise complain"`
	Offset    uint   `json:"offset" validate:"min=0"`
	Limit     uint   `json:"limit" validate:"required,min=1"`
	UserID    uint   `json:"user_id,omitempty"`
	Apartment string `json:"apartment,omitempty" validate:"omitempty,numeric"`
	Status    string `json:"status,omitempty" validate:"oneof=all new closed"`
}

type RequestGuardListResponse struct {
	Data  []*RequestForGuard `json:"data"`
	Count int                `json:"count"`
}

type RequestForGuard struct {
	ID          uint                `json:"id"`
	UserID      uint                `json:"user_id" gorm:"-"`
	Type        string              `json:"type"`
	Time        int64               `json:"time"`
	Description string              `json:"description,omitempty"`
	Status      string              `json:"status"`
	UserName    string              `json:"user_name"`
	Phone       string              `json:"phone"`
	Address     string              `json:"address"`
	Apartment   uint                `json:"apartment"`
	Images      []map[string]string `json:"images,omitempty"`
}

type GuardUpdateRequest struct {
	ID     uint   `json:"id" validate:"required"`
	Status string `json:"status" validate:"oneof=new closed"`
}

type UploadImageRequest struct {
	UserID    uint
	RequestID uint
	File      []byte
}

type UploadImageResponse struct {
	Img   string `json:"img"`
	Thumb string `json:"thumb"`
}

type DeleteImageRequest struct {
	UserID    uint
	RequestID uint
	Filepath  string `json:"filepath" validate:"required"`
}
