package app

import (
	"auth-service/internal/config"
	"auth-service/internal/handler"
	"auth-service/internal/infrastructure/jwt"
	"auth-service/internal/infrastructure/security"
	"auth-service/internal/repository/postgres"
	"auth-service/internal/service"
	transport "auth-service/internal/transport/http"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type App struct {
	Server *http.Server
	DB     *pgxpool.Pool
}

func New(cfg *config.Config, logger *log.Logger) (*App, error) {
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Name,
		cfg.DB.SSLMode,
	)

	db, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return nil, err
	}

	// repositories
	userRepo := postgres.NewUserRepository(db)
	tokenRepo := postgres.NewRefreshTokenRepository(db)

	// infrastructure
	hasher := security.NewBcryptHasher(bcrypt.DefaultCost)
	jwtSvc := jwt.NewJWTService(cfg.JWTSecret, 15*time.Minute)

	// services
	authSvc := service.NewAuthService(userRepo, tokenRepo, hasher, jwtSvc, logger)

	// handlers
	authHandler := handler.NewAuthHandler(authSvc)

	// router
	router := transport.NewRouter(authHandler)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return &App{
		Server: srv,
		DB:     db,
	}, nil
}
