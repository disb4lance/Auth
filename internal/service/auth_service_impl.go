package service

import (
	"errors"
	"time"

	"github.com/google/uuid"

	model "auth-service/internal/domain/models"
	"auth-service/internal/repository"
)

type authService struct {
	usersRepo  repository.UserRepository
	tokensRepo repository.RefreshTokenRepository
}

// Конструктор
func NewAuthService(u repository.UserRepository, t repository.RefreshTokenRepository) AuthService {
	return &authService{
		usersRepo:  u,
		tokensRepo: t,
	}
}

// Register — создаёт пользователя
func (s *authService) Register(email, password string) (*UserDTO, error) {
	// TODO: bcrypt hash password
	user := &model.User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: password, // временно plain-text
		CreatedAt:    time.Now().UTC(),
	}

	err := s.usersRepo.Create(user)
	if err != nil {
		return nil, err
	}

	return &UserDTO{
		ID:    user.ID.String(),
		Email: user.Email,
	}, nil
}

// Authenticate — проверяет credentials и возвращает токены
func (s *authService) Authenticate(creds Credentials) (*AuthenticatedUser, error) {
	user, err := s.usersRepo.GetByEmail(creds.Email)
	if err != nil {
		return nil, err
	}
	if user == nil || user.PasswordHash != creds.Password {
		return nil, errors.New("invalid credentials")
	}

	// TODO: здесь будет генерация JWT
	tokens := TokenPair{
		AccessToken:  "access-token-placeholder",
		RefreshToken: "refresh-token-placeholder",
	}

	// Создаём refresh token в БД
	rt := &model.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     tokens.RefreshToken,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		CreatedAt: time.Now().UTC(),
		IsRevoked: false,
	}

	err = s.tokensRepo.Create(rt)
	if err != nil {
		return nil, err
	}

	return &AuthenticatedUser{
		User: UserDTO{
			ID:    user.ID.String(),
			Email: user.Email,
		},
		Token: tokens,
	}, nil
}

// Refresh — обновляет токены по refresh token
func (s *authService) Refresh(refreshToken string) (*AuthenticatedUser, error) {
	rt, err := s.tokensRepo.GetByToken(refreshToken)
	if err != nil {
		return nil, err
	}
	if rt == nil || rt.IsRevoked || rt.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("invalid refresh token")
	}

	user, err := s.usersRepo.GetByID(rt.UserID)
	if err != nil || user == nil {
		return nil, errors.New("user not found")
	}

	// TODO: создать новые JWT
	newTokens := TokenPair{
		AccessToken:  "new-access-token-placeholder",
		RefreshToken: "new-refresh-token-placeholder",
	}

	// сохраняем новый refresh token
	newRT := &model.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     newTokens.RefreshToken,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		CreatedAt: time.Now().UTC(),
		IsRevoked: false,
	}

	err = s.tokensRepo.Create(newRT)
	if err != nil {
		return nil, err
	}

	// Отзываем старый токен
	err = s.tokensRepo.Revoke(rt.ID)
	if err != nil {
		return nil, err
	}

	return &AuthenticatedUser{
		User: UserDTO{
			ID:    user.ID.String(),
			Email: user.Email,
		},
		Token: newTokens,
	}, nil
}
