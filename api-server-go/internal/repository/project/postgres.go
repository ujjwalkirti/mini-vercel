package repository

import (
	"context"
	"database/sql"

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
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, name, git_url, subdomain, custom_domain, user_id, created_at, updated_at
		 FROM projects WHERE user_id = $1`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Initialize as empty slice to ensure JSON serializes to [] instead of null
	projects := make([]domain.Project, 0)
	for rows.Next() {
		var p domain.Project
		if err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.GitURL,
			&p.SubDomain,
			&p.CustomDomain,
			&p.UserID,
			&p.CreatedAt,
			&p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}

	return projects, nil
}

func (r *Repository) GetByIDAndUserID(ctx context.Context, projectID, userID string) (domain.Project, error) {
	var p domain.Project
	err := r.db.QueryRowContext(ctx, `SELECT id, name, git_url, subdomain, custom_domain, user_id, created_at, updated_at
		FROM projects WHERE id = $1 AND user_id = $2`, projectID, userID).Scan(
		&p.ID,
		&p.Name,
		&p.GitURL,
		&p.SubDomain,
		&p.CustomDomain,
		&p.UserID,
		&p.CreatedAt,
		&p.UpdatedAt,
	)

	return p, err
}

func (r *Repository) Delete(ctx context.Context, projectID string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM projects WHERE id = $1`, projectID)
	return err
}
