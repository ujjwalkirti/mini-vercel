package project

import (
	"time"

	"github.com/ujjwalkirti/mini-vercel-api-server/internal/domain/deployment"
)

type Project struct {
	ID           string                  `json:"id"`
	Name         string                  `json:"name"`
	GitURL       string                  `json:"gitURL"`
	SubDomain    string                  `json:"subDomain"`
	CustomDomain *string                 `json:"customDomain"`
	UserID       string                  `json:"userId"`
	CreatedAt    time.Time               `json:"createdAt"`
	UpdatedAt    time.Time               `json:"updatedAt"`
	Deployments  []deployment.Deployment `json:"Deployment,omitempty"`
}
