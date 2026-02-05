package repository

import (
	"github.com/google/uuid"

	model "auth-service/internal/domain/models"
)

type UserRepository interface {
	Create(user *model.User) error
	GetByID(id uuid.UUID) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
}
