package main

import (
	_ "auth-service/docs"
	"auth-service/internal/handler"
	"auth-service/internal/repository/postgres"
	"auth-service/internal/service"
	"context"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	httpSwagger "github.com/swaggo/http-swagger"
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

	// Сервис
	authSvc := service.NewAuthService(userRepo, tokenRepo)

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
