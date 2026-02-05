package service

import model "auth-service/internal/domain/models"

type AuthService interface {
	Authenticate(credentials model.Credentials) (*model.AuthenticatedUser, error)
	Register(credentials model.Credentials) error
	Refresh(token model.Token) (*model.AuthenticatedUser, error)
}
