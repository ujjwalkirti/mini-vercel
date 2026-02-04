# Quick Start: Your First Test

This guide will walk you through implementing and running your first test in the API server.

## Step 1: Choose Your First Test

We recommend starting with the health handler test since it's the simplest and has no external dependencies.

**File to edit:** `internal/handler/health/handler_test.go`

## Step 2: Implement a Simple Test

Let's implement the first test case. Open `internal/handler/health/handler_test.go` and replace the first test with:

```go
func TestHealthCheck(t *testing.T) {
	t.Run("should return 200 with status ok", func(t *testing.T) {
		// Arrange - Set up the test
		handler := NewHandler()
		req := httptest.NewRequest("GET", "/health", nil)
		rr := httptest.NewRecorder()

		// Act - Execute the test
		handler.Check(rr, req)

		// Assert - Verify the results
		if rr.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
		}

		// Parse the JSON response
		var response HealthResponse
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		// Check the status field
		if response.Status != "ok" {
			t.Errorf("Expected status 'ok', got '%s'", response.Status)
		}

		// Check that timestamp is present
		if response.Timestamp.IsZero() {
			t.Error("Expected timestamp to be set")
		}

		// Check that uptime is positive
		if response.Uptime <= 0 {
			t.Errorf("Expected positive uptime, got %f", response.Uptime)
		}
	})
}
```

Don't forget to add the required import at the top:

```go
import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)
```

## Step 3: Run the Test

Open your terminal in the `api-server-go` directory and run:

```bash
# Run just the health handler tests
go test ./internal/handler/health/...

# Or with verbose output to see what's happening
go test -v ./internal/handler/health/...
```

You should see output like:

```
=== RUN   TestHealthCheck
=== RUN   TestHealthCheck/should_return_200_with_status_ok
--- PASS: TestHealthCheck (0.00s)
    --- PASS: TestHealthCheck/should_return_200_with_status_ok (0.00s)
PASS
ok      github.com/ujjwalkirti/mini-vercel-api-server/internal/handler/health  0.123s
```

## Step 4: Implement More Tests

Now that you've got your first test working, let's add the second one:

```go
t.Run("should return valid JSON response", func(t *testing.T) {
	// Arrange
	handler := NewHandler()
	req := httptest.NewRequest("GET", "/health", nil)
	rr := httptest.NewRecorder()

	// Act
	handler.Check(rr, req)

	// Assert
	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}

	var response HealthResponse
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Response is not valid JSON: %v", err)
	}

	// Verify all required fields are present
	if response.Status == "" {
		t.Error("Status field is missing or empty")
	}
	if response.Timestamp.IsZero() {
		t.Error("Timestamp field is missing or zero")
	}
	if response.Uptime == 0 {
		t.Error("Uptime field is missing or zero")
	}
})
```

Run the tests again:

```bash
go test -v ./internal/handler/health/...
```

## Step 5: Implement the Uptime Test

This test requires a delay between calls:

```go
t.Run("should have increasing uptime on subsequent calls", func(t *testing.T) {
	// Arrange
	handler := NewHandler()

	// First call
	req1 := httptest.NewRequest("GET", "/health", nil)
	rr1 := httptest.NewRecorder()
	handler.Check(rr1, req1)

	var response1 HealthResponse
	json.Unmarshal(rr1.Body.Bytes(), &response1)

	// Wait a bit
	time.Sleep(100 * time.Millisecond)

	// Second call
	req2 := httptest.NewRequest("GET", "/health", nil)
	rr2 := httptest.NewRecorder()
	handler.Check(rr2, req2)

	var response2 HealthResponse
	json.Unmarshal(rr2.Body.Bytes(), &response2)

	// Assert
	if response2.Uptime <= response1.Uptime {
		t.Errorf("Expected second uptime (%f) to be greater than first (%f)",
			response2.Uptime, response1.Uptime)
	}
})
```

Add the time import:

```go
import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"  // Add this
)
```

## Step 6: Run All Tests

Now run all the health handler tests:

```bash
go test -v ./internal/handler/health/...
```

All three tests should pass!

## Step 7: Check Coverage

See how much of your code is covered by tests:

```bash
go test -cover ./internal/handler/health/...
```

You should see something like:

```
PASS
coverage: 85.7% of statements
ok      github.com/ujjwalkirti/mini-vercel-api-server/internal/handler/health  0.234s
```

## Step 8: Generate Coverage Report

For a visual coverage report:

```bash
# Generate coverage profile
go test -coverprofile=coverage.out ./internal/handler/health/...

# Open in browser
go tool cover -html=coverage.out
```

This will open an HTML report showing which lines are covered by tests.

## Common Issues and Solutions

### Issue: "undefined: httptest"

**Solution:** Add the import at the top of your test file:
```go
import "net/http/httptest"
```

### Issue: "cannot use handler (variable of type *Handler) as http.Handler"

**Solution:** The health handler's `Check` method takes `http.ResponseWriter` and `*http.Request` directly, not an `http.Handler`. Call it directly:
```go
handler.Check(rr, req)
```

### Issue: Tests fail with "expected status ok, got ''"

**Solution:** Make sure you're parsing the JSON response correctly and the handler is actually setting the status field.

## Next Steps

Once you've mastered the health handler tests, move on to:

1. **Implement Test Helpers** (`internal/testutil/helpers.go`)
   - Start with `CreateTestContext`
   - Then `ParseJSONResponse`
   - Then the assertion helpers

2. **Project Handler Tests** (`internal/handler/project/handler_test.go`)
   - These require mocking repositories
   - Start with the simplest tests (unauthenticated access)
   - Then move to successful operations

3. **Middleware Tests** (`internal/middleware/auth_test.go`)
   - Learn about JWT token testing
   - Mock the JWKS cache

## Testing Best Practices

1. **Use Table-Driven Tests** for similar test cases:
```go
tests := []struct {
	name           string
	expectedStatus int
}{
	{"valid request", 200},
	{"invalid request", 400},
}

for _, tt := range tests {
	t.Run(tt.name, func(t *testing.T) {
		// Test implementation
	})
}
```

2. **Use Subtests** with `t.Run()` to organize related tests

3. **Follow AAA Pattern**:
   - **Arrange**: Set up test data
   - **Act**: Execute the code
   - **Assert**: Verify the results

4. **Use Descriptive Test Names** that explain what's being tested

5. **Test Both Happy and Error Paths**

## Resources

- [Go Testing Package Docs](https://pkg.go.dev/testing)
- [httptest Package Docs](https://pkg.go.dev/net/http/httptest)
- [Table-Driven Tests in Go](https://dave.cheney.net/2019/05/07/prefer-table-driven-tests)
- [Learn Go with Tests](https://quii.gitbook.io/learn-go-with-tests/)

## Getting Help

If you get stuck:
1. Read the error message carefully
2. Check the TODO comments in the test files
3. Look at the [TESTING.md](./TESTING.md) guide
4. Review the [TEST_COVERAGE_PLAN.md](./TEST_COVERAGE_PLAN.md)

Happy testing! 🎉
