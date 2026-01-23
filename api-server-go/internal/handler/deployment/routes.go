package deployment

import (
	"database/sql"

	"github.com/go-chi/chi/v5"
)

func Routes(db *sql.DB) chi.Router {
	r := chi.NewRouter()
	h := NewHandler(db)

	r.Get("/", h.GetDeployments)
	r.Get("/{id}", h.GetDeployment)
	r.Post("/", h.CreateDeployment)
	return r
}
