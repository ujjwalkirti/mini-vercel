package project

import (
	"database/sql"

	"github.com/go-chi/chi/v5"
	"github.com/ujjwalkirti/mini-vercel-api-server/internal/auth"
	"github.com/ujjwalkirti/mini-vercel-api-server/internal/middleware"
	deploymentRepo "github.com/ujjwalkirti/mini-vercel-api-server/internal/repository/deployment"
	repo "github.com/ujjwalkirti/mini-vercel-api-server/internal/repository/project"
)

func Routes(db *sql.DB, jwks *auth.JWKSCache) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.AuthMiddleware(jwks))

	repository := repo.New(db)
	deploymentRepository := deploymentRepo.New(db)
	h := NewHandler(repository, deploymentRepository)

	r.Get("/", h.GetProjects)
	r.Get("/{id}", h.GetProject)
	r.Post("/", h.CreateProject)
	r.Put("/{id}", h.UpdateProject)
	r.Delete("/{id}", h.DeleteProject)

	return r
}
