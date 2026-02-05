package repository

import (
	"github.com/google/uuid"

	model "auth-service/internal/domain/models"
)

type RefreshTokenRepository interface {
	Create(token *model.RefreshToken) error
	GetByToken(token string) (*model.RefreshToken, error)
	Revoke(id uuid.UUID) error
}
