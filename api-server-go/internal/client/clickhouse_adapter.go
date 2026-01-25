package client

import (
	"context"
	"database/sql"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/ujjwalkirti/mini-vercel-api-server/internal/repository/logs"
)

// ClickHouseAdapter adapts ClickHouse driver.Conn to implement common database interfaces
// This allows the ClickHouse client to be used with generic repository interfaces
type ClickHouseAdapter struct {
	conn driver.Conn
}

// NewClickHouseAdapter creates a new adapter wrapping a ClickHouse connection
func NewClickHouseAdapter(conn driver.Conn) *ClickHouseAdapter {
	return &ClickHouseAdapter{
		conn: conn,
	}
}

// QueryContext executes a query that returns rows
func (a *ClickHouseAdapter) QueryContext(ctx context.Context, query string, args ...any) (logs.Rows, error) {
	rows, err := a.conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &clickhouseRows{rows: rows}, nil
}

// QueryRowContext executes a query that is expected to return at most one row
func (a *ClickHouseAdapter) QueryRowContext(ctx context.Context, query string, args ...any) logs.Row {
	return &clickhouseRow{
		conn:  a.conn,
		ctx:   ctx,
		query: query,
		args:  args,
	}
}

// ExecContext executes a query without returning any rows
func (a *ClickHouseAdapter) ExecContext(ctx context.Context, query string, args ...any) (logs.Result, error) {
	err := a.conn.Exec(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &clickhouseResult{}, nil
}

// clickhouseRows wraps driver.Rows
type clickhouseRows struct {
	rows driver.Rows
}

func (r *clickhouseRows) Close() error {
	return r.rows.Close()
}

func (r *clickhouseRows) Next() bool {
	return r.rows.Next()
}

func (r *clickhouseRows) Scan(dest ...any) error {
	return r.rows.Scan(dest...)
}

func (r *clickhouseRows) Err() error {
	return r.rows.Err()
}

// clickhouseRow implements Row interface for single row queries
type clickhouseRow struct {
	conn  driver.Conn
	ctx   context.Context
	query string
	args  []any
}

func (r *clickhouseRow) Scan(dest ...any) error {
	row := r.conn.QueryRow(r.ctx, r.query, r.args...)
	return row.Scan(dest...)
}

// clickhouseResult implements Result interface
type clickhouseResult struct{}

func (r *clickhouseResult) LastInsertId() (int64, error) {
	// ClickHouse doesn't support LastInsertId
	return 0, sql.ErrNoRows
}

func (r *clickhouseResult) RowsAffected() (int64, error) {
	// ClickHouse doesn't provide affected rows count in the same way
	// This could be enhanced if needed
	return 0, nil
}
