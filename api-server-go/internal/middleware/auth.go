package middleware

import (
	"context"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log"
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

// parseECPublicKey parses PEM-encoded ECDSA public key
func parseECPublicKey(pemStr string) (*ecdsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	ecdsaKey, ok := pub.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not an ECDSA public key")
	}

	return ecdsaKey, nil
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			respondWithError(w, http.StatusUnauthorized, "Unauthorized", "Missing or invalid authorization header")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		config.InitSupabase()

		// Parse and verify JWT token
		token, err := jwt.ParseWithClaims(tokenString, &SupabaseClaims{}, func(token *jwt.Token) (any, error) {
			// Verify signing method - Supabase uses ES256 (ECDSA)
			if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			// Parse the ECDSA public key from JWT secret
			publicKey, err := parseECPublicKey(config.SupabaseJWTSecret)
			if err != nil {
				return nil, fmt.Errorf("failed to parse public key: %w", err)
			}

			return publicKey, nil
		})

		if err != nil {
			log.Println(err.Error())
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

			token, err := jwt.ParseWithClaims(tokenString, &SupabaseClaims{}, func(token *jwt.Token) (any, error) {
				// Verify signing method - Supabase uses ES256 (ECDSA)
				if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}

				// Parse the ECDSA public key from JWT secret
				publicKey, err := parseECPublicKey(config.SupabaseJWTSecret)
				if err != nil {
					return nil, fmt.Errorf("failed to parse public key: %w", err)
				}

				return publicKey, nil
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
