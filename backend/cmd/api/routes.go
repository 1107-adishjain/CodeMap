package main

import (
	controller "codemap/backend/internal/controller"
	mw "codemap/backend/internal/middleware"
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (app *application) routes(db *sql.DB) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)

	// CORS settings
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // frontend origin
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by browsers
	}))
	r.Post("/api/v1/signup", controller.SignUp(db))
	r.Post("/api/v1/login", controller.Login(db))
	// API v1 routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(mw.Authenticate)
		r.Get("/healthcheck", app.healthCheckHandler)
		r.Post("/upload", app.uploadHandler)
		r.Post("/github", app.githubHandler)
		r.Post("/query", app.queryHandler)
	})
	return r
}