package deployment

import (
	"database/sql"

	"github.com/go-chi/chi/v5"
	"github.com/ujjwalkirti/mini-vercel-api-server/internal/middleware"
	repository "github.com/ujjwalkirti/mini-vercel-api-server/internal/repository/deployment"
)

func Routes(db *sql.DB) chi.Router {
	r := chi.NewRouter()

	// Apply auth middleware to all deployment routes
	r.Use(middleware.AuthMiddleware)

	repository := repository.New(db)
	h := NewHandler(repository)

	// GET /projects/:projectId/deployments - Get all deployments for a project
	r.Get("/projects/{projectId}/deployments", h.GetDeploymentsByProject)

	// GET /deployments/:id - Get specific deployment
	r.Get("/deployments/{id}", h.GetDeployment)

	// POST /deploy - Create new deployment
	r.Post("/deploy", h.CreateDeployment)

	// GET /deployments/:id/logs - Get deployment logs
	r.Get("/deployments/{id}/logs", h.GetDeploymentLogs)

	return r
}
