package transport

import (
	"github.com/microcosm-cc/bluemonday"

	"github.com/ivch/dynasty/server/handlers/users"
)

type UserByIDResponse struct {
	ID        uint            `json:"id"`
	Apartment uint            `json:"apartment"`
	FirstName string          `json:"first_name"`
	LastName  string          `json:"last_name"`
	Phone     string          `json:"phone"`
	Email     string          `json:"email"`
	Role      uint            `json:"role,omitempty"`
	Building  *users.Building `json:"building"`
	Entry     *users.Entry    `json:"entry,omitempty"`
	Active    bool            `json:"active" gorm:"active"`
}

type errorResponse struct {
	Error     string `json:"error"`
	ErrorCode uint   `json:"error_code"`
}

type userRegisterRequest struct {
	Password   string `json:"password,omitempty" validate:"required,min=6"`
	Phone      string `json:"phone,omitempty" validate:"required,numeric,min=12,max=13"`
	FirstName  string `json:"first_name,omitempty" validate:"required"`
	LastName   string `json:"last_name,omitempty" validate:"required"`
	BuildingID uint   `json:"building_id,omitempty" validate:"required"`
	EntryID    uint   `json:"entry_id" validate:"required,numeric"`
	Apartment  uint   `json:"apartment,omitempty" validate:"required"`
	Email      string `json:"email" validate:"required,email"`
	Code       string `json:"code"`
}

func (r *userRegisterRequest) Sanitize(p *bluemonday.Policy) {
	r.FirstName = p.Sanitize(r.FirstName)
	r.LastName = p.Sanitize(r.LastName)
	r.Code = p.Sanitize(r.Code)
}

type userRegisterResponse struct {
	ID    uint   `json:"id"`
	Phone string `json:"phone"`
}

type addFamilyMemberRequest struct {
	OwnerID uint
	Phone   string `json:"phone" validate:"required,numeric"`
}

type addFamilyMemberResponse struct {
	Code string `json:"code"`
}

type listFamilyMembersResponse struct {
	Data []*familyMember `json:"data"`
}

type familyMember struct {
	ID     uint   `json:"id"`
	Name   string `json:"name,omitempty"`
	Phone  string `json:"phone"`
	Code   string `json:"code"`
	Active bool   `json:"active"`
}

type UserUpdateRequest struct {
	ID                 uint    `gorm:"primary_key"`
	Email              *string `json:"email,omitempty"`
	Password           *string `json:"password,omitempty"`
	NewPassword        *string `json:"new_password,omitempty"`
	NewPasswordConfirm *string `json:"new_password_confirm,omitempty"`
	FirstName          *string `json:"first_name,omitempty"`
	LastName           *string `json:"last_name,omitempty"`
}

func (r *UserUpdateRequest) Sanitize(p *bluemonday.Policy) {
	if r.FirstName != nil {
		fname := p.Sanitize(*r.FirstName)
		r.FirstName = &fname
	}
	if r.LastName != nil {
		lname := p.Sanitize(*r.LastName)
		r.LastName = &lname
	}
}
