package project

import (
	"context"
	"database/sql"

	domain "github.com/ujjwalkirti/mini-vercel-api-server/internal/domain/project"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, p *domain.Project) error {
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

	var projects []domain.Project
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
