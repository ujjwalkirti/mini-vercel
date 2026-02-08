package logs

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"testing"
	"time"

	repologs "github.com/ujjwalkirti/mini-vercel-api-server/internal/repository/logs"
)

// MockRepository is a mock implementation of the logs repository
type MockRepository struct {
	QueryContextFunc    func(ctx context.Context, query string, args ...any) (repologs.Rows, error)
	QueryRowContextFunc func(ctx context.Context, query string, args ...any) repologs.Row
	ExecContextFunc     func(ctx context.Context, query string, args ...any) (repologs.Result, error)
}

func (m *MockRepository) QueryContext(ctx context.Context, query string, args ...any) (repologs.Rows, error) {
	if m.QueryContextFunc != nil {
		return m.QueryContextFunc(ctx, query, args...)
	}
	return nil, errors.New("QueryContextFunc not set")
}

func (m *MockRepository) QueryRowContext(ctx context.Context, query string, args ...any) repologs.Row {
	if m.QueryRowContextFunc != nil {
		return m.QueryRowContextFunc(ctx, query, args...)
	}
	return &MockRowImpl{err: errors.New("QueryRowContextFunc not set")}
}

func (m *MockRepository) ExecContext(ctx context.Context, query string, args ...any) (repologs.Result, error) {
	if m.ExecContextFunc != nil {
		return m.ExecContextFunc(ctx, query, args...)
	}
	return nil, errors.New("ExecContextFunc not set")
}

// MockRows implements the Rows interface
type MockRows interface {
	Close() error
	Next() bool
	Scan(dest ...any) error
	Err() error
}

type MockRowsImpl struct {
	rows    [][]any
	current int
	err     error
	closed  bool
}

func (m *MockRowsImpl) Close() error {
	m.closed = true
	return nil
}

func (m *MockRowsImpl) Next() bool {
	if m.current < len(m.rows) {
		m.current++
		return true
	}
	return false
}

func (m *MockRowsImpl) Scan(dest ...any) error {
	if m.err != nil {
		return m.err
	}
	if m.current == 0 || m.current > len(m.rows) {
		return errors.New("invalid row")
	}

	row := m.rows[m.current-1]
	for i, val := range row {
		if i < len(dest) {
			switch d := dest[i].(type) {
			case *string:
				*d = val.(string)
			case *time.Time:
				*d = val.(time.Time)
			case *int64:
				*d = val.(int64)
			}
		}
	}
	return nil
}

func (m *MockRowsImpl) Err() error {
	return m.err
}

// MockRow implements the Row interface
type MockRow interface {
	Scan(dest ...any) error
}

type MockRowImpl struct {
	values []any
	err    error
}

func (m *MockRowImpl) Scan(dest ...any) error {
	if m.err != nil {
		return m.err
	}
	for i, val := range m.values {
		if i < len(dest) {
			switch d := dest[i].(type) {
			case *int64:
				*d = val.(int64)
			}
		}
	}
	return nil
}

// MockResult implements the Result interface
type MockResult interface {
	LastInsertId() (int64, error)
	RowsAffected() (int64, error)
}

type MockResultImpl struct {
	lastInsertId int64
	rowsAffected int64
	err          error
}

func (m *MockResultImpl) LastInsertId() (int64, error) {
	return m.lastInsertId, m.err
}

func (m *MockResultImpl) RowsAffected() (int64, error) {
	return m.rowsAffected, m.err
}

// TestGetDeploymentLogs tests retrieving deployment logs
func TestGetDeploymentLogs(t *testing.T) {
	t.Run("should return logs for valid deployment ID", func(t *testing.T) {
		// Arrange
		deploymentID := "deploy-123"
		now := time.Now()

		mockRepo := &MockRepository{
			QueryContextFunc: func(ctx context.Context, query string, args ...any) (repologs.Rows, error) {
				return &MockRowsImpl{
					rows: [][]any{
						{"event-1", deploymentID, "Building project...", now},
						{"event-2", deploymentID, "Deploying...", now.Add(time.Second)},
					},
				}, nil
			},
		}

		service := New(mockRepo)

		// Act
		logs, err := service.GetDeploymentLogs(context.Background(), deploymentID)

		// Assert
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(logs) != 2 {
			t.Errorf("expected 2 logs, got %d", len(logs))
		}
		if logs[0].EventID != "event-1" {
			t.Errorf("expected event-1, got %s", logs[0].EventID)
		}
		if logs[0].Log != "Building project..." {
			t.Errorf("expected 'Building project...', got %s", logs[0].Log)
		}
	})

	t.Run("should return empty array when no logs exist", func(t *testing.T) {
		// Arrange
		mockRepo := &MockRepository{
			QueryContextFunc: func(ctx context.Context, query string, args ...any) (repologs.Rows, error) {
				return &MockRowsImpl{
					rows: [][]any{},
				}, nil
			},
		}

		service := New(mockRepo)

		// Act
		logs, err := service.GetDeploymentLogs(context.Background(), "deploy-123")

		// Assert
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		// logs can be nil or empty when there are no logs - both are acceptable
		if logs != nil && len(logs) != 0 {
			t.Errorf("expected 0 logs, got %d", len(logs))
		}
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		// Arrange
		mockRepo := &MockRepository{
			QueryContextFunc: func(ctx context.Context, query string, args ...any) (repologs.Rows, error) {
				return nil, errors.New("database error")
			},
		}

		service := New(mockRepo)

		// Act
		_, err := service.GetDeploymentLogs(context.Background(), "deploy-123")

		// Assert
		if err == nil {
			t.Error("expected error but got nil")
		}
	})

	t.Run("should handle invalid deployment ID", func(t *testing.T) {
		// Arrange
		mockRepo := &MockRepository{
			QueryContextFunc: func(ctx context.Context, query string, args ...any) (repologs.Rows, error) {
				return &MockRowsImpl{rows: [][]any{}}, nil
			},
		}

		service := New(mockRepo)

		// Act
		logs, err := service.GetDeploymentLogs(context.Background(), "invalid-id")

		// Assert
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(logs) != 0 {
			t.Errorf("expected 0 logs for invalid ID, got %d", len(logs))
		}
	})

	t.Run("should return logs in chronological order", func(t *testing.T) {
		// Arrange
		now := time.Now()
		time1 := now
		time2 := now.Add(time.Second)
		time3 := now.Add(2 * time.Second)

		mockRepo := &MockRepository{
			QueryContextFunc: func(ctx context.Context, query string, args ...any) (repologs.Rows, error) {
				return &MockRowsImpl{
					rows: [][]any{
						{"event-1", "deploy-123", "First log", time1},
						{"event-2", "deploy-123", "Second log", time2},
						{"event-3", "deploy-123", "Third log", time3},
					},
				}, nil
			},
		}

		service := New(mockRepo)

		// Act
		logs, err := service.GetDeploymentLogs(context.Background(), "deploy-123")

		// Assert
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(logs) != 3 {
			t.Fatalf("expected 3 logs, got %d", len(logs))
		}
		// Verify chronological order
		if !logs[0].Timestamp.Before(*logs[1].Timestamp) {
			t.Error("expected logs to be in chronological order")
		}
		if !logs[1].Timestamp.Before(*logs[2].Timestamp) {
			t.Error("expected logs to be in chronological order")
		}
	})

	t.Run("should handle context cancellation", func(t *testing.T) {
		// Arrange
		mockRepo := &MockRepository{
			QueryContextFunc: func(ctx context.Context, query string, args ...any) (repologs.Rows, error) {
				return nil, context.Canceled
			},
		}

		service := New(mockRepo)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// Act
		_, err := service.GetDeploymentLogs(ctx, "deploy-123")

		// Assert
		if err == nil {
			t.Error("expected error but got nil")
		}
	})

	t.Run("should include all log fields", func(t *testing.T) {
		// Arrange
		now := time.Now()
		mockRepo := &MockRepository{
			QueryContextFunc: func(ctx context.Context, query string, args ...any) (repologs.Rows, error) {
				return &MockRowsImpl{
					rows: [][]any{
						{"event-123", "deploy-456", "Test log message", now},
					},
				}, nil
			},
		}

		service := New(mockRepo)

		// Act
		logs, err := service.GetDeploymentLogs(context.Background(), "deploy-456")

		// Assert
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(logs) != 1 {
			t.Fatalf("expected 1 log, got %d", len(logs))
		}

		log := logs[0]
		if log.EventID != "event-123" {
			t.Errorf("expected EventID event-123, got %s", log.EventID)
		}
		if log.DeploymentID != "deploy-456" {
			t.Errorf("expected DeploymentID deploy-456, got %s", log.DeploymentID)
		}
		if log.Log != "Test log message" {
			t.Errorf("expected Log 'Test log message', got %s", log.Log)
		}
		if log.Timestamp == nil {
			t.Error("expected Timestamp to be set")
		}
	})
}

// TestLogEventStructure tests the LogEvent struct
func TestLogEventStructure(t *testing.T) {
	t.Run("should have required fields", func(t *testing.T) {
		// Arrange
		now := time.Now()
		log := LogEvent{
			EventID:      "event-123",
			DeploymentID: "deploy-456",
			Log:          "Test message",
			Timestamp:    &now,
		}

		// Assert
		if log.EventID == "" {
			t.Error("expected EventID to be set")
		}
		if log.DeploymentID == "" {
			t.Error("expected DeploymentID to be set")
		}
		if log.Log == "" {
			t.Error("expected Log to be set")
		}
		if log.Timestamp == nil {
			t.Error("expected Timestamp to be set")
		}
	})

	t.Run("should marshal to JSON correctly", func(t *testing.T) {
		// Arrange
		now := time.Now()
		log := LogEvent{
			EventID:      "event-123",
			DeploymentID: "deploy-456",
			Log:          "Test message",
			Timestamp:    &now,
		}

		// Act
		jsonData, err := json.Marshal(log)

		// Assert
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(jsonData) == 0 {
			t.Error("expected non-empty JSON")
		}

		// Verify JSON contains expected fields
		var result map[string]any
		if err := json.Unmarshal(jsonData, &result); err != nil {
			t.Errorf("failed to unmarshal JSON: %v", err)
		}
		if result["event_id"] != "event-123" {
			t.Error("expected event_id field in JSON")
		}
		if result["deployment_id"] != "deploy-456" {
			t.Error("expected deployment_id field in JSON")
		}
		if result["log"] != "Test message" {
			t.Error("expected log field in JSON")
		}
	})

	t.Run("should unmarshal from JSON correctly", func(t *testing.T) {
		// Arrange
		jsonStr := `{
			"event_id": "event-789",
			"deployment_id": "deploy-012",
			"log": "Another test message",
			"timestamp": "2024-01-01T12:00:00Z"
		}`

		// Act
		var log LogEvent
		err := json.Unmarshal([]byte(jsonStr), &log)

		// Assert
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if log.EventID != "event-789" {
			t.Errorf("expected EventID event-789, got %s", log.EventID)
		}
		if log.DeploymentID != "deploy-012" {
			t.Errorf("expected DeploymentID deploy-012, got %s", log.DeploymentID)
		}
		if log.Log != "Another test message" {
			t.Errorf("expected Log 'Another test message', got %s", log.Log)
		}
		if log.Timestamp == nil {
			t.Error("expected Timestamp to be set")
		}
	})
}

// TestInsertLog tests inserting a single log event
func TestInsertLog(t *testing.T) {
	t.Run("should insert log with all fields", func(t *testing.T) {
		// Arrange
		mockRepo := &MockRepository{
			ExecContextFunc: func(ctx context.Context, query string, args ...any) (repologs.Result, error) {
				if len(args) != 3 {
					t.Errorf("expected 3 args, got %d", len(args))
				}
				return &MockResultImpl{rowsAffected: 1}, nil
			},
		}

		service := New(mockRepo)
		log := LogEvent{
			EventID:      "event-123",
			DeploymentID: "deploy-456",
			Log:          "Test log",
		}

		// Act
		err := service.InsertLog(context.Background(), log)

		// Assert
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("should auto-generate event ID if empty", func(t *testing.T) {
		// Arrange
		var capturedEventID string
		mockRepo := &MockRepository{
			ExecContextFunc: func(ctx context.Context, query string, args ...any) (repologs.Result, error) {
				capturedEventID = args[0].(string)
				return &MockResultImpl{rowsAffected: 1}, nil
			},
		}

		service := New(mockRepo)
		log := LogEvent{
			EventID:      "", // Empty - should be generated
			DeploymentID: "deploy-456",
			Log:          "Test log",
		}

		// Act
		err := service.InsertLog(context.Background(), log)

		// Assert
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if capturedEventID == "" {
			t.Error("expected event ID to be generated")
		}
	})

	t.Run("should return error when insert fails", func(t *testing.T) {
		// Arrange
		mockRepo := &MockRepository{
			ExecContextFunc: func(ctx context.Context, query string, args ...any) (repologs.Result, error) {
				return nil, errors.New("insert failed")
			},
		}

		service := New(mockRepo)
		log := LogEvent{
			EventID:      "event-123",
			DeploymentID: "deploy-456",
			Log:          "Test log",
		}

		// Act
		err := service.InsertLog(context.Background(), log)

		// Assert
		if err == nil {
			t.Error("expected error but got nil")
		}
	})
}

// TestGetLogCount tests counting logs for a deployment
func TestGetLogCount(t *testing.T) {
	t.Run("should return correct count", func(t *testing.T) {
		// Arrange
		mockRepo := &MockRepository{
			QueryRowContextFunc: func(ctx context.Context, query string, args ...any) repologs.Row {
				return &MockRowImpl{
					values: []any{int64(42)},
				}
			},
		}

		service := New(mockRepo)

		// Act
		count, err := service.GetLogCount(context.Background(), "deploy-123")

		// Assert
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if count != 42 {
			t.Errorf("expected count 42, got %d", count)
		}
	})

	t.Run("should return 0 when no logs exist", func(t *testing.T) {
		// Arrange
		mockRepo := &MockRepository{
			QueryRowContextFunc: func(ctx context.Context, query string, args ...any) repologs.Row {
				return &MockRowImpl{
					values: []any{int64(0)},
				}
			},
		}

		service := New(mockRepo)

		// Act
		count, err := service.GetLogCount(context.Background(), "deploy-123")

		// Assert
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if count != 0 {
			t.Errorf("expected count 0, got %d", count)
		}
	})

	t.Run("should return error when query fails", func(t *testing.T) {
		// Arrange
		mockRepo := &MockRepository{
			QueryRowContextFunc: func(ctx context.Context, query string, args ...any) repologs.Row {
				return &MockRowImpl{
					err: sql.ErrConnDone,
				}
			},
		}

		service := New(mockRepo)

		// Act
		_, err := service.GetLogCount(context.Background(), "deploy-123")

		// Assert
		if err == nil {
			t.Error("expected error but got nil")
		}
	})
}

// TestInsert tests the convenience Insert method
func TestInsert(t *testing.T) {
	t.Run("should insert log with deployment ID and log text", func(t *testing.T) {
		// Arrange
		mockRepo := &MockRepository{
			ExecContextFunc: func(ctx context.Context, query string, args ...any) (repologs.Result, error) {
				if len(args) != 3 {
					t.Errorf("expected 3 args, got %d", len(args))
				}
				// args[0] = generated event_id, args[1] = deployment_id, args[2] = log
				if args[1].(string) != "deploy-456" {
					t.Errorf("expected deployment_id deploy-456, got %s", args[1])
				}
				if args[2].(string) != "Test log message" {
					t.Errorf("expected log 'Test log message', got %s", args[2])
				}
				return &MockResultImpl{rowsAffected: 1}, nil
			},
		}

		service := New(mockRepo)

		// Act
		err := service.Insert(context.Background(), "deploy-456", "Test log message")

		// Assert
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("should auto-generate event ID", func(t *testing.T) {
		// Arrange
		var capturedEventID string
		mockRepo := &MockRepository{
			ExecContextFunc: func(ctx context.Context, query string, args ...any) (repologs.Result, error) {
				capturedEventID = args[0].(string)
				return &MockResultImpl{rowsAffected: 1}, nil
			},
		}

		service := New(mockRepo)

		// Act
		err := service.Insert(context.Background(), "deploy-123", "Log message")

		// Assert
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if capturedEventID == "" {
			t.Error("expected event ID to be generated")
		}
	})
}

// TestDeleteDeploymentLogs tests deleting logs for a deployment
func TestDeleteDeploymentLogs(t *testing.T) {
	t.Run("should delete all logs for deployment", func(t *testing.T) {
		// Arrange
		mockRepo := &MockRepository{
			ExecContextFunc: func(ctx context.Context, query string, args ...any) (repologs.Result, error) {
				return &MockResultImpl{rowsAffected: 5}, nil
			},
		}

		service := New(mockRepo)

		// Act
		err := service.DeleteDeploymentLogs(context.Background(), "deploy-123")

		// Assert
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("should return error when delete fails", func(t *testing.T) {
		// Arrange
		mockRepo := &MockRepository{
			ExecContextFunc: func(ctx context.Context, query string, args ...any) (repologs.Result, error) {
				return nil, errors.New("delete failed")
			},
		}

		service := New(mockRepo)

		// Act
		err := service.DeleteDeploymentLogs(context.Background(), "deploy-123")

		// Assert
		if err == nil {
			t.Error("expected error but got nil")
		}
	})
}
