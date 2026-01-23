package deployment

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

func (h *Handler) GetDeployment(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) GetDeployments(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) CreateDeployment(w http.ResponseWriter, r *http.Request) {}
