package client

import (
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
func NewLogRepository(repoType LogRepositoryType) (logs.Repository, error) {
	switch repoType {
	case ClickHouseRepository:
		client, err := NewClickHouseClient()
		if err != nil {
			return nil, fmt.Errorf("failed to create clickhouse client: %w", err)
		}
		return NewClickHouseAdapter(client.GetConn()), nil

	case PostgresRepository:
		// Future implementation for PostgreSQL WAL or other backends
		return nil, fmt.Errorf("postgres repository not yet implemented")

	default:
		return nil, fmt.Errorf("unsupported repository type: %s", repoType)
	}
}
