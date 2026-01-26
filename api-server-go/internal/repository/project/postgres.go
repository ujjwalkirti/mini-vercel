package repository

import (
	"context"
	"database/sql"
	"encoding/json"

	deploymentdomain "github.com/ujjwalkirti/mini-vercel-api-server/internal/domain/deployment"
	domain "github.com/ujjwalkirti/mini-vercel-api-server/internal/domain/project"
	"github.com/ujjwalkirti/mini-vercel-api-server/internal/utils"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, p *domain.Project) error {
	// Generate UUID v4 if not provided
	if p.ID == "" {
		p.ID = utils.GenerateUUID()
	}

	query := `
		INSERT INTO projects (id, name, git_url, subdomain, custom_domain, user_id)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		p.ID,
		p.Name,
		p.GitURL,
		p.SubDomain,
		p.CustomDomain,
		p.UserID,
	)

	return err
}

func (r *Repository) ListByUser(ctx context.Context, userID string) ([]domain.Project, error) {
	query := `
		WITH latest_deployments AS (
			SELECT DISTINCT ON (project_id)
				id,
				project_id,
				status,
				created_at,
				updated_at
			FROM deployments
			ORDER BY project_id, created_at DESC
		)
		SELECT
			p.id,
			p.name,
			p.git_url,
			p.subdomain,
			p.custom_domain,
			p.user_id,
			p.created_at,
			p.updated_at,
			COALESCE(
				json_agg(
					json_build_object(
						'id', d.id,
						'projectId', d.project_id,
						'status', d.status,
						'createdAt', d.created_at,
						'updatedAt', d.updated_at
					)
				) FILTER (WHERE d.id IS NOT NULL),
				'[]'::json
			) AS deployments
		FROM projects p
		LEFT JOIN latest_deployments d ON p.id = d.project_id
		WHERE p.user_id = $1
		GROUP BY p.id, p.name, p.git_url, p.subdomain, p.custom_domain, p.user_id, p.created_at, p.updated_at
		ORDER BY p.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	projects := make([]domain.Project, 0)
	for rows.Next() {
		var p domain.Project
		var deploymentsJSON []byte

		if err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.GitURL,
			&p.SubDomain,
			&p.CustomDomain,
			&p.UserID,
			&p.CreatedAt,
			&p.UpdatedAt,
			&deploymentsJSON,
		); err != nil {
			return nil, err
		}

		// Always initialize as empty array
		p.Deployments = make([]deploymentdomain.Deployment, 0)

		// Only unmarshal if we have actual deployment data
		if len(deploymentsJSON) > 0 && string(deploymentsJSON) != "[]" && string(deploymentsJSON) != "null" {
			if err := json.Unmarshal(deploymentsJSON, &p.Deployments); err != nil {
				// If unmarshal fails, leave deployments as empty array
				p.Deployments = make([]deploymentdomain.Deployment, 0)
			}
		}

		projects = append(projects, p)
	}

	return projects, nil
}

func (r *Repository) GetByIDAndUserID(ctx context.Context, projectID, userID string) (domain.Project, error) {
	query := `
		SELECT
			p.id,
			p.name,
			p.git_url,
			p.subdomain,
			p.custom_domain,
			p.user_id,
			p.created_at,
			p.updated_at,
			COALESCE(
				json_agg(
					json_build_object(
						'id', d.id,
						'projectId', d.project_id,
						'status', d.status,
						'createdAt', d.created_at,
						'updatedAt', d.updated_at
					)
					ORDER BY d.created_at DESC
				) FILTER (WHERE d.id IS NOT NULL),
				'[]'::json
			) AS deployments
		FROM projects p
		LEFT JOIN deployments d ON p.id = d.project_id
		WHERE p.id = $1 AND p.user_id = $2
		GROUP BY p.id, p.name, p.git_url, p.subdomain, p.custom_domain, p.user_id, p.created_at, p.updated_at
	`

	var p domain.Project
	var deploymentsJSON []byte

	err := r.db.QueryRowContext(ctx, query, projectID, userID).Scan(
		&p.ID,
		&p.Name,
		&p.GitURL,
		&p.SubDomain,
		&p.CustomDomain,
		&p.UserID,
		&p.CreatedAt,
		&p.UpdatedAt,
		&deploymentsJSON,
	)

	if err != nil {
		return p, err
	}

	// Always initialize as empty array
	p.Deployments = make([]deploymentdomain.Deployment, 0)

	// Only unmarshal if we have actual deployment data
	if len(deploymentsJSON) > 0 && string(deploymentsJSON) != "[]" && string(deploymentsJSON) != "null" {
		if err := json.Unmarshal(deploymentsJSON, &p.Deployments); err != nil {
			// If unmarshal fails, leave deployments as empty array
			p.Deployments = make([]deploymentdomain.Deployment, 0)
		}
	}

	return p, nil
}

func (r *Repository) Delete(ctx context.Context, projectID string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM projects WHERE id = $1`, projectID)
	return err
}
