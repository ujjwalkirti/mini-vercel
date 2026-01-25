package logs

import (
	"context"
	"database/sql"
	"sync"
)

// MockRepository is a mock implementation of Repository for testing
type MockRepository struct {
	mu sync.RWMutex

	// Mock data storage
	QueryContextFunc    func(ctx context.Context, query string, args ...any) (Rows, error)
	QueryRowContextFunc func(ctx context.Context, query string, args ...any) Row
	ExecContextFunc     func(ctx context.Context, query string, args ...any) (Result, error)

	// Call tracking
	QueryContextCalls    []MockCall
	QueryRowContextCalls []MockCall
	ExecContextCalls     []MockCall
}

// MockCall represents a tracked function call
type MockCall struct {
	Query string
	Args  []any
}

// NewMockRepository creates a new mock repository
func NewMockRepository() *MockRepository {
	return &MockRepository{
		QueryContextCalls:    make([]MockCall, 0),
		QueryRowContextCalls: make([]MockCall, 0),
		ExecContextCalls:     make([]MockCall, 0),
	}
}

// QueryContext implements Repository interface
func (m *MockRepository) QueryContext(ctx context.Context, query string, args ...any) (Rows, error) {
	m.mu.Lock()
	m.QueryContextCalls = append(m.QueryContextCalls, MockCall{Query: query, Args: args})
	m.mu.Unlock()

	if m.QueryContextFunc != nil {
		return m.QueryContextFunc(ctx, query, args...)
	}
	return &MockRows{}, nil
}

// QueryRowContext implements Repository interface
func (m *MockRepository) QueryRowContext(ctx context.Context, query string, args ...any) Row {
	m.mu.Lock()
	m.QueryRowContextCalls = append(m.QueryRowContextCalls, MockCall{Query: query, Args: args})
	m.mu.Unlock()

	if m.QueryRowContextFunc != nil {
		return m.QueryRowContextFunc(ctx, query, args...)
	}
	return &MockRow{}
}

// ExecContext implements Repository interface
func (m *MockRepository) ExecContext(ctx context.Context, query string, args ...any) (Result, error) {
	m.mu.Lock()
	m.ExecContextCalls = append(m.ExecContextCalls, MockCall{Query: query, Args: args})
	m.mu.Unlock()

	if m.ExecContextFunc != nil {
		return m.ExecContextFunc(ctx, query, args...)
	}
	return &MockResult{}, nil
}

// Reset clears all tracked calls
func (m *MockRepository) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.QueryContextCalls = make([]MockCall, 0)
	m.QueryRowContextCalls = make([]MockCall, 0)
	m.ExecContextCalls = make([]MockCall, 0)
}

// MockRows is a mock implementation of Rows
type MockRows struct {
	data    [][]any
	current int
	err     error
}

func (m *MockRows) Close() error {
	return nil
}

func (m *MockRows) Next() bool {
	if m.current >= len(m.data) {
		return false
	}
	m.current++
	return true
}

func (m *MockRows) Scan(dest ...any) error {
	if m.current == 0 || m.current > len(m.data) {
		return sql.ErrNoRows
	}
	row := m.data[m.current-1]
	for i, v := range row {
		if i < len(dest) {
			switch d := dest[i].(type) {
			case *string:
				if s, ok := v.(string); ok {
					*d = s
				}
			case *int64:
				if n, ok := v.(int64); ok {
					*d = n
				}
			}
		}
	}
	return nil
}

func (m *MockRows) Err() error {
	return m.err
}

// MockRow is a mock implementation of Row
type MockRow struct {
	data []any
	err  error
}

func (m *MockRow) Scan(dest ...any) error {
	if m.err != nil {
		return m.err
	}
	for i, v := range m.data {
		if i < len(dest) {
			switch d := dest[i].(type) {
			case *string:
				if s, ok := v.(string); ok {
					*d = s
				}
			case *int64:
				if n, ok := v.(int64); ok {
					*d = n
				}
			}
		}
	}
	return nil
}

// MockResult is a mock implementation of Result
type MockResult struct {
	lastInsertId int64
	rowsAffected int64
}

func (m *MockResult) LastInsertId() (int64, error) {
	return m.lastInsertId, nil
}

func (m *MockResult) RowsAffected() (int64, error) {
	return m.rowsAffected, nil
}
