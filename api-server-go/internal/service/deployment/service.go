package deployment

import (
	"context"

	repository "github.com/ujjwalkirti/mini-vercel-api-server/internal/repository/deployment"
)

type DeploymentService struct {
	repo *repository.Repository
}

func NewDeploymentService(repo *repository.Repository) *DeploymentService {
	return &DeploymentService{
		repo: repo,
	}
}

func (s *DeploymentService) MarkInProgress(ctx context.Context, deploymentID string) error {
	return s.repo.UpdateStatus(ctx, deploymentID, "IN_PROGRESS")
}

func (s *DeploymentService) MarkReady(ctx context.Context, deploymentID string) error {
	return s.repo.UpdateStatus(ctx, deploymentID, "READY")
}

func (s *DeploymentService) MarkFailed(ctx context.Context, deploymentID string) error {
	return s.repo.UpdateStatus(ctx, deploymentID, "FAIL")
}
