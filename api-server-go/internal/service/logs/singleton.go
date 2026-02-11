package logs

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/ujjwalkirti/mini-vercel-api-server/internal/client"
	"github.com/ujjwalkirti/mini-vercel-api-server/internal/config"
	"go.uber.org/zap"
)

var (
	instance *Service
	once     sync.Once
	initErr  error
)

// GetInstance returns the singleton instance of the logs service.
// It initializes the service on first call and reuses it thereafter.
// This ensures only one logs service is created across the application.
//
// The repoType parameter allows you to choose between ClickHouse and PostgreSQL:
//   - "clickhouse": Use ClickHouse for logs storage
//   - "postgres": Use PostgreSQL for logs storage
//
// If repoType is empty, it auto-detects based on configuration:
//   - Uses ClickHouse if configured, otherwise falls back to PostgreSQL
func GetInstance(db *sql.DB, repoType string) (*Service, error) {
	once.Do(func() {
		logger := config.GetLogger()
		var selectedType client.LogRepositoryType

		// Auto-detect if not specified
		if repoType == "" {
			clickhouseCfg := config.GetClickHouseConfig()
			if clickhouseCfg.Host != "" {
				selectedType = client.ClickHouseRepository
				logger.Info("Auto-detected ClickHouse configuration")
			} else {
				selectedType = client.PostgresRepository
				logger.Info("Auto-detected PostgreSQL configuration")
			}
		} else {
			// Convert string to LogRepositoryType
			selectedType = client.LogRepositoryType(repoType)
		}

		// Initialize the repository
		logRepo, err := client.NewLogRepository(selectedType, db)
		if err != nil {
			// If the selected type fails and we auto-detected, try fallback
			if repoType == "" && selectedType == client.ClickHouseRepository {
				logger.Warn("Failed to initialize ClickHouse repository, falling back to PostgreSQL", zap.Error(err))
				logRepo, err = client.NewLogRepository(client.PostgresRepository, db)
				if err != nil {
					logger.Error("Failed to initialize PostgreSQL repository for logs", zap.Error(err))
					initErr = fmt.Errorf("failed to initialize logs repository: %w", err)
					return
				}
				selectedType = client.PostgresRepository
			} else {
				logger.Error("Failed to initialize logs repository",
					zap.String("type", string(selectedType)),
					zap.Error(err))
				initErr = fmt.Errorf("failed to initialize %s repository: %w", selectedType, err)
				return
			}
		}

		instance = New(logRepo)
		logger.Info("Logs service initialized", zap.String("repository", string(selectedType)))
	})

	return instance, initErr
}

// ResetInstance resets the singleton instance. This should only be used in tests.
func ResetInstance() {
	once = sync.Once{}
	instance = nil
	initErr = nil
}
