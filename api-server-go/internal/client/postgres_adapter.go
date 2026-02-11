package client

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/ujjwalkirti/mini-vercel-api-server/internal/repository/logs"
)

type PostgresAdapter struct {
	db *sql.DB
}

func NewPostgresAdapter(db *sql.DB) *PostgresAdapter {
	return &PostgresAdapter{db: db}
}

func convertPlaceholders(query string) string {
	placeholder := 1
	result := strings.Builder{}
	for i := 0; i < len(query); i++ {
		if query[i] == '?' {
			result.WriteString(fmt.Sprintf("$%d", placeholder))
			placeholder++
		} else {
			result.WriteByte(query[i])
		}
	}
	return result.String()
}

func (a *PostgresAdapter) QueryContext(ctx context.Context, query string, args ...any) (logs.Rows, error) {
	query = convertPlaceholders(query)
	rows, err := a.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &postgresRows{rows: rows}, nil
}

func (a *PostgresAdapter) QueryRowContext(ctx context.Context, query string, args ...any) logs.Row {
	query = convertPlaceholders(query)
	return &postgresRow{row: a.db.QueryRowContext(ctx, query, args...)}
}

func (a *PostgresAdapter) ExecContext(ctx context.Context, query string, args ...any) (logs.Result, error) {
	query = convertPlaceholders(query)
	result, err := a.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &postgresResult{result: result}, nil
}

// postgresRows wraps sql.Rows to implement logs.Rows interface
type postgresRows struct {
	rows *sql.Rows
}

func (r *postgresRows) Close() error {
	return r.rows.Close()
}

func (r *postgresRows) Next() bool {
	return r.rows.Next()
}

func (r *postgresRows) Scan(dest ...any) error {
	return r.rows.Scan(dest...)
}

func (r *postgresRows) Err() error {
	return r.rows.Err()
}

// postgresRow wraps sql.Row to implement logs.Row interface
type postgresRow struct {
	row *sql.Row
}

func (r *postgresRow) Scan(dest ...any) error {
	return r.row.Scan(dest...)
}

// postgresResult wraps sql.Result to implement logs.Result interface
type postgresResult struct {
	result sql.Result
}

func (r *postgresResult) LastInsertId() (int64, error) {
	return r.result.LastInsertId()
}

func (r *postgresResult) RowsAffected() (int64, error) {
	return r.result.RowsAffected()
}
