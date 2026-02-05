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
	hasher     PasswordHasher
	jwt        TokenService
}

type TokenResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

func NewAuthService(
	u repository.UserRepository,
	t repository.RefreshTokenRepository,
	h PasswordHasher,
	j TokenService,
) AuthService {
	return &authService{
		usersRepo:  u,
		tokensRepo: t,
		hasher:     h,
		jwt:        j,
	}
}

func (s *authService) Register(email, password string) (*UserDTO, error) {
	hash, err := s.hasher.Hash(password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: hash,
		CreatedAt:    time.Now().UTC(),
	}

	if err := s.usersRepo.Create(user); err != nil {
		return nil, err
	}

	return &UserDTO{
		ID:    user.ID.String(),
		Email: user.Email,
	}, nil
}

func (s *authService) Authenticate(creds Credentials) (*TokenResponse, error) {
	user, err := s.usersRepo.GetByEmail(creds.Email)
	if err != nil {
		return nil, err
	}

	if !s.hasher.Compare(user.PasswordHash, creds.Password) {
		return nil, errors.New("invalid credentials")
	}

	tokens, err := s.jwt.Generate(
		user.ID.String(),
		user.Email,
	)
	if err != nil {
		return nil, err
	}

	rt := &model.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     tokens.RefreshToken,
		ExpiresAt: tokens.ExpiresAt,
		CreatedAt: time.Now().UTC(),
		IsRevoked: false,
	}

	if err := s.tokensRepo.Create(rt); err != nil {
		return nil, err
	}

	return &TokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresAt:    tokens.ExpiresAt,
	}, nil
}

func (s *authService) Refresh(refreshToken string) (*TokenResponse, error) {
	// 1. ищем refresh token
	rt, err := s.tokensRepo.GetByToken(refreshToken)
	if err != nil {
		return nil, err
	}

	if rt == nil || rt.IsRevoked || rt.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("invalid refresh token")
	}

	user, err := s.usersRepo.GetByID(rt.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	tokens, err := s.jwt.Generate(
		user.ID.String(),
		user.Email,
	)
	if err != nil {
		return nil, err
	}

	newRT := &model.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     tokens.RefreshToken,
		ExpiresAt: tokens.ExpiresAt,
		CreatedAt: time.Now().UTC(),
		IsRevoked: false,
	}

	if err := s.tokensRepo.Create(newRT); err != nil {
		return nil, err
	}

	if err := s.tokensRepo.Revoke(rt.ID); err != nil {
		return nil, err
	}

	return &TokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresAt:    tokens.ExpiresAt,
	}, nil
}
