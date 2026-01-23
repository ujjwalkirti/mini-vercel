package project

import (
	"database/sql"

	"github.com/go-chi/chi/v5"
)

func Routes(db *sql.DB) chi.Router {
	r := chi.NewRouter()
	h := NewHandler(
		db,
	)

	r.Get("/", h.GetProjects)
	r.Get("/{id}", h.GetProject)
	r.Post("/", h.CreateProject)
	r.Put("/{id}", h.UpdateProject)
	r.Delete("/{id}", h.DeleteProject)
	return r
}
