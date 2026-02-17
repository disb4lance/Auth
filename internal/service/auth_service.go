package service

import (
	model "auth-service/internal/domain/models"
	"auth-service/internal/pkg/pkg_dto"
	"auth-service/internal/service/dto"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hash, password string) bool
}

type TokenService interface {
	Generate(userID, email string) (*pkg_dto.TokenPair, error)
}

type RefreshTokenRepository interface {
	Create(token *model.RefreshToken) error
	GetByToken(token string) (*model.RefreshToken, error)
	Revoke(id uuid.UUID) error
}

type UserRepository interface {
	Create(user *model.User) error
	GetByID(id uuid.UUID) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
}

type AuthService struct {
	usersRepo  UserRepository
	tokensRepo RefreshTokenRepository
	hasher     PasswordHasher
	jwt        TokenService
	logger     *log.Logger
}

func NewAuthService(
	u UserRepository,
	t RefreshTokenRepository,
	h PasswordHasher,
	j TokenService,
	l *log.Logger,
) *AuthService {
	return &AuthService{
		usersRepo:  u,
		tokensRepo: t,
		hasher:     h,
		jwt:        j,
		logger:     l,
	}
}

func (s *AuthService) Register(email, password string) (*dto.TokenResponse, error) {
	hash, err := s.hasher.Hash(password)
	if err != nil {
		s.logger.Printf("hash error: %v", err)
		return nil, err
	}

	user := &model.User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: hash,
		CreatedAt:    time.Now().UTC(),
	}

	if err := s.usersRepo.Create(user); err != nil {
		s.logger.Printf("user create failed email=%s err=%v", email, err)
		return nil, err
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

	return &dto.TokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresAt:    tokens.ExpiresAt,
	}, nil
}

func (s *AuthService) Authenticate(creds dto.Credentials) (*dto.TokenResponse, error) {
	user, err := s.usersRepo.GetByEmail(creds.Email)
	if err != nil {
		return nil, err
	}

	if !s.hasher.Compare(user.PasswordHash, creds.Password) {
		s.logger.Printf("invalid credentials email=%s", creds.Email)
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

	return &dto.TokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresAt:    tokens.ExpiresAt,
	}, nil
}

func (s *AuthService) Refresh(refreshToken string) (*dto.TokenResponse, error) {
	rt, err := s.tokensRepo.GetByToken(refreshToken)
	if err != nil {
		return nil, err
	}

	if rt == nil || rt.IsRevoked || rt.ExpiresAt.Before(time.Now()) {
		s.logger.Printf("invalid refresh attempt token=%s", refreshToken)
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

	return &dto.TokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresAt:    tokens.ExpiresAt,
	}, nil
}
