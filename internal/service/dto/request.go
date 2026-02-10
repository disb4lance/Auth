package dto

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterRequest
type RegisterRequest struct {
	Email    string `json:"email" example:"test@mail.com"`
	Password string `json:"password" example:"qwerty123"`
}

// LoginRequest
type LoginRequest struct {
	Email    string `json:"email" example:"test@mail.com"`
	Password string `json:"password" example:"qwerty123"`
}

// RefreshRequest
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" example:"refresh.jwt.token"`
}
