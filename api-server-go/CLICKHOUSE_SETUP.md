# ClickHouse Integration Setup Guide

This guide explains how to set up and use the ClickHouse integration for logging in the mini-vercel API server.

## Overview

The logs system uses a generic, technology-agnostic architecture that makes it easy to switch between different storage backends (ClickHouse, PostgreSQL, etc.) without changing business logic.

## Configuration

### Important: Environment Variable Loading

The config package uses **lazy initialization** (functions instead of variables) to ensure environment variables are read **after** the `.env` file is loaded. This prevents the common issue where config values are empty even though the `.env` file exists.

### Environment Variables

Add the following environment variables to your `.env` file:

```env
CLICKHOUSE_HOST=localhost
CLICKHOUSE_PORT=9000
CLICKHOUSE_DATABASE=logs
CLICKHOUSE_USERNAME=default
CLICKHOUSE_PASSWORD=your-password
```

### Database Setup

1. Install ClickHouse (if using Docker):

```bash
docker run -d \
  --name clickhouse-server \
  -p 9000:9000 \
  -p 8123:8123 \
  clickhouse/clickhouse-server
```

2. Create the database and table:

```bash
# Connect to ClickHouse
docker exec -it clickhouse-server clickhouse-client

# Run the schema (or execute clickhouse_schema.sql)
CREATE DATABASE IF NOT EXISTS logs;

USE logs;

CREATE TABLE IF NOT EXISTS log_events (
    event_id String,
    deployment_id String,
    log String,
    timestamp DateTime64(3)
) ENGINE = MergeTree()
ORDER BY (deployment_id, timestamp)
PARTITION BY toYYYYMM(timestamp)
TTL timestamp + INTERVAL 30 DAY;
```

## Usage

### Basic Usage

The logs service is automatically initialized in the deployment routes. You don't need to do anything special to use it.

```go
// The logs service is already available in deployment handlers
logs, err := logsService.GetDeploymentLogs(ctx, deploymentID)
```

### Available Operations

The logs service provides the following operations:

1. **Get all logs for a deployment**
   ```go
   logs, err := logsService.GetDeploymentLogs(ctx, deploymentID)
   ```

2. **Get logs with pagination**
   ```go
   logs, err := logsService.GetDeploymentLogsWithLimit(ctx, deploymentID, limit, offset)
   ```

3. **Get logs in a time range**
   ```go
   logs, err := logsService.GetLogsInTimeRange(ctx, deploymentID, startTime, endTime)
   ```

4. **Insert a log**
   ```go
   err := logsService.InsertLog(ctx, logsservice.LogEvent{
       EventID:      utils.GenerateUUID(),
       DeploymentID: "deployment-123",
       Log:          "Application started",
       Timestamp:    time.Now(),
   })
   ```

5. **Get log count**
   ```go
   count, err := logsService.GetLogCount(ctx, deploymentID)
   ```

6. **Delete deployment logs**
   ```go
   err := logsService.DeleteDeploymentLogs(ctx, deploymentID)
   ```

## Architecture

The system follows a layered architecture:

```
Handler → Service → Repository Interface → Adapter → ClickHouse Client
```

- **Service Layer**: Contains business logic, independent of storage
- **Repository Interface**: Defines storage operations
- **Adapter**: Translates between ClickHouse driver and repository interface
- **Factory**: Creates the appropriate repository implementation

## Switching to a Different Backend

To switch from ClickHouse to another backend (e.g., PostgreSQL):

1. Create a new adapter implementing `logs.Repository`
2. Add it to the factory in `internal/client/factory.go`
3. Change `client.ClickHouseRepository` to your new repository type

No changes needed to:
- Service layer
- Handler layer
- Business logic

## Testing

A mock repository is available for testing:

```go
import "github.com/ujjwalkirti/mini-vercel-api-server/internal/repository/logs"

mockRepo := logs.NewMockRepository()
mockRepo.QueryContextFunc = func(ctx context.Context, query string, args ...any) (logs.Rows, error) {
    // Your mock implementation
}

logsService := logsservice.New(mockRepo)
```

## Files Changed/Created

### New Files
- `internal/client/clickhouse.go` - ClickHouse client with config integration
- `internal/client/clickhouse_adapter.go` - Adapter implementing Repository interface
- `internal/client/factory.go` - Factory for creating repositories
- `internal/repository/logs/repository.go` - Repository interface definition
- `internal/repository/logs/mock_repository.go` - Mock for testing
- `internal/utils/uuid.go` - UUID generation utility (updated)

### Modified Files
- `internal/service/logs/service.go` - Updated to use repository interface
- `internal/handler/deployment/routes.go` - Updated to use factory pattern
- `internal/config/clickhouse.go` - ClickHouse configuration
- `.env.example` - Updated with ClickHouse configuration

## Troubleshooting

### Connection Issues

If you see "Failed to initialize ClickHouse client":
- Check that ClickHouse is running
- Verify environment variables are set correctly
- Ensure the database exists
- Check network connectivity

### Empty Logs

If logs endpoint returns empty:
- Verify the table exists and has data
- Check that deployment_id matches
- Ensure timestamps are correct (ClickHouse uses UTC)

## Performance Considerations

- ClickHouse is optimized for analytical queries
- Logs are partitioned by month for better query performance
- TTL of 30 days automatically deletes old logs
- Use pagination for large result sets
