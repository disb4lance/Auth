package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port      int
	Env       string
	JWTSecret string
	DB        DBConfig
}

type DBConfig struct {
	User     string
	Password string
	Name     string
	Host     string
	Port     string
	SSLMode  string
}

func Load() (*Config, error) {
	port, err := strconv.Atoi(getEnv("PORT", "8080"))
	if err != nil {
		return nil, fmt.Errorf("invalid PORT value: %w", err)
	}

	cfg := &Config{
		Port:      port,
		Env:       getEnv("ENV", "development"),
		JWTSecret: getEnv("JWT_SECRET", "k3+Ta<8$9Vc4,*<My|B.=uF6[Tl>S_N)L$+)EHB#XRu"),
		DB: DBConfig{
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			Name:     getEnv("DB_NAME", "authdb"),
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5433"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
	}

	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
