package deployment

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/ujjwalkirti/mini-vercel-api-server/internal/domain/deployment"
	"github.com/ujjwalkirti/mini-vercel-api-server/internal/middleware"
	repository "github.com/ujjwalkirti/mini-vercel-api-server/internal/repository/deployment"
	projectRepo "github.com/ujjwalkirti/mini-vercel-api-server/internal/repository/project"
	"github.com/ujjwalkirti/mini-vercel-api-server/internal/service/ecs"
	"github.com/ujjwalkirti/mini-vercel-api-server/internal/service/logs"
	"github.com/ujjwalkirti/mini-vercel-api-server/internal/utils"
)

type Handler struct {
	repo        *repository.Repository
	projectRepo *projectRepo.Repository
	ecsService  *ecs.Service
	logsService *logs.Service
}

func NewHandler(repo *repository.Repository, projectRepo *projectRepo.Repository, ecsService *ecs.Service, logsService *logs.Service) *Handler {
	return &Handler{
		repo:        repo,
		projectRepo: projectRepo,
		ecsService:  ecsService,
		logsService: logsService,
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
	if !utils.IsValidUUID(projectID) {
		utils.BadRequest(w, "Invalid project ID")
		return
	}

	// TODO: Fetch all deployments for the project
	deployments, err := h.repo.GetByProjectID(r.Context(), projectID, user.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.NotFound(w, "Failed to fetch deployments "+err.Error())
			return
		}

		utils.InternalServerError(w, "Failed to fetch deployments")
		return
	}

	// Return deployments array (matching Express API response)
	utils.Success(w, deployments, "Deployments fetched successfully.")
}

// GetDeployment handles GET /deployments/:id
// Returns a specific deployment
// Verifies user owns the parent project
func (h *Handler) GetDeployment(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		utils.Unauthorized(w, "Unauthorized")
		return
	}

	id := chi.URLParam(r, "id")

	// TODO: Validate UUID format for id
	if !utils.IsValidUUID(id) {
		utils.BadRequest(w, "Invalid deployment ID")
		return
	}

	// TODO: Fetch deployment with project data
	deployment, err := h.repo.GetByIDWithProject(r.Context(), id, user.ID)
	if err != nil {
		utils.InternalServerError(w, "Failed to fetch deployment")
		return
	}

	// Return deployment object (matching Express API response)
	utils.Success(w, deployment)
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
	if req.ProjectID == "" {
		utils.BadRequest(w, "Invalid request body. Project ID is required")
		return
	}

	// TODO: Verify project exists and user owns it
	project, err := h.projectRepo.GetByIDAndUserID(r.Context(), req.ProjectID, user.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.NotFound(w, "Project not found")
			return
		}
		utils.InternalServerError(w, "Failed to fetch project")
		return
	}

	// TODO: Create deployment with status "QUEUED"
	deployment, err := h.repo.Create(r.Context(), &deployment.Deployment{ProjectID: req.ProjectID, Status: "QUEUED"})

	// TODO: Prepare environment variables for ECS task
	envVars := []ecs.EnvVar{
		{Name: "PROJECT_ID", Value: req.ProjectID},
		{Name: "GIT_REPOSITORY_URL", Value: project.GitURL},
		{Name: "DEPLOYMENT_ID", Value: deployment.ID},
		{Name: "KAFKA_BROKERS", Value: os.Getenv("KAFKA_BROKERS")},
		{Name: "KAFKA_CLIENT_ID", Value: os.Getenv("KAFKA_CLIENT_ID")},
		{Name: "KAFKA_USERNAME", Value: os.Getenv("KAFKA_USERNAME")},
		{Name: "KAFKA_PASSWORD", Value: os.Getenv("KAFKA_PASSWORD")},
		{Name: "R2_ACCOUNT_ID", Value: os.Getenv("R2_ACCOUNT_ID")},
		{Name: "R2_ACCESS_KEY_ID", Value: os.Getenv("R2_ACCESS_KEY_ID")},
		{Name: "R2_SECRET_ACCESS_KEY", Value: os.Getenv("R2_SECRET_ACCESS_KEY")},
		{Name: "R2_BUCKET_NAME", Value: os.Getenv("R2_BUCKET_NAME")},
	}

	// TODO: Trigger AWS ECS task to build and deploy
	_, err = h.ecsService.RunTask(r.Context(), envVars)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {

		}
		utils.InternalServerError(w, "Failed to trigger ECS task")
		return
	}

	// Return deployment response (matching Express API response format)
	// Express returns: { deploymentId, status: "Queued", url: "{subdomain}.localhost:8001" }
	utils.Success(w, map[string]interface{}{
		"deploymentId": deployment.ID, // Will be: deployment.ID
		"status":       "Queued",
		"url":          project.SubDomain + ".localhost:8001",
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
	if !utils.IsValidUUID(id) {
		utils.BadRequest(w, "Invalid deployment ID")
		return
	}

	// TODO: Fetch deployment with project data
	deployment, err := h.repo.GetByIDWithProject(r.Context(), id, user.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.NotFound(w, "Deployment not found")
			return
		}
		utils.InternalServerError(w, "Failed to fetch deployment")
		return
	}

	// Query logs service for deployment logs
	var logEvents []logs.LogEvent
	if h.logsService != nil {
		logEvents, err = h.logsService.GetDeploymentLogs(r.Context(), id)
		if err != nil {
			utils.InternalServerError(w, "Failed to fetch logs")
			return
		}
	} else {
		// If logs service is not configured, return empty logs
		logEvents = []logs.LogEvent{}
	}

	// Return logs response
	utils.Success(w, map[string]interface{}{
		"deployment": map[string]interface{}{
			"id":     deployment.ID,
			"status": deployment.Status,
		},
		"logs": logEvents,
	}, "Logs retrieved successfully")
}
