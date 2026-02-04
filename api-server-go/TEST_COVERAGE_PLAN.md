# Test Coverage Plan

This document provides an overview of all test skeleton files created for the Mini Vercel API Server.

## Overview

All test files have been created with skeleton functions containing TODO comments. Each test function outlines what needs to be tested but does not contain the implementation. This allows you to learn Go testing gradually by implementing one test at a time.

## Test Files Created

### 1. Handler Tests

#### Health Handler
**File:** `internal/handler/health/handler_test.go`
- ✅ TestHealthCheck
  - should return 200 with status ok
  - should return valid JSON response
  - should have increasing uptime on subsequent calls
- ✅ BenchmarkHealthCheck

**Test Count:** 3 tests + 1 benchmark

#### Project Handler
**File:** `internal/handler/project/handler_test.go`
- ✅ TestGetProjects (5 test cases)
  - Unauthenticated access
  - Empty projects list
  - All projects for user
  - Database errors
  - Including recent deployments
- ✅ TestGetProject (6 test cases)
  - Unauthenticated access
  - Invalid UUID
  - Project not found
  - Different user's project
  - Valid project retrieval
  - Database errors
- ✅ TestCreateProject (9 test cases)
  - Unauthenticated access
  - Invalid JSON body
  - Missing name
  - Missing github_url
  - Empty fields
  - Successful creation
  - Unique subdomain generation
  - Database errors
  - Whitespace trimming
- ✅ TestUpdateProject (8 test cases)
  - Unauthenticated access
  - Invalid UUID
  - Invalid JSON
  - Project not found
  - Different user's project
  - Update name
  - Update github_url
  - Subdomain protection
  - Database errors
- ✅ TestDeleteProject (8 test cases)
  - Unauthenticated access
  - Invalid UUID
  - Project not found
  - Different user's project
  - Successful deletion
  - Cascade delete deployments
  - Deployment deletion failure
  - Project deletion failure

**Test Count:** 36 tests + helper functions

#### Deployment Handler
**File:** `internal/handler/deployment/handler_test.go`
- ✅ TestGetDeploymentsByProject (8 test cases)
  - Unauthenticated access
  - Invalid project UUID
  - Empty deployments list
  - All deployments for project
  - Project not found
  - Unauthorized access
  - Database errors
  - Status field presence
- ✅ TestGetDeployment (6 test cases)
  - Unauthenticated access
  - Invalid deployment UUID
  - Deployment not found
  - Deployment with project data
  - Database errors
  - User ownership verification
- ✅ TestCreateDeployment (11 test cases)
  - Unauthenticated access
  - Invalid JSON
  - Missing project_id
  - Empty project_id
  - Project not found
  - Unauthorized access
  - Successful creation with QUEUED status
  - ECS task with environment variables
  - ECS task failure
  - Response with deployment URL
  - Response format
  - Nil ECS service handling
- ✅ TestGetDeploymentLogs (9 test cases)
  - Unauthenticated access
  - Invalid deployment UUID
  - Deployment not found
  - Empty logs array
  - Logs from ClickHouse
  - Logs service failure
  - Nil logs service handling
  - User ownership verification
  - Chronological log ordering

**Test Count:** 34 tests + helper functions

### 2. Middleware Tests

#### Authentication Middleware
**File:** `internal/middleware/auth_test.go`
- ✅ TestAuthMiddleware (15 test cases)
  - Missing Authorization header
  - Invalid Authorization format
  - Empty Bearer token
  - Token verification failure
  - Expired token
  - Invalid signature
  - Valid token processing
  - User context injection
  - User field extraction
  - No role restrictions
  - Role mismatch (403)
  - Single allowed role
  - Exact role match
  - Multiple allowed roles
  - Bearer prefix handling
  - Case sensitivity
- ✅ TestGetUserFromContext (4 test cases)
  - User present in context
  - User not in context
  - Wrong type in context
  - Nil context handling

**Test Count:** 19 tests + helper functions

### 3. Router Tests

#### Router Integration
**File:** `internal/router/router_test.go`
- ✅ TestRouterSetup (5 test cases)
  - Router creation
  - CORS middleware
  - Logger middleware
  - Recoverer middleware
  - Request ID middleware
- ✅ TestHealthRoutes (2 test cases)
  - Health route mounting
  - No authentication required
- ✅ TestProjectRoutes (7 test cases)
  - Project routes mounting
  - Authentication required
  - GET /projects
  - GET /projects/:id
  - POST /projects
  - PUT /projects/:id
  - DELETE /projects/:id
- ✅ TestDeploymentRoutes (6 test cases)
  - Deployment routes mounting
  - Authentication required
  - GET /projects/:projectId (deployments)
  - GET /deployments/:id
  - POST /deploy
  - GET /deployments/:id/logs
- ✅ TestRouteNotFound (2 test cases)
  - Non-existent routes
  - Invalid HTTP methods
- ✅ TestCORSConfiguration (5 test cases)
  - Allowed origins
  - Allowed methods
  - Allowed headers
  - Max age
  - Credentials policy
- ✅ TestMiddlewareOrder (3 test cases)
  - Logger middleware order
  - Recoverer middleware order
  - Auth middleware order

**Test Count:** 30 tests + helper functions

### 4. Service Tests

#### ECS Service
**File:** `internal/service/ecs/service_test.go`
- ✅ TestRunTask (14 test cases)
  - Task with correct configuration
  - Environment variables
  - Cluster name
  - Task definition
  - Network configuration with subnets
  - Security group
  - Public IP assignment
  - Launch type
  - Container image override
  - Task count
  - AWS API failures
  - Context cancellation
  - Empty environment variables
  - Environment variable sanitization
- ✅ TestEnvVar (2 test cases)
  - Struct fields
  - AWS SDK conversion
- ✅ TestNew (3 test cases)
  - Service creation with all parameters
  - Empty subnets handling
  - Missing optional parameters

**Test Count:** 19 tests + helper functions

#### Logs Service
**File:** `internal/service/logs/service_test.go`
- ✅ TestGetDeploymentLogs (7 test cases)
  - Valid deployment ID
  - Empty logs array
  - Repository failures
  - Invalid deployment ID
  - Chronological ordering
  - Context cancellation
  - All log fields present
- ✅ TestLogEventStructure (3 test cases)
  - Required fields
  - JSON marshaling
  - JSON unmarshaling

**Test Count:** 10 tests + helper functions

### 5. Repository Tests

#### Project Repository
**File:** `internal/repository/project/repository_test.go`
- ✅ TestListByUser (5 test cases)
  - All projects for user
  - Empty projects list
  - Database query failures
  - Recent deployment inclusion
  - Ordering by created_at
- ✅ TestGetByIDAndUserID (4 test cases)
  - Valid project retrieval
  - Project not found
  - Unauthorized access
  - Database query failures
- ✅ TestCreate (5 test cases)
  - Create with all fields
  - UUID generation
  - Timestamp setting
  - Database insert failures
  - Duplicate subdomain handling
- ✅ TestUpdate (5 test cases)
  - Field updates
  - updated_at timestamp update
  - created_at preservation
  - Project not found
  - Database update failures
- ✅ TestDelete (4 test cases)
  - Successful deletion
  - Project not found
  - Database delete failures
  - Cascading deletes

**Test Count:** 23 tests + helper functions

#### Deployment Repository
**File:** `internal/repository/deployment/repository_test.go`
- ✅ TestGetByProjectID (5 test cases)
  - All deployments for project
  - Empty deployments list
  - User ownership verification
  - Ordering by created_at
  - Database query failures
- ✅ TestGetByIDWithProject (5 test cases)
  - Deployment with project data
  - User ownership verification
  - Deployment not found
  - Database query failures
  - JOIN query verification
- ✅ TestCreate (7 test cases)
  - Create with all fields
  - UUID generation
  - Default QUEUED status
  - Timestamp setting
  - Foreign key validation
  - Database insert failures
  - Return created deployment
- ✅ TestDeleteByProjectID (4 test cases)
  - Delete all for project
  - Other projects unaffected
  - Empty project handling
  - Database delete failures
- ✅ TestUpdateStatus (4 test cases)
  - Status update
  - updated_at timestamp update
  - Valid status transitions
  - Deployment not found

**Test Count:** 25 tests + helper functions

### 6. Test Utilities

#### Test Helpers
**File:** `internal/testutil/helpers.go`

Functions to implement:
- CreateTestContext - Create authenticated context
- CreateAuthenticatedRequest - Create HTTP request with user context
- CreateRequestWithBearerToken - Create request with JWT token
- ParseJSONResponse - Parse JSON response from httptest recorder
- ParseJSONResponseArray - Parse JSON array response
- CreateJSONRequestBody - Marshal data to JSON for requests
- AssertStatusCode - Assert HTTP status code
- AssertResponseContains - Assert response contains substring
- AssertJSONField - Assert JSON field value
- AssertJSONFieldExists - Assert JSON field exists
- AssertHeaderExists - Assert response header exists
- AssertHeaderValue - Assert response header value
- CreateTestUser - Create test AuthUser
- AssertErrorResponse - Assert error response format
- AssertSuccessResponse - Assert success response
- CompareJSON - Compare JSON strings
- GenerateValidUUID - Generate valid UUID
- GenerateInvalidUUID - Return invalid UUID

#### Mock Implementations
**File:** `internal/testutil/mocks.go`

Mock types to implement:
- MockProjectRepository - Mock project repository with configurable behavior
- MockDeploymentRepository - Mock deployment repository
- MockECSService - Mock AWS ECS service
- MockLogsService - Mock logs service
- MockJWKSCache - Mock JWKS cache for JWT verification

Helper functions:
- CreateMockProject - Create test project data
- CreateMockDeployment - Create test deployment data
- CreateMockLogEvent - Create test log event
- CreateMockProjects - Create multiple test projects
- CreateMockDeployments - Create multiple test deployments

## Total Test Coverage

| Category | Files | Test Functions | Test Cases |
|----------|-------|----------------|------------|
| Handlers | 3 | 18 | 73 |
| Middleware | 1 | 2 | 19 |
| Router | 1 | 7 | 30 |
| Services | 2 | 5 | 29 |
| Repositories | 2 | 10 | 48 |
| **Total** | **9** | **42** | **199** |

Plus 1 benchmark test and 36+ helper functions across 2 utility files.

## Implementation Priority

Recommended order for implementing tests:

### Phase 1: Basics (Start Here)
1. `internal/testutil/helpers.go` - Implement core helper functions first
2. `internal/handler/health/handler_test.go` - Simplest tests, no dependencies
3. `internal/testutil/mocks.go` - Implement mock types

### Phase 2: Foundation
4. `internal/middleware/auth_test.go` - Learn JWT/auth testing
5. `internal/handler/project/handler_test.go` - Learn handler testing with mocks

### Phase 3: Complex Flows
6. `internal/handler/deployment/handler_test.go` - Multiple dependencies
7. `internal/service/ecs/service_test.go` - External service mocking
8. `internal/service/logs/service_test.go` - ClickHouse mocking

### Phase 4: Data Layer
9. `internal/repository/project/repository_test.go` - Database testing
10. `internal/repository/deployment/repository_test.go` - Database testing

### Phase 5: Integration
11. `internal/router/router_test.go` - Full integration tests

## Running Tests

```bash
# Run all tests
go test ./...

# Run specific package tests
go test ./internal/handler/health/...

# Run with coverage
go test -cover ./...

# Run specific test function
go test -run TestHealthCheck ./internal/handler/health/...

# Run with verbose output
go test -v ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Learning Resources

See [TESTING.md](./TESTING.md) for:
- Detailed testing guidelines
- Mocking strategies
- Common patterns
- Best practices
- Example implementations

## Notes

- All test files contain skeleton functions with TODO comments
- Each TODO describes what needs to be tested
- Implement tests gradually as you learn Go testing concepts
- Start with simple tests (health check) and progress to complex ones (deployment with ECS)
- Use the helper functions in `testutil` to reduce boilerplate
- Follow the Arrange-Act-Assert pattern in all tests

## Next Steps

1. Read [TESTING.md](./TESTING.md) for testing guidelines
2. Start with Phase 1: Implement test helpers
3. Work through health handler tests
4. Gradually implement other tests as you learn
5. Aim for 80%+ code coverage
6. Run tests frequently during development

Good luck with your testing journey! 🚀
