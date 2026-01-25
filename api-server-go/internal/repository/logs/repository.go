package logs

import (
	"context"
)

// Repository defines the interface for log storage operations
// This interface is technology-agnostic and can be implemented by any storage backend
// (ClickHouse, PostgreSQL WAL, ElasticSearch, etc.)
type Repository interface {
	QueryContext(ctx context.Context, query string, args ...any) (Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) Row
	ExecContext(ctx context.Context, query string, args ...any) (Result, error)
}

// Rows interface for iterating over query results
type Rows interface {
	Close() error
	Next() bool
	Scan(dest ...any) error
	Err() error
}

// Row interface for single row queries
type Row interface {
	Scan(dest ...any) error
}

// Result interface for execution results
type Result interface {
	LastInsertId() (int64, error)
	RowsAffected() (int64, error)
}
