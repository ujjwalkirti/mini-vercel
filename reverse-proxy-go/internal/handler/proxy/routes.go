package proxy

import (
	"github.com/go-chi/chi/v5"
)

func (h *Handler) RegisterRoutes(r *chi.Mux) {
	// Catch-all route for proxying
	r.HandleFunc("/*", h.ProxyRequest)
}
