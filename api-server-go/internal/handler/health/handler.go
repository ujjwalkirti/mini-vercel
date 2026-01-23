package health

import (
	"net/http"
	"time"

	"github.com/ujjwalkirti/mini-vercel-api-server/internal/utils"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

func (h *Handler) Check(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status:    "ok",
		Timestamp: time.Now(),
	}

	utils.Success(w, response)
}
