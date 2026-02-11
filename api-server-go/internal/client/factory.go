package client

import (
	"database/sql"
	"fmt"

	"github.com/ujjwalkirti/mini-vercel-api-server/internal/repository/logs"
)

// LogRepositoryType represents the type of log repository to create
type LogRepositoryType string

const (
	// ClickHouseRepository uses ClickHouse for log storage
	ClickHouseRepository LogRepositoryType = "clickhouse"
	// PostgresRepository uses PostgreSQL for log storage (future implementation)
	PostgresRepository LogRepositoryType = "postgres"
)

// NewLogRepository creates a new log repository based on the specified type
// This factory pattern allows easy switching between different storage backends
func NewLogRepository(repoType LogRepositoryType, db *sql.DB) (logs.Repository, error) {
	switch repoType {
	case ClickHouseRepository:
		client, err := NewClickHouseClient()
		if err != nil {
			return nil, err
		}
		return NewClickHouseAdapter(client.GetConn()), nil

	case PostgresRepository:
		if db == nil {
			return nil, fmt.Errorf("database connection is required")
		}
		return NewPostgresAdapter(db), nil

	default:
		return nil, fmt.Errorf("unsupported repository type: %s", repoType)
	}
}
