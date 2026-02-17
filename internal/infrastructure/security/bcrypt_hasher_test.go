package security

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestBcryptHasher_Hash(t *testing.T) {
	t.Run("successful hash", func(t *testing.T) {
		hasher := NewBcryptHasher(bcrypt.DefaultCost)
		password := "mySecretPassword123!"

		hash, err := hasher.Hash(password)

		require.NoError(t, err)
		assert.NotEmpty(t, hash)
		assert.NotEqual(t, password, hash)
	})

	t.Run("hash with different costs", func(t *testing.T) {
		password := "testPassword"
		costs := []int{bcrypt.MinCost, bcrypt.DefaultCost, 12}

		for _, cost := range costs {
			hasher := NewBcryptHasher(cost)
			hash, err := hasher.Hash(password)

			require.NoError(t, err, "cost: %d", cost)
			assert.NotEmpty(t, hash)
		}
	})

	t.Run("same password produces different hashes", func(t *testing.T) {
		hasher := NewBcryptHasher(bcrypt.DefaultCost)
		password := "samePassword"

		hash1, _ := hasher.Hash(password)
		hash2, _ := hasher.Hash(password)

		assert.NotEqual(t, hash1, hash2)
	})

	t.Run("empty password", func(t *testing.T) {
		hasher := NewBcryptHasher(bcrypt.DefaultCost)
		hash, err := hasher.Hash("")

		require.NoError(t, err)
		assert.NotEmpty(t, hash)
	})
}

func TestBcryptHasher_Compare(t *testing.T) {
	hasher := NewBcryptHasher(bcrypt.DefaultCost)

	t.Run("correct password", func(t *testing.T) {
		password := "myPassword123"
		hash, err := hasher.Hash(password)
		require.NoError(t, err)

		result := hasher.Compare(hash, password)
		assert.True(t, result)
	})

	t.Run("incorrect password", func(t *testing.T) {
		password := "myPassword123"
		wrongPassword := "wrongPassword"
		hash, err := hasher.Hash(password)
		require.NoError(t, err)

		result := hasher.Compare(hash, wrongPassword)
		assert.False(t, result)
	})

	t.Run("compare with invalid hash", func(t *testing.T) {
		hasher := NewBcryptHasher(bcrypt.DefaultCost)
		invalidHash := "invalidHashString"
		password := "anyPassword"

		result := hasher.Compare(invalidHash, password)
		assert.False(t, result)
	})

	t.Run("empty password compare", func(t *testing.T) {
		emptyPassword := ""
		hash, err := hasher.Hash(emptyPassword)
		require.NoError(t, err)

		result := hasher.Compare(hash, emptyPassword)
		assert.True(t, result)

		result = hasher.Compare(hash, "notEmpty")
		assert.False(t, result)
	})
}

func BenchmarkBcryptHasher_Hash(b *testing.B) {
	hasher := NewBcryptHasher(bcrypt.DefaultCost)
	password := "benchmarkPassword"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = hasher.Hash(password)
	}
}

func BenchmarkBcryptHasher_Compare(b *testing.B) {
	hasher := NewBcryptHasher(bcrypt.DefaultCost)
	password := "benchmarkPassword"
	hash, _ := hasher.Hash(password)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = hasher.Compare(hash, password)
	}
}

func TestBcryptHasher_Compare_TableDriven(t *testing.T) {
	hasher := NewBcryptHasher(bcrypt.DefaultCost)
	password := "testPassword123"
	hash, _ := hasher.Hash(password)

	tests := []struct {
		name     string
		hash     string
		password string
		want     bool
	}{
		{
			name:     "правильный пароль",
			hash:     hash,
			password: password,
			want:     true,
		},
		{
			name:     "неправильный пароль",
			hash:     hash,
			password: "wrong",
			want:     false,
		},
		{
			name:     "пустой пароль",
			hash:     hash,
			password: "",
			want:     false,
		},
		{
			name:     "невалидный хеш",
			hash:     "invalid",
			password: password,
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hasher.Compare(tt.hash, tt.password)
			assert.Equal(t, tt.want, got)
		})
	}
}
