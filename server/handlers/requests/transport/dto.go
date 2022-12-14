package transport

import (
	"time"

	"github.com/microcosm-cc/bluemonday"
)

type errorResponse struct {
	Error     string `json:"error"`
	ErrorCode uint   `json:"error_code"`
	Ru        string `json:"ru"`
	Ua        string `json:"ua"`
}

type RequestCreateRequest struct {
	Type        string `json:"type"`
	Time        int64  `json:"time"`
	UserID      uint   `json:"user_id"`
	Description string `json:"description"`
}

func (r *RequestCreateRequest) Sanitize(p *bluemonday.Policy) {
	r.Description = p.Sanitize(r.Description)
}

type RequestCreateResponse struct {
	ID uint `json:"id"`
}

type RequestUpdateRequest struct {
	ID          uint
	Type        *string `json:"type,omitempty"`
	Time        *int64  `json:"time,omitempty"`
	Description *string `json:"description,omitempty"`
	Status      *string `json:"status,omitempty"`
}

func (r *RequestUpdateRequest) Sanitize(p *bluemonday.Policy) {
	if r.Description == nil {
		return
	}
	desc := p.Sanitize(*r.Description)
	r.Description = &desc
}

type ListByUserResponse struct {
	Data []*RequestByIDResponse `json:"data"`
}

type RequestByIDResponse struct {
	ID          uint                `json:"id"`
	Type        string              `json:"type"`
	UserID      uint                `json:"user_id"`
	Time        int64               `json:"time"`
	Description string              `json:"description"`
	Status      string              `json:"status"`
	Images      []map[string]string `json:"images,omitempty"`
	CreatedAt   *time.Time          `json:"created_at,omitempty"`
}

type UploadImageResponse struct {
	Img   string `json:"img"`
	Thumb string `json:"thumb"`
}

type DeleteImageRequest struct {
	Filepath string `json:"filepath"`
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
	CreatedAt   *time.Time          `json:"created_at,omitempty"`
}

type RequestGuardListResponse struct {
	Data  []*RequestForGuard `json:"data"`
	Count int                `json:"count"`
}

type GuardUpdateRequest struct {
	Status string `json:"status"`
}
