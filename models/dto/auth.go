package dto

type AuthLoginRequest struct {
	Phone    string `json:"phone" validate:"required"`
	Password string `json:"password" validate:"required"`
	// IP       net.IP `json:"ip" validate:"required"`
	// Ua       string `json:"ua" validate:"required"`
}

type AuthLoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type AuthRefreshTokenRequest struct {
	Token string `validate:"required"`
}
