package dto

import "github.com/ivch/dynasty/models/entities"

type UserRegisterRequest struct {
	Password   string `json:"password,omitempty" validate:"required,min=6"`
	Phone      string `json:"phone,omitempty" validate:"required"`
	FirstName  string `json:"first_name,omitempty" validate:"required"`
	LastName   string `json:"last_name,omitempty" validate:"required"`
	BuildingID int    `json:"building_id,omitempty" validate:"required"`
	Apartment  uint   `json:"apartment,omitempty" validate:"required"`
	Email      string `json:"email,omitempty" validate:"email"`
	Code       string `json:"code"`
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
