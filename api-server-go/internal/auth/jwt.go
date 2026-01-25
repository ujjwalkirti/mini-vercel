package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ujjwalkirti/mini-vercel-api-server/internal/config"
)

type SupabaseClaims struct {
	Email string `json:"email"`
	Role  string `json:"role"`
	Aud   string `json:"aud"`
	jwt.RegisteredClaims
}

func VerifyToken(tokenString string, cache *JWKSCache) (*SupabaseClaims, error) {
	claims := &SupabaseClaims{}

	config.InitSupabase()

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(t *jwt.Token) (interface{}, error) {
			// Algorithm check
			switch t.Method.Alg() {
			case jwt.SigningMethodRS256.Alg():
			case jwt.SigningMethodES256.Alg():
			default:
				return nil, fmt.Errorf("unexpected signing method: %s", t.Method.Alg())
			}

			// kid extraction
			kid, ok := t.Header["kid"].(string)
			if !ok {
				return nil, fmt.Errorf("missing kid header")
			}

			// fetch key from JWKS cache
			return cache.Get(kid)
		},
		jwt.WithIssuer(config.SupabaseURL+"/auth/v1"),
		jwt.WithLeeway(2*time.Minute),
	)

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
