package deployment

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ujjwalkirti/mini-vercel-api-server/internal/middleware"
	repository "github.com/ujjwalkirti/mini-vercel-api-server/internal/repository/deployment"
	"github.com/ujjwalkirti/mini-vercel-api-server/internal/utils"
)

type Handler struct {
	repo *repository.Repository
}

func NewHandler(repo *repository.Repository) *Handler {
	return &Handler{
		repo: repo,
	}
}

// GetDeploymentsByProject handles GET /projects/:projectId/deployments
// Returns all deployments for a specific project
// Verifies user owns the project
func (h *Handler) GetDeploymentsByProject(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		utils.Unauthorized(w, "Unauthorized")
		return
	}

	projectID := chi.URLParam(r, "projectId")

	// TODO: Validate UUID format for projectID

	// TODO: Verify project exists and user owns it
	// project, err := h.projectRepo.GetByIDAndUserID(projectID, user.ID)
	// if err != nil {
	//     if errors.Is(err, sql.ErrNoRows) {
	//         utils.NotFound(w, "Project not found")
	//         return
	//     }
	//     utils.InternalServerError(w, "Failed to fetch project")
	//     return
	// }

	// TODO: Fetch all deployments for the project
	// deployments, err := h.repo.GetByProjectID(projectID)

	// Return deployments array (matching Express API response)
	utils.Success(w, []interface{}{})
}

// GetDeployment handles GET /deployments/:id
// Returns a specific deployment
// Verifies user owns the parent project
func (h *Handler) GetDeployment(w http.ResponseWriter, r *http.Request) {
	_, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		utils.Unauthorized(w, "Unauthorized")
		return
	}

	_ = chi.URLParam(r, "id")

	// TODO: Validate UUID format for id

	// TODO: Fetch deployment with project data
	// deployment, err := h.repo.GetByIDWithProject(id)
	// if err != nil {
	//     if errors.Is(err, sql.ErrNoRows) {
	//         utils.NotFound(w, "Deployment not found")
	//         return
	//     }
	//     utils.InternalServerError(w, "Failed to fetch deployment")
	//     return
	// }

	// TODO: Verify user owns the parent project
	// if deployment.Project.UserID != user.ID {
	//     utils.Forbidden(w, "Access denied")
	//     return
	// }

	// Return deployment object (matching Express API response)
	// When implemented, uncomment: utils.Success(w, deployment)
	utils.InternalServerError(w, "Not implemented")
}

// CreateDeployment handles POST /deploy
// Creates a new deployment and triggers the build process
// Request body: { "project_id": string }
// Queues the deployment to AWS ECS
func (h *Handler) CreateDeployment(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		utils.Unauthorized(w, "Unauthorized")
		return
	}

	type CreateDeploymentRequest struct {
		ProjectID string `json:"project_id"`
	}

	var req CreateDeploymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid request body")
		return
	}

	// TODO: Validate request (project_id is required)

	// TODO: Verify project exists and user owns it
	// project, err := h.projectRepo.GetByIDAndUserID(req.ProjectID, user.ID)
	// if err != nil {
	//     if errors.Is(err, sql.ErrNoRows) {
	//         utils.NotFound(w, "Project not found")
	//         return
	//     }
	//     utils.InternalServerError(w, "Failed to fetch project")
	//     return
	// }

	// TODO: Create deployment with status "QUEUED"
	// deployment, err := h.repo.Create(req.ProjectID, "QUEUED")

	// TODO: Prepare environment variables for ECS task
	// envVars := []EnvVar{
	//     {Name: "PROJECT_ID", Value: req.ProjectID},
	//     {Name: "GIT_REPOSITORY_URL", Value: project.GitURL},
	//     {Name: "DEPLOYMENT_ID", Value: deployment.ID},
	//     {Name: "KAFKA_BROKERS", Value: os.Getenv("KAFKA_BROKERS")},
	//     {Name: "KAFKA_CLIENT_ID", Value: os.Getenv("KAFKA_CLIENT_ID")},
	//     {Name: "KAFKA_USERNAME", Value: os.Getenv("KAFKA_USERNAME")},
	//     {Name: "KAFKA_PASSWORD", Value: os.Getenv("KAFKA_PASSWORD")},
	//     {Name: "R2_ACCOUNT_ID", Value: os.Getenv("R2_ACCOUNT_ID")},
	//     {Name: "R2_ACCESS_KEY_ID", Value: os.Getenv("R2_ACCESS_KEY_ID")},
	//     {Name: "R2_SECRET_ACCESS_KEY", Value: os.Getenv("R2_SECRET_ACCESS_KEY")},
	//     {Name: "R2_BUCKET_NAME", Value: os.Getenv("R2_BUCKET_NAME")},
	// }

	// TODO: Trigger AWS ECS task to build and deploy
	// err = h.awsECSService.RunTask(envVars)

	// Return deployment response (matching Express API response format)
	// Express returns: { deploymentId, status: "Queued", url: "{subdomain}.localhost:8001" }
	utils.Success(w, map[string]interface{}{
		"deploymentId": "placeholder-id",        // Will be: deployment.ID
		"status":       "Queued",
		"url":          "subdomain.localhost:8001", // Will be: project.SubDomain + ".localhost:8001"
	}, "Build queued successfully")
}

// GetDeploymentLogs handles GET /deployments/:id/logs
// Returns logs for a specific deployment from ClickHouse
// Verifies user owns the parent project
func (h *Handler) GetDeploymentLogs(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		utils.Unauthorized(w, "Unauthorized")
		return
	}

	id := chi.URLParam(r, "id")

	// TODO: Validate UUID format for id

	// TODO: Fetch deployment with project data
	// deployment, err := h.repo.GetByIDWithProject(id)
	// if err != nil {
	//     if errors.Is(err, sql.ErrNoRows) {
	//         utils.NotFound(w, "Deployment not found")
	//         return
	//     }
	//     utils.InternalServerError(w, "Failed to fetch deployment")
	//     return
	// }

	// TODO: Verify user owns the parent project
	// if deployment.Project.UserID != user.ID {
	//     utils.Forbidden(w, "Access denied")
	//     return
	// }

	// TODO: Query ClickHouse for deployment logs
	// query := "SELECT event_id, deployment_id, log, timestamp FROM log_events WHERE deployment_id = ? ORDER BY timestamp ASC"
	// logs, err := h.clickhouseClient.Query(query, id)

	// Placeholder response
	utils.Success(w, map[string]interface{}{
		"deployment": map[string]interface{}{
			"id": id,
		},
		"logs": []interface{}{},
	}, "Logs retrieved successfully")
}
