package health

import (
	"net/http"
	"time"

	"github.com/ujjwalkirti/mini-vercel-api-server/internal/utils"
)

var startTime = time.Now()

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Uptime    float64   `json:"uptime"`
}

func (h *Handler) Check(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status:    "ok",
		Timestamp: time.Now(),
		Uptime:    time.Since(startTime).Seconds(),
	}

	utils.Success(w, response)
}
