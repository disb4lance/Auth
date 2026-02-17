// internal/service/auth_test.go
package service

import (
	models "auth-service/internal/domain/models"
	"auth-service/internal/pkg/pkg_dto"
	"auth-service/internal/service/dto"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAuthService_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем моки
	mockUserRepo := NewMockUserRepository(ctrl)
	mockTokenRepo := NewMockRefreshTokenRepository(ctrl)
	mockHasher := NewMockPasswordHasher(ctrl)
	mockJWT := NewMockTokenService(ctrl)

	service := NewAuthService(
		mockUserRepo,
		mockTokenRepo,
		mockHasher,
		mockJWT,
	)

	email := "test@example.com"
	password := "password123"
	hash := "hashed_password"

	t.Run("hash error", func(t *testing.T) {
		expectedErr := errors.New("hash error")

		mockHasher.EXPECT().
			Hash(password).
			Return("", expectedErr)

		resp, err := service.Register(email, password)

		assert.ErrorIs(t, err, expectedErr)
		assert.Nil(t, resp)
	})

	t.Run("user creation error", func(t *testing.T) {
		expectedErr := errors.New("db error")

		mockHasher.EXPECT().
			Hash(password).
			Return(hash, nil)

		mockUserRepo.EXPECT().
			Create(gomock.Any()).
			Return(expectedErr)

		resp, err := service.Register(email, password)

		assert.ErrorIs(t, err, expectedErr)
		assert.Nil(t, resp)
	})

	t.Run("jwt generation error", func(t *testing.T) {
		expectedErr := errors.New("jwt error")

		mockHasher.EXPECT().
			Hash(password).
			Return(hash, nil)

		mockUserRepo.EXPECT().
			Create(gomock.Any()).
			Return(nil)

		mockJWT.EXPECT().
			Generate(gomock.Any(), email).
			Return(nil, expectedErr)

		resp, err := service.Register(email, password)

		assert.ErrorIs(t, err, expectedErr)
		assert.Nil(t, resp)
	})
}

func TestAuthService_Authenticate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := NewMockUserRepository(ctrl)
	mockTokenRepo := NewMockRefreshTokenRepository(ctrl)
	mockHasher := NewMockPasswordHasher(ctrl)
	mockJWT := NewMockTokenService(ctrl)

	service := NewAuthService(
		mockUserRepo,
		mockTokenRepo,
		mockHasher,
		mockJWT,
	)

	creds := dto.Credentials{
		Email:    "test@example.com",
		Password: "password123",
	}
	userID := uuid.New()
	hash := "hashed_password"

	t.Run("successful authentication", func(t *testing.T) {
		user := &models.User{
			ID:           userID,
			Email:        creds.Email,
			PasswordHash: hash,
		}

		mockUserRepo.EXPECT().
			GetByEmail(creds.Email).
			Return(user, nil)

		mockHasher.EXPECT().
			Compare(hash, creds.Password).
			Return(true)

		tokens := &pkg_dto.TokenPair{
			AccessToken:  "access_token",
			RefreshToken: "refresh_token",
			ExpiresAt:    time.Now().Add(15 * time.Minute),
		}

		mockJWT.EXPECT().
			Generate(userID.String(), creds.Email).
			Return(tokens, nil)

		mockTokenRepo.EXPECT().
			Create(gomock.Any()).
			Return(nil)

		resp, err := service.Authenticate(creds)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("user not found", func(t *testing.T) {
		mockUserRepo.EXPECT().
			GetByEmail(creds.Email).
			Return(nil, errors.New("not found"))

		resp, err := service.Authenticate(creds)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("invalid password", func(t *testing.T) {
		user := &models.User{
			ID:           userID,
			Email:        creds.Email,
			PasswordHash: hash,
		}

		mockUserRepo.EXPECT().
			GetByEmail(creds.Email).
			Return(user, nil)

		mockHasher.EXPECT().
			Compare(hash, creds.Password).
			Return(false)

		resp, err := service.Authenticate(creds)

		assert.Error(t, err)
		assert.Equal(t, "invalid credentials", err.Error())
		assert.Nil(t, resp)
	})
}

func TestAuthService_Refresh(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := NewMockUserRepository(ctrl)
	mockTokenRepo := NewMockRefreshTokenRepository(ctrl)
	mockHasher := NewMockPasswordHasher(ctrl)
	mockJWT := NewMockTokenService(ctrl)

	service := NewAuthService(
		mockUserRepo,
		mockTokenRepo,
		mockHasher,
		mockJWT,
	)

	refreshToken := "valid_refresh_token"
	userID := uuid.New()
	rtID := uuid.New()
	now := time.Now()

	t.Run("successful refresh", func(t *testing.T) {
		oldRT := &models.RefreshToken{
			ID:        rtID,
			UserID:    userID,
			Token:     refreshToken,
			ExpiresAt: now.Add(24 * time.Hour),
			IsRevoked: false,
		}

		user := &models.User{
			ID:    userID,
			Email: "test@example.com",
		}

		mockTokenRepo.EXPECT().
			GetByToken(refreshToken).
			Return(oldRT, nil)

		mockUserRepo.EXPECT().
			GetByID(userID).
			Return(user, nil)

		tokens := &pkg_dto.TokenPair{
			AccessToken:  "new_access_token",
			RefreshToken: "new_refresh_token",
			ExpiresAt:    now.Add(15 * time.Minute),
		}

		mockJWT.EXPECT().
			Generate(userID.String(), user.Email).
			Return(tokens, nil)

		mockTokenRepo.EXPECT().
			Create(gomock.Any()).
			DoAndReturn(func(rt *models.RefreshToken) error {
				assert.Equal(t, userID, rt.UserID)
				assert.Equal(t, tokens.RefreshToken, rt.Token)
				return nil
			})

		mockTokenRepo.EXPECT().
			Revoke(rtID).
			Return(nil)

		resp, err := service.Refresh(refreshToken)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, tokens.AccessToken, resp.AccessToken)
		assert.Equal(t, tokens.RefreshToken, resp.RefreshToken)
	})

	t.Run("token not found", func(t *testing.T) {
		mockTokenRepo.EXPECT().
			GetByToken(refreshToken).
			Return(nil, errors.New("not found"))

		resp, err := service.Refresh(refreshToken)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("token is revoked", func(t *testing.T) {
		oldRT := &models.RefreshToken{
			ID:        rtID,
			UserID:    userID,
			Token:     refreshToken,
			ExpiresAt: now.Add(24 * time.Hour),
			IsRevoked: true,
		}

		mockTokenRepo.EXPECT().
			GetByToken(refreshToken).
			Return(oldRT, nil)

		resp, err := service.Refresh(refreshToken)

		assert.Error(t, err)
		assert.Equal(t, "invalid refresh token", err.Error())
		assert.Nil(t, resp)
	})

	t.Run("token expired", func(t *testing.T) {
		oldRT := &models.RefreshToken{
			ID:        rtID,
			UserID:    userID,
			Token:     refreshToken,
			ExpiresAt: now.Add(-1 * time.Hour), // просрочен
			IsRevoked: false,
		}

		mockTokenRepo.EXPECT().
			GetByToken(refreshToken).
			Return(oldRT, nil)

		resp, err := service.Refresh(refreshToken)

		assert.Error(t, err)
		assert.Equal(t, "invalid refresh token", err.Error())
		assert.Nil(t, resp)
	})

	t.Run("user not found after token check", func(t *testing.T) {
		oldRT := &models.RefreshToken{
			ID:        rtID,
			UserID:    userID,
			Token:     refreshToken,
			ExpiresAt: now.Add(24 * time.Hour),
			IsRevoked: false,
		}

		mockTokenRepo.EXPECT().
			GetByToken(refreshToken).
			Return(oldRT, nil)

		mockUserRepo.EXPECT().
			GetByID(userID).
			Return(nil, errors.New("user not found"))

		resp, err := service.Refresh(refreshToken)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}
