package router

import (
	"database/sql"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

	// Public routes (no auth required)
	r.Mount("/health", health.Routes())

	// Protected routes (auth required)
	r.Mount("/projects", project.Routes(db))
	r.Mount("/", deployment.Routes(db)) // Mounts /deploy and /deployments/* routes

	return r
}
