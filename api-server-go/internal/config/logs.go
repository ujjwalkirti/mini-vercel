package config

import (
	"os"
	"strings"
)

// LogsConfig holds logs service configuration
type LogsConfig struct {
	// RepositoryType specifies which repository to use for logs storage
	// Valid values: "clickhouse", "postgres"
	// If empty, auto-detection will be used based on ClickHouse configuration
	RepositoryType string
}

// GetLogsConfig returns logs configuration from environment variables
func GetLogsConfig() LogsConfig {
	repoType := strings.ToLower(strings.TrimSpace(os.Getenv("LOGS_REPOSITORY_TYPE")))

	// Normalize the repository type
	switch repoType {
	case "clickhouse":
		repoType = "clickhouse"
	case "postgres", "postgresql":
		repoType = "postgres"
	default:
		// Empty string triggers auto-detection
		repoType = ""
	}

	return LogsConfig{
		RepositoryType: repoType,
	}
}
