package main

import (
	"log"

	"auth-service/internal/app"
	"auth-service/internal/config"
)

func main() {
	cfg := config.Load()

	application := app.New(cfg.HTTPPort)

	if err := application.Run(); err != nil {
		log.Fatal(err)
	}
}
