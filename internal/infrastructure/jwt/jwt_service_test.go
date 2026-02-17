package jwt

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJWTService_Generate(t *testing.T) {
	secret := "test-secret-key"
	ttl := 15 * time.Minute
	service := NewJWTService(secret, ttl)

	userID := uuid.New().String()
	email := "test@example.com"

	tokenPair, err := service.Generate(userID, email)

	require.NoError(t, err)
	require.NotNil(t, tokenPair)

	assert.NotEmpty(t, tokenPair.AccessToken)
	assert.NotEmpty(t, tokenPair.RefreshToken)
	assert.True(t, tokenPair.ExpiresAt.After(time.Now()))
}

func TestJWTService_Generate_InvalidInput(t *testing.T) {
	service := NewJWTService("secret", time.Minute)

	tests := []struct {
		name    string
		userID  string
		email   string
		wantErr bool
	}{
		{
			name:    "пустой userID",
			userID:  "",
			email:   "test@test.com",
			wantErr: false,
		},
		{
			name:    "пустой email",
			userID:  uuid.New().String(),
			email:   "",
			wantErr: false,
		},
		{
			name:    "оба поля пустые",
			userID:  "",
			email:   "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pair, err := service.Generate(tt.userID, tt.email)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, pair)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, pair)
			}
		})
	}
}

func TestJWTService_RefreshTokenUniqueness(t *testing.T) {
	service := NewJWTService("secret", time.Minute)

	pair1, _ := service.Generate("user1", "test@test.com")
	pair2, _ := service.Generate("user1", "test@test.com")

	assert.NotEqual(t, pair1.RefreshToken, pair2.RefreshToken)
}

func BenchmarkJWTService_Generate(b *testing.B) {
	service := NewJWTService("secret", time.Hour)
	userID := uuid.New().String()
	email := "test@test.com"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.Generate(userID, email)
		if err != nil {
			b.Fatal(err)
		}
	}
}
