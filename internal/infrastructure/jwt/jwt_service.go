// jwt/jwt_service.go
package jwt

import (
	"auth-service/internal/pkg/pkg_dto"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTService struct {
	secret         []byte
	accessTokenTTL time.Duration
}

type AccessClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func NewJWTService(secret string, accessTTL time.Duration) *JWTService {
	return &JWTService{
		secret:         []byte(secret),
		accessTokenTTL: accessTTL,
	}
}

func (s *JWTService) Generate(userID, email string) (*pkg_dto.TokenPair, error) {
	now := time.Now()

	claims := AccessClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.secret)
	if err != nil {
		return nil, err
	}

	return &pkg_dto.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: uuid.NewString(),
		ExpiresAt:    claims.ExpiresAt.Time,
	}, nil
}
