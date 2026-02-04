# Testing Guide for Mini Vercel API Server

This document provides guidelines for implementing tests in the API server.

## Test Structure

The test suite is organized to mirror the application structure:

```
api-server-go/
├── internal/
│   ├── handler/
│   │   ├── health/
│   │   │   └── handler_test.go
│   │   ├── project/
│   │   │   └── handler_test.go
│   │   └── deployment/
│   │       └── handler_test.go
│   ├── middleware/
│   │   └── auth_test.go
│   ├── router/
│   │   └── router_test.go
│   ├── service/
│   │   ├── ecs/
│   │   │   └── service_test.go
│   │   └── logs/
│   │       └── service_test.go
│   └── testutil/
│       ├── helpers.go      # Test helper functions
│       └── mocks.go        # Mock implementations
```

## Test Categories

### 1. Handler Tests (`internal/handler/*/handler_test.go`)

Test HTTP handlers with focus on:
- Request validation (missing fields, invalid data types, malformed JSON)
- Authentication and authorization (401 for unauthenticated, 403 for unauthorized)
- Business logic (correct responses for valid inputs)
- Error handling (404 for not found, 500 for server errors)
- Response format (JSON structure, status codes, headers)

**Key scenarios to test:**
- Unauthenticated requests → 401
- Invalid UUIDs → 400
- Missing required fields → 400
- Resource not found → 404
- User doesn't own resource → 404
- Database errors → 500
- Successful operations → 200/201

### 2. Middleware Tests (`internal/middleware/*_test.go`)

Test middleware functionality:
- Authentication middleware (token validation, user context injection)
- CORS configuration
- Request logging
- Error recovery

**Key scenarios to test:**
- Missing Authorization header → 401
- Invalid token format → 401
- Expired token → 401
- Invalid signature → 401
- Valid token → add user to context
- Role-based access control

### 3. Router Tests (`internal/router/router_test.go`)

Test router configuration and integration:
- Route mounting (correct paths)
- Middleware order and execution
- CORS configuration
- Public vs protected routes

### 4. Service Tests (`internal/service/*/service_test.go`)

Test business logic and external service integration:
- ECS task execution
- Log retrieval from ClickHouse
- Error handling for external service failures

### 5. Integration Tests

Test complete request flow from router to handler to service.

## Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...

# Run tests in a specific package
go test ./internal/handler/project/...

# Run a specific test
go test -run TestGetProjects ./internal/handler/project/...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Test Implementation Guidelines

### Using Test Helpers

The `testutil` package provides helpers for common testing tasks:

```go
import "github.com/ujjwalkirti/mini-vercel-api-server/internal/testutil"

// Create authenticated request
req := testutil.CreateAuthenticatedRequest("GET", "/projects", nil, userID, email, role)

// Parse JSON response
data := testutil.ParseJSONResponse(t, rr)

// Assert status code
testutil.AssertStatusCode(t, rr, http.StatusOK)

// Assert JSON field
testutil.AssertJSONField(t, data, "status", "ok")
```

### Using Mocks

Mock implementations are provided in `testutil/mocks.go`:

```go
// Create mock repository
mockRepo := &testutil.MockProjectRepository{
    ListByUserFunc: func(ctx context.Context, userID string) ([]*project.Project, error) {
        return testutil.CreateMockProjects(3, userID), nil
    },
}
```

### Test Naming Convention

Use descriptive test names that explain the scenario:

```go
func TestGetProject(t *testing.T) {
    t.Run("should return 401 when user is not authenticated", func(t *testing.T) {
        // Test implementation
    })

    t.Run("should return project when valid and user owns it", func(t *testing.T) {
        // Test implementation
    })
}
```

### Test Structure Pattern

Follow the Arrange-Act-Assert pattern:

```go
t.Run("description", func(t *testing.T) {
    // Arrange - Set up test data and mocks
    mockRepo := createMockRepository()
    handler := NewHandler(mockRepo)
    req := createTestRequest()

    // Act - Execute the code being tested
    rr := httptest.NewRecorder()
    handler.ServeHTTP(rr, req)

    // Assert - Verify the results
    if rr.Code != http.StatusOK {
        t.Errorf("Expected status 200, got %d", rr.Code)
    }
})
```

## What to Test

### ✅ Do Test

- All API endpoints (happy path and error cases)
- Input validation
- Authentication and authorization
- Error handling
- Business logic
- Edge cases (empty arrays, null values, etc.)

### ❌ Don't Test

- Third-party library internals
- Database driver functionality
- Standard library functions

## Mocking External Dependencies

### Database

Use mock repositories that implement the repository interface:

```go
type MockProjectRepository struct {
    GetByIDFunc func(ctx context.Context, id string) (*project.Project, error)
}

func (m *MockProjectRepository) GetByID(ctx context.Context, id string) (*project.Project, error) {
    if m.GetByIDFunc != nil {
        return m.GetByIDFunc(ctx, id)
    }
    return nil, errors.New("not implemented")
}
```

### AWS ECS

Mock the ECS service for testing deployment handlers:

```go
mockECS := &testutil.MockECSService{
    RunTaskFunc: func(ctx context.Context, envVars []ecs.EnvVar) (string, error) {
        return "arn:aws:ecs:region:account:task/cluster/task-id", nil
    },
}
```

### ClickHouse (Logs)

Mock the logs service:

```go
mockLogs := &testutil.MockLogsService{
    GetDeploymentLogsFunc: func(ctx context.Context, deploymentID string) ([]logs.LogEvent, error) {
        return testutil.CreateTestLogEvents(deploymentID, 5), nil
    },
}
```

### JWKS/JWT Authentication

Mock JWKS cache for testing authentication:

```go
mockJWKS := &testutil.MockJWKSCache{
    GetKeyFunc: func(keyID string) (interface{}, error) {
        // Return test public key
        return testPublicKey, nil
    },
}
```

## Common Test Scenarios

### Testing Authenticated Endpoints

```go
// Create authenticated context
ctx := testutil.CreateTestContext("user-id-123", "user@example.com", "user")
req := httptest.NewRequest("GET", "/projects", nil).WithContext(ctx)

// Or create request with Bearer token
req := testutil.CreateRequestWithBearerToken("GET", "/projects", validToken, nil)
```

### Testing Error Responses

```go
testutil.AssertErrorResponse(t, rr, http.StatusNotFound, "Project not found")
```

### Testing JSON Responses

```go
data := testutil.ParseJSONResponse(t, rr)
testutil.AssertJSONFieldExists(t, data, "id")
testutil.AssertJSONField(t, data, "name", "My Project")
```

## Coverage Goals

- Aim for 80%+ code coverage
- Prioritize critical paths (authentication, data mutation)
- Focus on meaningful tests over coverage metrics

## Resources

- [Go Testing Package](https://pkg.go.dev/testing)
- [httptest Package](https://pkg.go.dev/net/http/httptest)
- [Table-Driven Tests in Go](https://dave.cheney.net/2019/05/07/prefer-table-driven-tests)

## TODO: Learning Path

As you learn Go testing, implement tests in this order:

1. **Start Simple:** Health handler tests (no auth, simple responses)
2. **Add Mocking:** Project handler tests (mock repositories)
3. **Authentication:** Auth middleware tests (JWT validation)
4. **Complex Flows:** Deployment handler tests (multiple dependencies)
5. **Integration:** Router tests (full request flow)
6. **Services:** ECS and Logs service tests (external services)

Each test file has TODO comments indicating what needs to be implemented. Work through them gradually as you learn testing concepts.
