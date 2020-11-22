package transport

type errorResponse struct {
	Error     string `json:"error"`
	ErrorCode uint   `json:"error_code"`
}

type loginRequest struct {
	Phone    string `json:"phone" validate:"required,numeric"`
	Password string `json:"password" validate:"required"`
	// IP       net.IP `json:"ip" validate:"required"`
	// Ua       string `json:"ua" validate:"required"`
}

type loginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type authRefreshTokenRequest struct {
	Token string `validate:"required,uuid"`
}
