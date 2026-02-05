// internal/service/token_service.go
package service

import "time"

type TokenPair struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

type TokenService interface {
	Generate(userID, email string) (*TokenPair, error)
}
