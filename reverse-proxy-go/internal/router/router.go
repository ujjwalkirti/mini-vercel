package router

import (
	"database/sql"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	proxyHandler "reverse-proxy/internal/handler/proxy"
	customMiddleware "reverse-proxy/internal/middleware"
	"reverse-proxy/internal/repository/project"
)

func New(db *sql.DB, r2PublicURL string) *chi.Mux {
	r := chi.NewRouter()

	// Middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(customMiddleware.Logger)
	r.Use(middleware.Recoverer)

	// CORS configuration
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Initialize repository
	projectRepo := project.NewRepository(db)

	// Initialize handlers
	proxy := proxyHandler.NewHandler(projectRepo, r2PublicURL)

	// Register routes
	proxy.RegisterRoutes(r)

	return r
}
