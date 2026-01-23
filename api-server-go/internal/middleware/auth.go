package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ujjwalkirti/mini-vercel-api-server/internal/config"
)

type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

type SupabaseClaims struct {
	Sub   string `json:"sub"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			respondWithError(w, http.StatusUnauthorized, "Unauthorized", "Missing or invalid authorization header")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse and verify JWT token
		token, err := jwt.ParseWithClaims(tokenString, &SupabaseClaims{}, func(token *jwt.Token) (interface{}, error) {
			// Verify signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(config.SupabaseJWTSecret), nil
		})

		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Unauthorized", "Invalid token")
			return
		}

		claims, ok := token.Claims.(*SupabaseClaims)
		if !ok || !token.Valid {
			respondWithError(w, http.StatusUnauthorized, "Unauthorized", "Invalid token claims")
			return
		}

		// Add user to context
		authUser := &AuthUser{
			ID:    claims.Sub,
			Email: claims.Email,
		}

		ctx := context.WithValue(r.Context(), UserContextKey, authUser)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func OptionalAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			token, err := jwt.ParseWithClaims(tokenString, &SupabaseClaims{}, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(config.SupabaseJWTSecret), nil
			})

			if err == nil {
				if claims, ok := token.Claims.(*SupabaseClaims); ok && token.Valid {
					authUser := &AuthUser{
						ID:    claims.Sub,
						Email: claims.Email,
					}
					ctx := context.WithValue(r.Context(), UserContextKey, authUser)
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}

func respondWithError(w http.ResponseWriter, statusCode int, message, errorDetail string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{
		Success: false,
		Message: message,
		Error:   errorDetail,
	})
}
