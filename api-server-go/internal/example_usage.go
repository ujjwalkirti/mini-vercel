package internal

import (
	"context"
	"time"

	"github.com/ujjwalkirti/mini-vercel-api-server/internal/client"
	logsservice "github.com/ujjwalkirti/mini-vercel-api-server/internal/service/logs"
	"github.com/ujjwalkirti/mini-vercel-api-server/internal/utils"
)

// ExampleLogsUsage demonstrates how to use the logs service
// This example shows how the architecture allows for easy replacement of storage backends
func ExampleLogsUsage() error {
	ctx := context.Background()

	// Step 1: Create a log repository using the factory
	// To switch to PostgreSQL in the future, simply change ClickHouseRepository to PostgresRepository
	logRepo, err := client.NewLogRepository(client.ClickHouseRepository)
	if err != nil {
		return err
	}

	// Step 2: Create the logs service with the repository
	logsService := logsservice.New(logRepo)

	// Step 3: Use the service to perform log operations

	// Insert a log
	err = logsService.InsertLog(ctx, logsservice.LogEvent{
		EventID:      utils.GenerateUUID(),
		DeploymentID: "deployment-123",
		Log:          "Application started successfully",
		Timestamp:    time.Now(),
	})
	if err != nil {
		return err
	}

	// Get all logs for a deployment
	logs, err := logsService.GetDeploymentLogs(ctx, "deployment-123")
	if err != nil {
		return err
	}
	_ = logs

	// Get logs with pagination
	logsWithLimit, err := logsService.GetDeploymentLogsWithLimit(ctx, "deployment-123", 10, 0)
	if err != nil {
		return err
	}
	_ = logsWithLimit

	// Get logs in a time range
	startTime := time.Now().Add(-24 * time.Hour)
	endTime := time.Now()
	logsInRange, err := logsService.GetLogsInTimeRange(ctx, "deployment-123", startTime, endTime)
	if err != nil {
		return err
	}
	_ = logsInRange

	// Get log count
	count, err := logsService.GetLogCount(ctx, "deployment-123")
	if err != nil {
		return err
	}
	_ = count

	return nil
}

// To switch to a different storage backend (e.g., PostgreSQL WAL):
// 1. Implement a PostgreSQL adapter that satisfies the logs.Repository interface
// 2. Add the implementation to the factory in client/factory.go
// 3. Change client.ClickHouseRepository to client.PostgresRepository in your initialization code
// No changes needed to the service layer or business logic!
