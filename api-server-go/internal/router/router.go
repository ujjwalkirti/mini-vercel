package router

import (
	"database/sql"

	"github.com/go-chi/chi/v5"
	"github.com/ujjwalkirti/mini-vercel-api-server/internal/handler/deployment"
	"github.com/ujjwalkirti/mini-vercel-api-server/internal/handler/project"
)

func New(db *sql.DB) *chi.Mux {
	r := chi.NewRouter()
	r.Mount("/projects", project.Routes(db))
	r.Mount("/deployments", deployment.Routes(db))
	return r
}
