# Kafka Consumer Implementation - Completed

This document summarizes the completed Kafka consumer implementation in the Go API server, mirroring the Express API server's functionality.

## Overview

The Kafka consumer implementation follows the same pattern as the Express server:
1. Listens for messages from the `build-events` topic (previously `mini-vercel-build-logs`)
2. Processes build log events
3. Updates deployment status based on specific log messages
4. Inserts all logs into ClickHouse

## Architecture

```
┌─────────────────────────────────────────┐
│         Application Layer               │
│  (internal/app/app.go)                  │
│  - Initializes all services             │
│  - Starts Kafka consumer                │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│      Kafka Consumer Group               │
│  (internal/kafka/consumer/group.go)     │
│  - Manages consumer group               │
│  - Subscribes to topics                 │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│         Handler                         │
│  (internal/kafka/consumer/handler.go)   │
│  - Implements Sarama interfaces         │
│  - Uses worker pool for concurrency     │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│         Processor                       │
│  (internal/kafka/consumer/process.go)   │
│  - Processes individual messages        │
│  - Updates deployment status            │
│  - Inserts logs to ClickHouse           │
└─────────────────┬───────────────────────┘
                  │
        ┌─────────┴─────────┐
        │                   │
        ▼                   ▼
┌───────────────┐   ┌───────────────┐
│  Deployment   │   │  Logs         │
│  Service      │   │  Service      │
└───────────────┘   └───────────────┘
```

## Implementation Files

### 1. Consumer Group (`internal/kafka/consumer/group.go`)
- Creates Sarama consumer group with configuration
- Subscribes to `build-events` topic
- Consumer group ID: `mini-vercel-build-logs`
- Runs in infinite loop with context cancellation support

### 2. Handler (`internal/kafka/consumer/handler.go`)
- Implements `sarama.ConsumerGroupHandler` interface
- Uses worker pool with 50 workers for concurrent message processing
- Marks messages after successful processing

### 3. Processor (`internal/kafka/consumer/process.go`)
- Contains business logic for processing messages
- Parses JSON events with structure:
  ```go
  {
    "project_id": "string",
    "deployment_id": "string",
    "log": "string"
  }
  ```
- Updates deployment status based on log content:
  - `"info: starting build pipeline..."` → Status: `IN_PROGRESS`
  - `"info: pipeline completed successfully."` → Status: `READY`
  - `"error: ... pipeline failed ..."` → Status: `FAIL`
- Always inserts log to ClickHouse

### 4. Worker Pool (`internal/kafka/consumer/worker_pool.go`)
- Semaphore-based worker pool
- Limits concurrent message processing to 50
- Prevents overwhelming the system

### 5. Logs Service (`internal/service/logs/service.go`)
- Added `Insert(ctx, deploymentID, log)` method
- Generates event IDs automatically
- Sets timestamp to current time
- Inserts to ClickHouse via repository pattern

### 6. Application Wiring (`internal/app/app.go`)
- Initializes ClickHouse repository
- Creates logs service and deployment service
- Instantiates Kafka processor with services
- Starts Kafka consumer in background goroutine
- Handles graceful shutdown with context cancellation

## Configuration

Kafka configuration is loaded from environment variables via `config.LoadKafkaConfig()`:

```env
KAFKA_BROKERS=broker1:9092,broker2:9092
KAFKA_CLIENT_ID=api-server
KAFKA_USERNAME=your-username
KAFKA_PASSWORD=your-password
```

TLS configuration is managed by `internal/kafka/tls/tls.go`:
- Development: uses `ca.pem` in the same directory
- Production: uses `/secrets/kafka-consumer-ca`

## Message Flow

1. Message arrives on `build-events` topic
2. Handler receives message via `ConsumeClaim`
3. Handler submits to worker pool
4. Processor parses JSON event
5. Processor checks log content and updates deployment status if needed
6. Processor inserts log to ClickHouse
7. Handler marks message as processed

## Differences from Express Implementation

| Aspect | Express (TypeScript) | Go |
|--------|---------------------|-----|
| Kafka Library | kafkajs | IBM Sarama |
| Topic Name | mini-vercel-build-logs | build-events |
| Concurrency | Sequential processing in batch | Worker pool (50 concurrent) |
| Error Handling | Logs errors, continues | Returns errors, marks only on success |
| Event ID Generation | uuid v4 | Timestamp-based |

## Testing the Implementation

1. Start the API server:
   ```bash
   cd api-server-go
   go run main.go
   ```

2. Send a test message to Kafka:
   ```bash
   # Using kafkacat or similar tool
   echo '{"project_id":"test-project","deployment_id":"test-deployment","log":"INFO: Starting build pipeline..."}' | \
     kafkacat -b localhost:9092 -t build-events -P
   ```

3. Check logs for processing confirmation

4. Verify in ClickHouse:
   ```sql
   SELECT * FROM log_events WHERE deployment_id = 'test-deployment';
   ```

## Graceful Shutdown

The implementation supports graceful shutdown:
1. Signal handler listens for `SIGINT` and `SIGTERM`
2. Context is cancelled
3. Consumer loop exits
4. In-flight messages complete processing

## Future Improvements

- Add metrics/observability (processed messages, errors, latency)
- Add retry logic for failed message processing
- Add dead letter queue for permanently failed messages
- Use crypto/rand for better event ID generation
- Add configurable worker pool size
- Add message schema validation
