package logs

import (
	"context"
	"fmt"
	"time"

	repologs "github.com/ujjwalkirti/mini-vercel-api-server/internal/repository/logs"
)

// LogEvent represents a log entry
type LogEvent struct {
	EventID      string     `json:"event_id"`
	DeploymentID string     `json:"deployment_id"`
	Log          string     `json:"log"`
	Timestamp    *time.Time `json:"timestamp,omitempty"`
}

// Service handles business logic for deployment logs
type Service struct {
	repo repologs.Repository
}

// New creates a new logs service
func New(repo repologs.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// GetDeploymentLogs retrieves all logs for a specific deployment
func (s *Service) GetDeploymentLogs(ctx context.Context, deploymentID string) ([]LogEvent, error) {
	query := `
		SELECT event_id, deployment_id, log, timestamp
		FROM log_events
		WHERE deployment_id = ?
		ORDER BY timestamp ASC
	`

	rows, err := s.repo.QueryContext(ctx, query, deploymentID)
	if err != nil {
		return nil, fmt.Errorf("failed to query logs: %w", err)
	}
	defer rows.Close()

	var logs []LogEvent
	for rows.Next() {
		var log LogEvent
		var timestamp time.Time
		if err := rows.Scan(&log.EventID, &log.DeploymentID, &log.Log, &timestamp); err != nil {
			return nil, fmt.Errorf("failed to scan log row: %w", err)
		}
		log.Timestamp = &timestamp
		logs = append(logs, log)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating logs: %w", err)
	}

	return logs, nil
}

// GetDeploymentLogsWithLimit retrieves logs with pagination
func (s *Service) GetDeploymentLogsWithLimit(ctx context.Context, deploymentID string, limit, offset int) ([]LogEvent, error) {
	query := `
		SELECT event_id, deployment_id, log, timestamp
		FROM log_events
		WHERE deployment_id = ?
		ORDER BY timestamp ASC
		LIMIT ? OFFSET ?
	`

	rows, err := s.repo.QueryContext(ctx, query, deploymentID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query logs: %w", err)
	}
	defer rows.Close()

	var logs []LogEvent
	for rows.Next() {
		var log LogEvent
		var timestamp time.Time
		if err := rows.Scan(&log.EventID, &log.DeploymentID, &log.Log, &timestamp); err != nil {
			return nil, fmt.Errorf("failed to scan log row: %w", err)
		}
		log.Timestamp = &timestamp
		logs = append(logs, log)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating logs: %w", err)
	}

	return logs, nil
}

// GetLogsInTimeRange retrieves logs within a specific time range
func (s *Service) GetLogsInTimeRange(ctx context.Context, deploymentID string, startTime, endTime time.Time) ([]LogEvent, error) {
	query := `
		SELECT event_id, deployment_id, log, timestamp
		FROM log_events
		WHERE deployment_id = ?
		  AND timestamp >= ?
		  AND timestamp <= ?
		ORDER BY timestamp ASC
	`

	rows, err := s.repo.QueryContext(ctx, query, deploymentID, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to query logs: %w", err)
	}
	defer rows.Close()

	var logs []LogEvent
	for rows.Next() {
		var log LogEvent
		var timestamp time.Time
		if err := rows.Scan(&log.EventID, &log.DeploymentID, &log.Log, &timestamp); err != nil {
			return nil, fmt.Errorf("failed to scan log row: %w", err)
		}
		log.Timestamp = &timestamp
		logs = append(logs, log)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating logs: %w", err)
	}

	return logs, nil
}

// InsertLog inserts a single log event
func (s *Service) InsertLog(ctx context.Context, log LogEvent) error {
	query := `
		INSERT INTO log_events (event_id, deployment_id, log, timestamp)
		VALUES (?, ?, ?, ?)
	`

	var timestamp time.Time
	if log.Timestamp != nil {
		timestamp = *log.Timestamp
	} else {
		timestamp = time.Now()
	}

	_, err := s.repo.ExecContext(ctx, query, log.EventID, log.DeploymentID, log.Log, timestamp)
	if err != nil {
		return fmt.Errorf("failed to insert log: %w", err)
	}

	return nil
}

// GetLogCount returns the total number of logs for a deployment
func (s *Service) GetLogCount(ctx context.Context, deploymentID string) (int64, error) {
	query := `
		SELECT count(*)
		FROM log_events
		WHERE deployment_id = ?
	`

	var count int64
	err := s.repo.QueryRowContext(ctx, query, deploymentID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count logs: %w", err)
	}

	return count, nil
}

// DeleteDeploymentLogs deletes all logs for a specific deployment
func (s *Service) DeleteDeploymentLogs(ctx context.Context, deploymentID string) error {
	query := `
		ALTER TABLE log_events DELETE WHERE deployment_id = ?
	`

	_, err := s.repo.ExecContext(ctx, query, deploymentID)
	if err != nil {
		return fmt.Errorf("failed to delete logs: %w", err)
	}

	return nil
}

// Insert is a simpler method that inserts a log with just deployment_id and log text
// Similar to the Express API server's insertLog method
func (s *Service) Insert(ctx context.Context, deploymentID, log string) error {
	query := `
		INSERT INTO log_events (event_id, deployment_id, log, timestamp)
		VALUES (?, ?, ?, ?)
	`

	eventID := generateEventID()
	timestamp := time.Now()

	_, err := s.repo.ExecContext(ctx, query, eventID, deploymentID, log, timestamp)
	if err != nil {
		return fmt.Errorf("failed to insert log: %w", err)
	}

	return nil
}

// generateEventID generates a simple event ID
func generateEventID() string {
	return fmt.Sprintf("%d-%d", time.Now().UnixNano(), time.Now().Unix())
}
