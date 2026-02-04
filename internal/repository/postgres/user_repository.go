package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	model "auth-service/internal/domain/models"
	"auth-service/internal/repository"
)

type PostgresUserRepository struct {
	db *pgxpool.Pool
}

// Конструктор
func NewPostgresUserRepository(db *pgxpool.Pool) repository.UserRepository {
	return &PostgresUserRepository{db: db}
}

// Create вставляет нового пользователя
func (r *PostgresUserRepository) Create(user *model.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.db.Exec(ctx,
		"INSERT INTO users (id, email, password_hash, created_at) VALUES ($1, $2, $3, $4)",
		user.Id, user.Email, user.PasswordHash, user.CreatedAt,
	)
	return err
}

// GetByID возвращает пользователя по ID
func (r *PostgresUserRepository) GetByID(id uuid.UUID) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user model.User

	row := r.db.QueryRow(ctx,
		"SELECT id, email, password_hash, created_at FROM users WHERE id = $1",
		id,
	)

	err := row.Scan(&user.Id, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetByEmail возвращает пользователя по email
func (r *PostgresUserRepository) GetByEmail(email string) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user model.User

	row := r.db.QueryRow(ctx,
		"SELECT id, email, password_hash, created_at FROM users WHERE email = $1",
		email,
	)

	err := row.Scan(&user.Id, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Delete удаляет пользователя по ID
func (r *PostgresUserRepository) Delete(id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.db.Exec(ctx,
		"DELETE FROM users WHERE id = $1",
		id,
	)
	return err
}
