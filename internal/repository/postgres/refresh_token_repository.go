package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	model "auth-service/internal/domain/models"
)

type RefreshTokenRepository struct {
	db *pgxpool.Pool
}

func NewRefreshTokenRepository(db *pgxpool.Pool) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

// Create - сохраняет новый токен
func (r *RefreshTokenRepository) Create(token *model.RefreshToken) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.db.Exec(ctx,
		`INSERT INTO refresh_tokens (id, token, expires_at, is_revoked, user_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		token.ID, token.Token, token.ExpiresAt, token.IsRevoked, token.UserID, token.CreatedAt,
	)
	return err
}

// GetByToken - ищет токен по значению
func (r *RefreshTokenRepository) GetByToken(tokenStr string) (*model.RefreshToken, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var token model.RefreshToken
	row := r.db.QueryRow(ctx,
		`SELECT id, token, expires_at, is_revoked, user_id, created_at 
		 FROM refresh_tokens 
		 WHERE token = $1`,
		tokenStr,
	)

	err := row.Scan(&token.ID, &token.Token, &token.ExpiresAt, &token.IsRevoked, &token.UserID, &token.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &token, nil
}

// Revoke - помечает токен как отозванный
func (r *RefreshTokenRepository) Revoke(id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.db.Exec(ctx,
		`UPDATE refresh_tokens SET is_revoked = true WHERE id = $1`,
		id,
	)
	return err
}
