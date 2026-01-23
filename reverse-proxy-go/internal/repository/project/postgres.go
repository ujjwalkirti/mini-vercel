package project

import (
	"context"
	"database/sql"
	"fmt"

	"reverse-proxy/internal/domain/deployment"
	"reverse-proxy/internal/domain/project"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// FindBySubdomain finds a project by subdomain with its latest READY deployment
func (r *Repository) FindBySubdomain(ctx context.Context, subdomain string) (*project.Project, error) {
	query := `
		SELECT
			p.id, p.name, p.git_url, p.subdomain, p.custom_domain,
			p.user_id, p.created_at, p.updated_at,
			d.id, d.project_id, d.status, d.created_at, d.updated_at
		FROM projects p
		LEFT JOIN deployments d ON p.id = d.project_id
		WHERE p.subdomain = $1
		AND d.status = $2
		ORDER BY d.created_at DESC
		LIMIT 1
	`

	var proj project.Project
	var deploy deployment.Deployment
	var customDomain sql.NullString

	err := r.db.QueryRowContext(ctx, query, subdomain, deployment.StatusReady).Scan(
		&proj.ID,
		&proj.Name,
		&proj.GitURL,
		&proj.Subdomain,
		&customDomain,
		&proj.UserID,
		&proj.CreatedAt,
		&proj.UpdatedAt,
		&deploy.ID,
		&deploy.ProjectID,
		&deploy.Status,
		&deploy.CreatedAt,
		&deploy.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no project found with subdomain: %s", subdomain)
		}
		return nil, fmt.Errorf("failed to query project: %w", err)
	}

	if customDomain.Valid {
		proj.CustomDomain = &customDomain.String
	}

	proj.Deployments = []deployment.Deployment{deploy}

	return &proj, nil
}

// FindByCustomDomain finds a project by custom domain with its latest READY deployment
func (r *Repository) FindByCustomDomain(ctx context.Context, customDomain string) (*project.Project, error) {
	query := `
		SELECT
			p.id, p.name, p.git_url, p.subdomain, p.custom_domain,
			p.user_id, p.created_at, p.updated_at,
			d.id, d.project_id, d.status, d.created_at, d.updated_at
		FROM projects p
		LEFT JOIN deployments d ON p.id = d.project_id
		WHERE p.custom_domain = $1
		AND d.status = $2
		ORDER BY d.created_at DESC
		LIMIT 1
	`

	var proj project.Project
	var deploy deployment.Deployment
	var customDomainVal sql.NullString

	err := r.db.QueryRowContext(ctx, query, customDomain, deployment.StatusReady).Scan(
		&proj.ID,
		&proj.Name,
		&proj.GitURL,
		&proj.Subdomain,
		&customDomainVal,
		&proj.UserID,
		&proj.CreatedAt,
		&proj.UpdatedAt,
		&deploy.ID,
		&deploy.ProjectID,
		&deploy.Status,
		&deploy.CreatedAt,
		&deploy.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no project found with custom domain: %s", customDomain)
		}
		return nil, fmt.Errorf("failed to query project: %w", err)
	}

	if customDomainVal.Valid {
		proj.CustomDomain = &customDomainVal.String
	}

	proj.Deployments = []deployment.Deployment{deploy}

	return &proj, nil
}
