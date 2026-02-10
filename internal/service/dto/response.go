package dto

import "time"

// возвращаем клиенту только безопасные поля
type UserDTO struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

// DTO для ответа после аутентификации
type AuthenticatedUser struct {
	User  UserDTO   `json:"user"`
	Token TokenPair `json:"token"`
}

type TokenPair struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

type TokenResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}
