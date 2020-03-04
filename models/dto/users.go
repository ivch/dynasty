package dto

import (
	"github.com/microcosm-cc/bluemonday"

	"github.com/ivch/dynasty/models/entities"
)

type UserRegisterRequest struct {
	Password   string `json:"password,omitempty" validate:"required,min=6"`
	Phone      string `json:"phone,omitempty" validate:"required,numeric"`
	FirstName  string `json:"first_name,omitempty" validate:"required"`
	LastName   string `json:"last_name,omitempty" validate:"required"`
	BuildingID int    `json:"building_id,omitempty" validate:"required"`
	Apartment  uint   `json:"apartment,omitempty" validate:"required"`
	Email      string `json:"email" validate:"required,email"`
	Code       string `json:"code"`
}

func (r *UserRegisterRequest) Sanitize(p *bluemonday.Policy) {
	r.FirstName = p.Sanitize(r.FirstName)
	r.LastName = p.Sanitize(r.LastName)
	r.Code = p.Sanitize(r.Code)
}

type UserRegisterResponse struct {
	ID    uint   `json:"id"`
	Phone string `json:"phone"`
}

type UserAuthResponse struct {
	ID        uint   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      uint   `json:"role"`
}

type UserByIDResponse struct {
	ID        uint              `json:"id"`
	Apartment uint              `json:"apartment"`
	FirstName string            `json:"first_name"`
	LastName  string            `json:"last_name"`
	Phone     string            `json:"phone"`
	Email     string            `json:"email"`
	Role      uint              `json:"role"`
	Building  entities.Building `json:"building"`
}

type AddFamilyMemberRequest struct {
	OwnerID uint
	Phone   string `json:"phone" validate:"required,numeric"`
}

type AddFamilyMemberResponse struct {
	Code string `json:"code"`
}

type ListFamilyMembersResponse struct {
	Data []*FamilyMember `json:"data"`
}

type FamilyMember struct {
	ID     uint   `json:"id"`
	Name   string `json:"name,omitempty"`
	Phone  string `json:"phone"`
	Code   string `json:"code"`
	Active bool   `json:"active"`
}

type DeleteFamilyMemberRequest struct {
	OwnerID  uint
	MemberID uint
}
