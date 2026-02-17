package main

import (
	_ "auth-service/docs"
	"auth-service/internal/app"
	"auth-service/internal/config"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("no .env file found")
	}

	logger := log.New(os.Stdout, "[auth-service] ", log.LstdFlags)

	cfg, err := config.Load()
	if err != nil {
		logger.Fatal(err)
	}

	application, err := app.New(cfg, logger)
	if err != nil {
		logger.Fatal(err)
	}

	go func() {
		logger.Printf("starting %s server on %s", cfg.Env, application.Server.Addr)
		if err := application.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := application.Server.Shutdown(ctx); err != nil {
		logger.Fatal("server forced to shutdown:", err)
	}

	application.DB.Close()

	logger.Println("server exiting")
}
