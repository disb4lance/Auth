package repository

import (
	model "auth-service/internal/domain/models"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(user *model.User) error
	GetByID(id uuid.UUID) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	Delete(id uuid.UUID) error
}
