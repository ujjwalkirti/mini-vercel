package project

import (
	"time"

	"reverse-proxy/internal/domain/deployment"
)

type Project struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	GitURL       string                 `json:"git_url"`
	Subdomain    string                 `json:"subdomain"`
	CustomDomain *string                `json:"custom_domain"`
	UserID       string                 `json:"user_id"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
	Deployments  []deployment.Deployment `json:"deployments,omitempty"`
}
