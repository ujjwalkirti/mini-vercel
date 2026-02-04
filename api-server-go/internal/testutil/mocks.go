package testutil

import (
	"context"

	"github.com/ujjwalkirti/mini-vercel-api-server/internal/domain/deployment"
	"github.com/ujjwalkirti/mini-vercel-api-server/internal/domain/project"
	"github.com/ujjwalkirti/mini-vercel-api-server/internal/service/ecs"
	"github.com/ujjwalkirti/mini-vercel-api-server/internal/service/logs"
)

// MockProjectRepository is a mock implementation of project repository
type MockProjectRepository struct {
	// TODO: Add fields for controlling mock behavior
	// TODO: Add fields for capturing method calls

	// Example fields:
	// ListByUserFunc func(ctx context.Context, userID string) ([]*project.Project, error)
	// GetByIDAndUserIDFunc func(ctx context.Context, id, userID string) (*project.Project, error)
	// CreateFunc func(ctx context.Context, project *project.Project) error
	// UpdateFunc func(ctx context.Context, project *project.Project) error
	// DeleteFunc func(ctx context.Context, id string) error
}

// ListByUser mocks listing projects by user
func (m *MockProjectRepository) ListByUser(ctx context.Context, userID string) ([]*project.Project, error) {
	// TODO: Implement - call ListByUserFunc if set
	// TODO: Return default behavior or error
	return nil, nil
}

// GetByIDAndUserID mocks getting a project by ID and user ID
func (m *MockProjectRepository) GetByIDAndUserID(ctx context.Context, id, userID string) (*project.Project, error) {
	// TODO: Implement - call GetByIDAndUserIDFunc if set
	// TODO: Return default behavior or error
	return nil, nil
}

// Create mocks creating a project
func (m *MockProjectRepository) Create(ctx context.Context, project *project.Project) error {
	// TODO: Implement - call CreateFunc if set
	// TODO: Return default behavior or error
	return nil
}

// Update mocks updating a project
func (m *MockProjectRepository) Update(ctx context.Context, project *project.Project) error {
	// TODO: Implement - call UpdateFunc if set
	// TODO: Return default behavior or error
	return nil
}

// Delete mocks deleting a project
func (m *MockProjectRepository) Delete(ctx context.Context, id string) error {
	// TODO: Implement - call DeleteFunc if set
	// TODO: Return default behavior or error
	return nil
}

// MockDeploymentRepository is a mock implementation of deployment repository
type MockDeploymentRepository struct {
	// TODO: Add fields for controlling mock behavior
	// TODO: Add fields for capturing method calls

	// Example fields:
	// GetByProjectIDFunc func(ctx context.Context, projectID, userID string) ([]*deployment.Deployment, error)
	// GetByIDWithProjectFunc func(ctx context.Context, id, userID string) (*deployment.DeploymentWithProject, error)
	// CreateFunc func(ctx context.Context, deployment *deployment.Deployment) (*deployment.Deployment, error)
	// DeleteByProjectIDFunc func(ctx context.Context, projectID string) error
}

// GetByProjectID mocks getting deployments by project ID
func (m *MockDeploymentRepository) GetByProjectID(ctx context.Context, projectID, userID string) ([]*deployment.Deployment, error) {
	// TODO: Implement - call GetByProjectIDFunc if set
	// TODO: Return default behavior or error
	return nil, nil
}

// GetByIDWithProject mocks getting deployment with project data
func (m *MockDeploymentRepository) GetByIDWithProject(ctx context.Context, id, userID string) (*deployment.Deployment, error) {
	// TODO: Implement - call GetByIDWithProjectFunc if set
	// TODO: Return default behavior or error
	return nil, nil
}

// Create mocks creating a deployment
func (m *MockDeploymentRepository) Create(ctx context.Context, dep *deployment.Deployment) (*deployment.Deployment, error) {
	// TODO: Implement - call CreateFunc if set
	// TODO: Return default behavior or error
	return nil, nil
}

// DeleteByProjectID mocks deleting deployments by project ID
func (m *MockDeploymentRepository) DeleteByProjectID(ctx context.Context, projectID string) error {
	// TODO: Implement - call DeleteByProjectIDFunc if set
	// TODO: Return default behavior or error
	return nil
}

// MockECSService is a mock implementation of ECS service
type MockECSService struct {
	// TODO: Add fields for controlling mock behavior
	// TODO: Add fields for capturing method calls

	// Example fields:
	// RunTaskFunc func(ctx context.Context, envVars []ecs.EnvVar) (string, error)
	// RunTaskCalls [][]ecs.EnvVar  // Track all calls
}

// RunTask mocks running an ECS task
func (m *MockECSService) RunTask(ctx context.Context, envVars []ecs.EnvVar) (string, error) {
	// TODO: Implement - call RunTaskFunc if set
	// TODO: Capture call arguments
	// TODO: Return task ARN and error
	return "", nil
}

// MockLogsService is a mock implementation of logs service
type MockLogsService struct {
	// TODO: Add fields for controlling mock behavior
	// TODO: Add fields for capturing method calls

	// Example fields:
	// GetDeploymentLogsFunc func(ctx context.Context, deploymentID string) ([]logs.LogEvent, error)
}

// GetDeploymentLogs mocks getting deployment logs
func (m *MockLogsService) GetDeploymentLogs(ctx context.Context, deploymentID string) ([]logs.LogEvent, error) {
	// TODO: Implement - call GetDeploymentLogsFunc if set
	// TODO: Return default behavior or error
	return nil, nil
}

// MockJWKSCache is a mock implementation of JWKS cache
type MockJWKSCache struct {
	// TODO: Add fields for controlling mock behavior

	// Example fields:
	// GetKeyFunc func(keyID string) (interface{}, error)
}

// GetKey mocks getting a public key from JWKS
func (m *MockJWKSCache) GetKey(keyID string) (interface{}, error) {
	// TODO: Implement - call GetKeyFunc if set
	// TODO: Return test key or error
	return nil, nil
}

// Helper functions for creating mock data

// CreateMockProject creates a mock project for testing
func CreateMockProject(id, name, gitURL, subdomain, userID string) *project.Project {
	// TODO: Implement - create and return project with given fields
	// TODO: Set timestamps to current time
	return nil
}

// CreateMockDeployment creates a mock deployment for testing
func CreateMockDeployment(id, projectID, status string) *deployment.Deployment {
	// TODO: Implement - create and return deployment with given fields
	// TODO: Set timestamps to current time
	return nil
}

// CreateMockLogEvent creates a mock log event for testing
func CreateMockLogEvent(deploymentID, logMessage string) logs.LogEvent {
	// TODO: Implement - create and return log event
	// TODO: Set timestamp to current time
	return logs.LogEvent{}
}

// CreateMockProjects creates multiple mock projects
func CreateMockProjects(count int, userID string) []*project.Project {
	// TODO: Implement - create slice of projects
	// TODO: Use different IDs and names for each
	// TODO: Return slice
	return nil
}

// CreateMockDeployments creates multiple mock deployments
func CreateMockDeployments(count int, projectID string) []*deployment.Deployment {
	// TODO: Implement - create slice of deployments
	// TODO: Use different IDs and statuses for each
	// TODO: Return slice
	return nil
}
