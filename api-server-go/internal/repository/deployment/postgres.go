package deployment

import (
	"context"
	"database/sql"

	domain "github.com/ujjwalkirti/mini-vercel-api-server/internal/domain/deployment"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, d *domain.Deployment) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO deployments (id, project_id, status)
		 VALUES ($1, $2, $3)`,
		d.ID,
		d.ProjectID,
		d.Status,
	)
	return err
}
