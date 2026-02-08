package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ujjwalkirti/mini-vercel-api-server/internal/auth"
)

// TestAuthMiddleware tests the authentication middleware
func TestAuthMiddleware(t *testing.T) {
	t.Run("should return 401 when Authorization header is missing", func(t *testing.T) {
		// Arrange - Create mock JWKS cache (nil is fine for this test as auth check happens before JWKS usage)
		var jwks *auth.JWKSCache = nil

		// Create middleware with JWKS
		middleware := AuthMiddleware(jwks)

		// Create test handler that should not be called
		handlerCalled := false
		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handlerCalled = true
			w.WriteHeader(http.StatusOK)
		})

		// Wrap test handler with middleware
		wrappedHandler := middleware(testHandler)

		// Create test request without Authorization header
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rr := httptest.NewRecorder()

		// Act - Call middleware
		wrappedHandler.ServeHTTP(rr, req)

		// Assert - Assert status code is 401
		if rr.Code != http.StatusUnauthorized {
			t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, rr.Code)
		}

		// Assert error message is "Unauthorized"
		if rr.Body.String() != "Unauthorized\n" {
			t.Errorf("Expected error message 'Unauthorized', got '%s'", rr.Body.String())
		}

		// Assert test handler was not called
		if handlerCalled {
			t.Error("Expected handler to not be called when authorization is missing")
		}
	})

	t.Run("should return 401 when Authorization header does not start with Bearer", func(t *testing.T) {
		// Arrange
		var jwks *auth.JWKSCache = nil
		middleware := AuthMiddleware(jwks)

		handlerCalled := false
		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handlerCalled = true
			w.WriteHeader(http.StatusOK)
		})
		wrappedHandler := middleware(testHandler)

		// Create test request with invalid Authorization header
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Basic sometoken")
		rr := httptest.NewRecorder()

		// Act
		wrappedHandler.ServeHTTP(rr, req)

		// Assert
		if rr.Code != http.StatusUnauthorized {
			t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, rr.Code)
		}
		if handlerCalled {
			t.Error("Expected handler to not be called when authorization format is invalid")
		}
	})

	t.Run("should return 401 when token is empty after Bearer", func(t *testing.T) {
		// Arrange
		var jwks *auth.JWKSCache = nil
		middleware := AuthMiddleware(jwks)

		handlerCalled := false
		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handlerCalled = true
			w.WriteHeader(http.StatusOK)
		})
		wrappedHandler := middleware(testHandler)

		// Create test request with empty token
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Bearer ")
		rr := httptest.NewRecorder()

		// Act
		wrappedHandler.ServeHTTP(rr, req)

		// Assert
		if rr.Code != http.StatusUnauthorized {
			t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, rr.Code)
		}
		if handlerCalled {
			t.Error("Expected handler to not be called when token is empty")
		}
	})

	t.Run("should return 401 when token verification fails", func(t *testing.T) {
		t.Skip("Requires JWT mocking - auth.VerifyToken needs refactoring for dependency injection")
		// NOTE: This test would require refactoring auth.VerifyToken to accept a verifier interface
		// Current implementation directly calls auth.VerifyToken which is hard to mock
	})

	t.Run("should return 401 when token is expired", func(t *testing.T) {
		t.Skip("Requires JWT mocking - auth.VerifyToken needs refactoring for dependency injection")
	})

	t.Run("should return 401 when token signature is invalid", func(t *testing.T) {
		t.Skip("Requires JWT mocking - auth.VerifyToken needs refactoring for dependency injection")
	})

	t.Run("should call next handler when token is valid", func(t *testing.T) {
		t.Skip("Requires JWT mocking - auth.VerifyToken needs refactoring for dependency injection")
	})

	t.Run("should add user to request context when token is valid", func(t *testing.T) {
		t.Skip("Requires JWT mocking - auth.VerifyToken needs refactoring for dependency injection")
	})

	t.Run("should extract correct user fields from token claims", func(t *testing.T) {
		t.Skip("Requires JWT mocking - auth.VerifyToken needs refactoring for dependency injection")
	})

	t.Run("should allow request when no roles are specified", func(t *testing.T) {
		t.Skip("Requires JWT mocking - auth.VerifyToken needs refactoring for dependency injection")
	})

	t.Run("should return 403 when user role does not match allowed roles", func(t *testing.T) {
		t.Skip("Requires JWT mocking - auth.VerifyToken needs refactoring for dependency injection")
	})

	t.Run("should allow request when user role matches one of allowed roles", func(t *testing.T) {
		t.Skip("Requires JWT mocking - auth.VerifyToken needs refactoring for dependency injection")
	})

	t.Run("should allow request when user role matches exactly", func(t *testing.T) {
		t.Skip("Requires JWT mocking - auth.VerifyToken needs refactoring for dependency injection")
	})

	t.Run("should handle multiple allowed roles correctly", func(t *testing.T) {
		t.Skip("Requires JWT mocking - auth.VerifyToken needs refactoring for dependency injection")
	})

	t.Run("should trim Bearer prefix correctly with extra spaces", func(t *testing.T) {
		t.Skip("Requires JWT mocking - auth.VerifyToken needs refactoring for dependency injection")
	})

	t.Run("should handle case-sensitive Bearer prefix", func(t *testing.T) {
		// Arrange
		var jwks *auth.JWKSCache = nil
		middleware := AuthMiddleware(jwks)

		handlerCalled := false
		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handlerCalled = true
			w.WriteHeader(http.StatusOK)
		})
		wrappedHandler := middleware(testHandler)

		// Create test request with lowercase "bearer"
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "bearer sometoken")
		rr := httptest.NewRecorder()

		// Act
		wrappedHandler.ServeHTTP(rr, req)

		// Assert - case matters for "Bearer"
		if rr.Code != http.StatusUnauthorized {
			t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, rr.Code)
		}
		if handlerCalled {
			t.Error("Expected handler to not be called when Bearer prefix is lowercase")
		}
	})
}

// TestGetUserFromContext tests retrieving user from context
func TestGetUserFromContext(t *testing.T) {
	t.Run("should return user when present in context", func(t *testing.T) {
		// Arrange
		expectedUser := &AuthUser{
			ID:    "user-123",
			Email: "test@example.com",
			Role:  "admin",
		}
		ctx := context.WithValue(context.Background(), UserContextKey, expectedUser)

		// Act
		user, ok := GetUserFromContext(ctx)

		// Assert
		if !ok {
			t.Error("Expected ok to be true")
		}
		if user == nil {
			t.Fatal("Expected user to not be nil")
		}
		if user.ID != expectedUser.ID {
			t.Errorf("Expected user ID %s, got %s", expectedUser.ID, user.ID)
		}
		if user.Email != expectedUser.Email {
			t.Errorf("Expected user email %s, got %s", expectedUser.Email, user.Email)
		}
		if user.Role != expectedUser.Role {
			t.Errorf("Expected user role %s, got %s", expectedUser.Role, user.Role)
		}
	})

	t.Run("should return false when user is not in context", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		// Act
		user, ok := GetUserFromContext(ctx)

		// Assert
		if ok {
			t.Error("Expected ok to be false")
		}
		if user != nil {
			t.Error("Expected user to be nil")
		}
	})

	t.Run("should return false when context value is wrong type", func(t *testing.T) {
		// Arrange
		ctx := context.WithValue(context.Background(), UserContextKey, "not-a-user")

		// Act
		user, ok := GetUserFromContext(ctx)

		// Assert
		if ok {
			t.Error("Expected ok to be false when value is wrong type")
		}
		if user != nil {
			t.Error("Expected user to be nil when value is wrong type")
		}
	})

	t.Run("should handle nil context gracefully", func(t *testing.T) {
		// Note: Passing nil context will panic - this is expected Go behavior
		// Context should never be nil according to Go conventions
		// We document this rather than test it
		t.Skip("Nil context is not supported per Go conventions - would panic")
	})
}

// Note: Additional helper functions for JWT token creation and verification mocking
// would require refactoring the auth package to use dependency injection.
// For integration tests, consider using real JWT tokens with test keys.
