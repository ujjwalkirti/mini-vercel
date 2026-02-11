# Logs Service Configuration

The logs service uses a **singleton pattern** to ensure only one instance is created and shared across the application. This prevents duplicate initialization and ensures consistent behavior.

## Repository Backends

The logs service supports two storage backends:

1. **ClickHouse** - High-performance columnar database optimized for logs
2. **PostgreSQL** - Relational database (fallback option)

## Configuration

### Option 1: Auto-Detection (Recommended)

By default, the service auto-detects which backend to use:
- If `CLICKHOUSE_HOST` is configured → Uses ClickHouse
- Otherwise → Falls back to PostgreSQL

**No additional configuration needed!**

### Option 2: Explicit Selection

Set the `LOGS_REPOSITORY_TYPE` environment variable to explicitly choose:

```bash
# Use ClickHouse
LOGS_REPOSITORY_TYPE=clickhouse

# Use PostgreSQL
LOGS_REPOSITORY_TYPE=postgres
```

## Usage in Code

### Getting the Singleton Instance

```go
import (
    "github.com/ujjwalkirti/mini-vercel-api-server/internal/config"
    "github.com/ujjwalkirti/mini-vercel-api-server/internal/service/logs"
)

// Auto-detect based on environment configuration
logsConfig := config.GetLogsConfig()
logService, err := logs.GetInstance(db, logsConfig.RepositoryType)
if err != nil {
    logger.Fatal("Failed to initialize logs service", zap.Error(err))
}

// Or use auto-detection directly
logService, err := logs.GetInstance(db, "")
```

### Explicit Backend Selection

```go
import "github.com/ujjwalkirti/mini-vercel-api-server/internal/service/logs"

// Force ClickHouse
logService, err := logs.GetInstance(db, "clickhouse")

// Force PostgreSQL
logService, err := logs.GetInstance(db, "postgres")
```

## How It Works

1. **First Call**: `GetInstance()` initializes the service with the specified or auto-detected backend
2. **Subsequent Calls**: Returns the same instance, regardless of parameters
3. **Thread-Safe**: Uses `sync.Once` to prevent race conditions

## Testing

For unit tests, you can reset the singleton:

```go
import "github.com/ujjwalkirti/mini-vercel-api-server/internal/service/logs"

func TestSomething(t *testing.T) {
    // Reset before test
    logs.ResetInstance()

    // Your test code
    logService, _ := logs.GetInstance(db, "postgres")

    // Clean up after test
    defer logs.ResetInstance()
}
```

## Benefits

✅ **Single Initialization** - Service created only once across the application
✅ **Thread-Safe** - Uses `sync.Once` to prevent race conditions
✅ **Flexible** - Supports multiple backends with easy switching
✅ **Auto-Fallback** - Automatically falls back to PostgreSQL if ClickHouse fails
✅ **Environment-Based** - Can be configured via environment variables
✅ **Testable** - Includes reset functionality for unit tests

## Environment Variables

See [.env.example](../.env.example) for all available configuration options.
