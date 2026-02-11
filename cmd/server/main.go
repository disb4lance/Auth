package main

import (
	"auth-service/internal/handler"
	"auth-service/internal/infrastructure/jwt"
	"auth-service/internal/infrastructure/security"
	"auth-service/internal/repository/postgres"
	"auth-service/internal/service"
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
}

type application struct {
	config config
	logger *log.Logger
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 8080, "API server port") // Изменено на 8080 для соответствия ListenAndServe
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	app := &application{
		config: cfg,
		logger: logger,
	}

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
		"super-secret-key",
		15*time.Minute,
	)

	// Сервис
	authSvc := service.NewAuthService(userRepo, tokenRepo, hasher, jwtSvc)

	// Хендлер
	authHandler := handler.NewAuthHandler(authSvc)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(authHandler),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Printf("starting %s server on %s", cfg.env, srv.Addr)
	logger.Fatal(srv.ListenAndServe())
}
