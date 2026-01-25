package project

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sio/coolname"
	"github.com/ujjwalkirti/mini-vercel-api-server/internal/domain/project"
	"github.com/ujjwalkirti/mini-vercel-api-server/internal/middleware"
	deploymentRepository "github.com/ujjwalkirti/mini-vercel-api-server/internal/repository/deployment"
	projectRepository "github.com/ujjwalkirti/mini-vercel-api-server/internal/repository/project"
	"github.com/ujjwalkirti/mini-vercel-api-server/internal/utils"
)

type Handler struct {
	repo *projectRepository.Repository

	deploymentRepo *deploymentRepository.Repository
}

func NewHandler(repo *projectRepository.Repository, deploymentRepo *deploymentRepository.Repository) *Handler {
	return &Handler{
		repo:           repo,
		deploymentRepo: deploymentRepo,
	}
}

// GetProjects handles GET /projects
// Returns all projects for the authenticated user
// Response includes the most recent deployment for each project
func (h *Handler) GetProjects(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		utils.Unauthorized(w, "Unauthorized")
		return
	}

	// TODO: Fetch all projects for user.ID from database
	// TODO: Include the most recent deployment for each project
	projects, err := h.repo.ListByUser(r.Context(), user.ID)

	if err != nil {
		utils.InternalServerError(w, "Failed to fetch projects")
		return
	}

	// Return projects array directly (matching Express API response)
	utils.Success(w, projects)
}

// GetProject handles GET /projects/:id
// Returns a specific project with all its deployments
// Verifies user owns the project
func (h *Handler) GetProject(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		utils.Unauthorized(w, "Unauthorized")
		return
	}

	id := chi.URLParam(r, "id")

	// TODO: Validate UUID format for id
	if !utils.IsValidUUID(id) {
		utils.BadRequest(w, "Invalid project ID")
		return
	}

	// TODO: Fetch project by ID and verify ownership
	project, err := h.repo.GetByIDAndUserID(r.Context(), id, user.ID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.NotFound(w, "Project not found")
			return
		}
		utils.InternalServerError(w, "Failed to fetch project")
		return
	}

	// Return project object (matching Express API response)
	utils.Success(w, project)
}

func generateSubdomain() string {
	slug, _ := coolname.SlugN(3)
	return slug
}

// CreateProject handles POST /projects
// Creates a new project for the authenticated user
// Request body: { "name": string, "github_url": string }
// Generates a random subdomain for the project
func (h *Handler) CreateProject(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		utils.Unauthorized(w, "Unauthorized")
		return
	}

	type CreateProjectRequest struct {
		Name      string `json:"name"`
		GithubURL string `json:"github_url"`
	}

	var req CreateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid request body")
		return
	}

	// TODO: Validate request (name and github_url are required)
	if req.Name == "" || req.GithubURL == "" {
		utils.BadRequest(w, "Invalid request body. Name and github_url are required")
		return
	}

	// TODO: Generate random subdomain (3-word slug)
	subdomain := generateSubdomain()

	// Create project in database
	newProject := &project.Project{
		Name:      req.Name,
		GitURL:    req.GithubURL,
		SubDomain: subdomain,
		UserID:    user.ID,
	}

	err := h.repo.Create(r.Context(), newProject)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.BadRequest(w, "Error creating porject: "+err.Error())
			return
		}
		utils.InternalServerError(w, "Failed to create project")
		return
	}

	// Return the created project object (matching Express API response)
	utils.Created(w, newProject, "Project created successfully")
}

// UpdateProject handles PUT /projects/:id
// Updates an existing project
// Verifies user owns the project
func (h *Handler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		utils.Unauthorized(w, "Unauthorized")
		return
	}

	id := chi.URLParam(r, "id")

	// TODO: Decode request body
	// TODO: Validate ownership
	// TODO: Update project in database

	// Placeholder response
	utils.Success(w, map[string]interface{}{
		"id":     id,
		"userId": user.ID,
	}, "Project updated successfully")
}

// DeleteProject handles DELETE /projects/:id
// Deletes a project and all its deployments
// Verifies user owns the project
func (h *Handler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		utils.Unauthorized(w, "Unauthorized")
		return
	}

	id := chi.URLParam(r, "id")

	// TODO: Validate UUID format for id
	if !utils.IsValidUUID(id) {
		utils.BadRequest(w, "Invalid project ID")
		return
	}

	// TODO: Verify project exists and user owns it
	project, err := h.repo.GetByIDAndUserID(r.Context(), id, user.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.NotFound(w, "Project not found")
			return
		}
		utils.InternalServerError(w, "Failed to fetch project")
		return
	}

	// TODO: Delete all deployments for this project
	err = h.deploymentRepo.DeleteByProjectID(r.Context(), id)

	// TODO: Delete the project
	err = h.repo.Delete(r.Context(), project.ID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.NotFound(w, "Error deleting project: "+err.Error())
			return
		}

		utils.InternalServerError(w, "Failed to delete project")
		return
	}

	// Placeholder response
	utils.Success(w, nil, "Project deleted successfully")
}
