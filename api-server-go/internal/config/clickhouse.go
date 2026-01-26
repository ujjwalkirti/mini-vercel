package config

import (
	"os"
)

// ClickHouseConfig holds ClickHouse connection configuration
type ClickHouseConfig struct {
	Host     string
	Port     string
	Database string
	Username string
	Password string
}

// GetClickHouseConfig returns ClickHouse configuration from environment variables
// Using a function for lazy evaluation ensures env vars are read after .env file is loaded
func GetClickHouseConfig() ClickHouseConfig {
	return ClickHouseConfig{
		Host:     os.Getenv("CLICKHOUSE_HOST"),
		Port:     os.Getenv("CLICKHOUSE_PORT"),
		Database: os.Getenv("CLICKHOUSE_DATABASE"),
		Username: os.Getenv("CLICKHOUSE_USERNAME"),
		Password: os.Getenv("CLICKHOUSE_PASSWORD"),
	}
}
