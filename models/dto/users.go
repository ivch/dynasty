package dto

type UserAuthResponse struct {
	ID        uint   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      uint   `json:"role"`
	Active    bool   `json:"active"`
}
