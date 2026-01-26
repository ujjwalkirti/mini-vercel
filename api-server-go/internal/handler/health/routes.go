package health

import (
	"github.com/go-chi/chi/v5"
)

func Routes() chi.Router {
	r := chi.NewRouter()
	h := NewHandler()

	r.Get("/", h.Check)

	return r
}
