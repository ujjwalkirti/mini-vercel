# Logs Architecture

This document explains the architecture of the logs system, which is designed to be storage-backend agnostic.

## Architecture Overview

The logs system follows a layered architecture with clear separation of concerns:

```
┌─────────────────────────────────────────┐
│         Handler Layer                   │
│  (HTTP handlers, API endpoints)         │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│         Service Layer                   │
│  (Business logic, domain operations)    │
│  internal/service/logs/                 │
└─────────────────┬───────────────────────┘
                  │
                  │ Uses Repository Interface
                  │
┌─────────────────▼───────────────────────┐
│      Repository Interface               │
│  (Storage abstraction)                  │
│  internal/repository/logs/              │
└─────────────────┬───────────────────────┘
                  │
                  │ Implemented by
                  │
┌─────────────────▼───────────────────────┐
│         Client Layer                    │
│  (ClickHouse, PostgreSQL, etc.)         │
│  internal/client/                       │
└─────────────────────────────────────────┘
```

## Components

### 1. Repository Interface (`internal/repository/logs/repository.go`)

Defines a generic interface for log storage operations:

```go
type Repository interface {
    QueryContext(ctx context.Context, query string, args ...any) (Rows, error)
    QueryRowContext(ctx context.Context, query string, args ...any) Row
    ExecContext(ctx context.Context, query string, args ...any) (Result, error)
}
```

This interface is technology-agnostic and can be implemented by any storage backend.

### 2. Service Layer (`internal/service/logs/service.go`)

Contains business logic for log operations:
- `GetDeploymentLogs()` - Get all logs for a deployment
- `GetDeploymentLogsWithLimit()` - Get logs with pagination
- `GetLogsInTimeRange()` - Get logs within a time range
- `InsertLog()` - Insert a new log event
- `GetLogCount()` - Count logs for a deployment
- `DeleteDeploymentLogs()` - Delete all logs for a deployment

The service depends only on the `Repository` interface, not on any specific implementation.

### 3. ClickHouse Client (`internal/client/clickhouse.go`)

Manages the ClickHouse connection:
- Reads configuration from `internal/config/clickhouse.go`
- Establishes and maintains connection
- Provides connection access to adapters

### 4. ClickHouse Adapter (`internal/client/clickhouse_adapter.go`)

Adapts the ClickHouse driver to implement the `Repository` interface:
- Wraps ClickHouse-specific types
- Implements `QueryContext`, `QueryRowContext`, `ExecContext`
- Translates between generic interface and ClickHouse driver

### 5. Factory (`internal/client/factory.go`)

Provides a factory function to create repositories:

```go
func NewLogRepository(repoType LogRepositoryType) (logs.Repository, error)
```

Currently supports:
- `ClickHouseRepository` - ClickHouse implementation
- `PostgresRepository` - Placeholder for future PostgreSQL implementation

## Configuration

ClickHouse configuration is managed in `internal/config/clickhouse.go`:

```go
type ClickHouseConfig struct {
    Host     string
    Port     string
    Database string
    Username string
    Password string
}

func GetClickHouseConfig() ClickHouseConfig {
    return ClickHouseConfig{
        Host:     os.Getenv("CLICKHOUSE_HOST"),
        Port:     os.Getenv("CLICKHOUSE_PORT"),
        Database: os.Getenv("CLICKHOUSE_DATABASE"),
        Username: os.Getenv("CLICKHOUSE_USERNAME"),
        Password: os.Getenv("CLICKHOUSE_PASSWORD"),
    }
}
```

The function-based approach ensures environment variables are read after the `.env` file is loaded, preventing empty config values.

Set these environment variables in your `.env` file:

```env
CLICKHOUSE_HOST=localhost
CLICKHOUSE_PORT=9000
CLICKHOUSE_DATABASE=logs
CLICKHOUSE_USERNAME=default
CLICKHOUSE_PASSWORD=your-password
```

## Usage Example

```go
// 1. Create a repository using the factory
logRepo, err := client.NewLogRepository(client.ClickHouseRepository)
if err != nil {
    return err
}

// 2. Create the logs service
logsService := logsservice.New(logRepo)

// 3. Use the service
logs, err := logsService.GetDeploymentLogs(ctx, "deployment-123")
if err != nil {
    return err
}
```

## How to Switch Storage Backends

The architecture makes it easy to switch from ClickHouse to another storage backend (e.g., PostgreSQL WAL):

### Step 1: Implement the Repository Interface

Create a new adapter (e.g., `postgres_adapter.go`):

```go
type PostgresAdapter struct {
    db *sql.DB
}

func (a *PostgresAdapter) QueryContext(ctx context.Context, query string, args ...any) (Rows, error) {
    // PostgreSQL implementation
}

func (a *PostgresAdapter) QueryRowContext(ctx context.Context, query string, args ...any) Row {
    // PostgreSQL implementation
}

func (a *PostgresAdapter) ExecContext(ctx context.Context, query string, args ...any) (Result, error) {
    // PostgreSQL implementation
}
```

### Step 2: Update the Factory

Add the new implementation to `factory.go`:

```go
case PostgresRepository:
    client, err := NewPostgresClient()
    if err != nil {
        return nil, fmt.Errorf("failed to create postgres client: %w", err)
    }
    return NewPostgresAdapter(client.GetDB()), nil
```

### Step 3: Update Your Initialization Code

Change the repository type:

```go
// Before
logRepo, err := client.NewLogRepository(client.ClickHouseRepository)

// After
logRepo, err := client.NewLogRepository(client.PostgresRepository)
```

**That's it!** No changes needed to:
- Service layer (`internal/service/logs/`)
- Business logic
- Handler layer
- API endpoints

## Benefits of This Architecture

1. **Separation of Concerns**: Business logic is separated from storage implementation
2. **Testability**: Easy to mock the repository interface for testing
3. **Flexibility**: Switch storage backends with minimal code changes
4. **Maintainability**: Changes to storage don't affect business logic
5. **Type Safety**: Compile-time guarantees through interfaces
6. **Single Responsibility**: Each component has one clear purpose

## Database Schema

The ClickHouse schema is defined in `clickhouse_schema.sql`:

```sql
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

This schema is optimized for ClickHouse but the same data model can be used with other databases.
