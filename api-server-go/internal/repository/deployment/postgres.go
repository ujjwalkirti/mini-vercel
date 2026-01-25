package repository

import (
	"context"
	"database/sql"
	"encoding/json"

	domain "github.com/ujjwalkirti/mini-vercel-api-server/internal/domain/deployment"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, d *domain.Deployment) (domain.Deployment, error) {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO deployments (id, project_id, status)
		 VALUES ($1, $2, $3)`,
		d.ID,
		d.ProjectID,
		d.Status,
	)
	if err != nil {
		return domain.Deployment{}, err
	}
	return *d, nil
}

func (r *Repository) GetByProjectID(ctx context.Context, projectID string, userID string) ([]domain.Deployment, error) {
	var jsonData []byte
	err := r.db.QueryRowContext(ctx, `
		SELECT COALESCE(json_agg(
			json_build_object(
				'id', d.id,
				'project_id', d.project_id,
				'status', d.status
			)
		), '[]'::json)
		FROM deployments d
		INNER JOIN projects p ON d.project_id = p.id
		WHERE d.project_id = $1 AND p.user_id = $2
	`, projectID, userID).Scan(&jsonData)

	if err != nil {
		return nil, err
	}

	var deployments []domain.Deployment
	if err := json.Unmarshal(jsonData, &deployments); err != nil {
		return nil, err
	}

	return deployments, nil
}

func (r *Repository) GetByIDWithProject(ctx context.Context, id string, userID string) (domain.Deployment, error) {
	var d domain.Deployment
	err := r.db.QueryRowContext(ctx, `
		SELECT d.id, d.project_id, d.status, d.created_at, d.updated_at
		FROM deployments d
		INNER JOIN projects p ON d.project_id = p.id
		WHERE d.id = $1 AND p.user_id = $2
	`, id, userID).Scan(
		&d.ID,
		&d.ProjectID,
		&d.Status,
		&d.CreatedAt,
		&d.UpdatedAt,
	)
	return d, err
}

func (r *Repository) DeleteByProjectID(ctx context.Context, projectID string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM deployments WHERE project_id = $1`, projectID)
	return err
}
