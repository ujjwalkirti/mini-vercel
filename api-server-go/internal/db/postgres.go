package db

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/ujjwalkirti/mini-vercel-api-server/internal/config"

	_ "github.com/lib/pq"
)

func Connect() (*sql.DB, error) {
	logger := config.GetLogger()

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	logger.Info("Connected to database")
	return db, nil
}
