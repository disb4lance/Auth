package main

import (
	"auth-service/internal/handler"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

func (app *application) routes(handler *handler.AuthHandler) *chi.Mux {
	router := chi.NewRouter()

	router.Get("/auth/register", handler.RegisterHandler)
	router.Post("/auth/tokens", handler.AuthenticateHandler)
	router.Get("/auth/refresh", handler.RefreshHandler)
	router.Get("/swagger/", httpSwagger.WrapHandler)

	return router
}
