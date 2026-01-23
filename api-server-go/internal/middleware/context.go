package middleware

import (
	"context"
)

type contextKey string

const UserContextKey contextKey = "user"

type AuthUser struct {
	ID    string
	Email string
}

func GetUserFromContext(ctx context.Context) (*AuthUser, bool) {
	user, ok := ctx.Value(UserContextKey).(*AuthUser)
	return user, ok
}
