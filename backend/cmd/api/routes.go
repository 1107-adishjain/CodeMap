package main

import (
	controller "codemap/backend/internal/controller"
	mw "codemap/backend/internal/middleware"
	"database/sql"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
)

func (app *application) routes(db *sql.DB) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(mw.SecureHeaders)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	r.Use(httprate.LimitByIP(100, 1*time.Minute)) //means allow 100 requests per minute per IP
	r.Use(middleware.Timeout(15 * time.Minute))
	r.Use(middleware.Compress(5))

	r.Post("/api/v1/signup", controller.SignUp(db))
	r.Post("/api/v1/login", controller.Login(db))
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(mw.Authenticate)
		r.Get("/healthcheck", app.healthCheckHandler)
		r.Post("/upload", app.uploadHandler)
		r.Post("/github", app.githubHandler)
		r.Post("/query", app.queryHandler)
	})

	return http.MaxBytesHandler(r, 300*1024*1024) 
}
