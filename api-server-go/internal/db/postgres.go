package db

import (
	"database/sql"
	"fmt"

	"github.com/ujjwalkirti/mini-vercel-api-server/internal/config"

	_ "github.com/lib/pq"
)

func Connect() (*sql.DB, error) {
	logger := config.GetLogger()

	databaseConfig := config.GetDBConfig()
	if databaseConfig.DATABASE_URL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	db, err := sql.Open(databaseConfig.DATABASE_TYPE, databaseConfig.DATABASE_URL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Configure connection pool
	db.SetMaxOpenConns(databaseConfig.DATABASE_MAX_CONNECTIONS)
	db.SetMaxIdleConns(databaseConfig.DATABASE_IDLE_CONNECTIONS)
	db.SetConnMaxLifetime(databaseConfig.DATABASE_CONN_MAX_LIFETIME)
	db.SetConnMaxIdleTime(databaseConfig.DATABASE_CONN_MAX_IDLE_TIME)

	logger.Info("Connected to database")
	return db, nil
}
