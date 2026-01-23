package project

import (
	"database/sql"
	"net/http"
)

type Handler struct {
	db *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{
		db: db,
	}
}

func (h *Handler) GetProjects(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) GetProject(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) CreateProject(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) UpdateProject(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) DeleteProject(w http.ResponseWriter, r *http.Request) {}
