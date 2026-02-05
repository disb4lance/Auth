package main

import (
	_ "auth-service/docs"
	"auth-service/internal/handler"
	"auth-service/internal/infrastructure/jwt"
	"auth-service/internal/infrastructure/security"
	"auth-service/internal/repository/postgres"
	"auth-service/internal/service"
	"context"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	httpSwagger "github.com/swaggo/http-swagger"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	ctx := context.Background()
	db, err := pgxpool.New(ctx, "postgres://myuser:mypassword@localhost:5432/mydatabase?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	// Репозитории
	userRepo := postgres.NewUserRepository(db)
	tokenRepo := postgres.NewRefreshTokenRepository(db)
	hasher := security.NewBcryptHasher(bcrypt.DefaultCost)
	jwtSvc := jwt.NewJWTService(
		"super-secret-key", // secret
		15*time.Minute,     // access TTL
	)

	// Сервис
	authSvc := service.NewAuthService(userRepo, tokenRepo, hasher, jwtSvc)

	// Хендлер
	authHandler := handler.NewAuthHandler(authSvc)

	mux := http.NewServeMux()
	mux.HandleFunc("/auth/register", authHandler.RegisterHandler)
	mux.HandleFunc("/auth/tokens", authHandler.AuthenticateHandler)
	mux.HandleFunc("/auth/refresh", authHandler.RefreshHandler)
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
