package router

import (
	"database/sql"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/ujjwalkirti/mini-vercel-api-server/internal/handler/deployment"
	"github.com/ujjwalkirti/mini-vercel-api-server/internal/handler/health"
	"github.com/ujjwalkirti/mini-vercel-api-server/internal/handler/project"
)

func New(db *sql.DB) *chi.Mux {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	// CORS middleware
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{os.Getenv("ALLOWED_ORIGINS")}, // Allow all origins for development
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Public routes (no auth required)
	r.Mount("/health", health.Routes())

	// Protected routes (auth required)
	r.Mount("/projects", project.Routes(db))
	r.Mount("/", deployment.Routes(db)) // Mounts /deploy and /deployments/* routes

	return r
}
