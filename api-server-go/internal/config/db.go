package config

import (
	"os"
	"time"
)

type DBConfig struct {
	DATABASE_URL                string
	DATABASE_TYPE               string
	DATABASE_MAX_CONNECTIONS    int
	DATABASE_IDLE_CONNECTIONS   int
	DATABASE_CONN_MAX_LIFETIME  time.Duration
	DATABASE_CONN_MAX_IDLE_TIME time.Duration
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func GetDBConfig() DBConfig {
	return DBConfig{
		DATABASE_URL:                getEnvOrDefault("DATABASE_URL", ""),
		DATABASE_TYPE:               getEnvOrDefault("DATABASE_TYPE", "postgres"),
		DATABASE_MAX_CONNECTIONS:    getEnvAsInt("DATABASE_MAX_CONNECTIONS", 10),
		DATABASE_IDLE_CONNECTIONS:   getEnvAsInt("DATABASE_IDLE_CONNECTIONS", 10),
		DATABASE_CONN_MAX_LIFETIME:  getEnvAsDuration("DATABASE_CONN_MAX_LIFETIME", 5*time.Minute),
		DATABASE_CONN_MAX_IDLE_TIME: getEnvAsDuration("DATABASE_CONN_MAX_IDLE_TIME", 5*time.Minute),
	}
}
