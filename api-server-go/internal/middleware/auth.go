package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/ujjwalkirti/mini-vercel-api-server/internal/auth"
	"github.com/ujjwalkirti/mini-vercel-api-server/internal/config"
	"go.uber.org/zap"
)

func AuthMiddleware(jwks *auth.JWKSCache, allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			claims, err := auth.VerifyToken(tokenString, jwks)
			if err != nil {
				logger := config.GetLogger()
				logger.Error("Error in verifying token", zap.Error(err))
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Role / policy validation
			if len(allowedRoles) > 0 {
				allowed := false
				for _, role := range allowedRoles {
					if claims.Role == role {
						allowed = true
						break
					}
				}
				if !allowed {
					http.Error(w, "Forbidden", http.StatusForbidden)
					return
				}
			}

			user := &AuthUser{
				ID:    claims.Subject,
				Email: claims.Email,
				Role:  claims.Role,
			}

			ctx := context.WithValue(r.Context(), UserContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
