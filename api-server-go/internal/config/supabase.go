package config

import (
	"os"
)

var (
	SupabaseURL            string
	SupabaseJWTSecret      string
	SupabaseServiceRoleKey string
)

func InitSupabase() error {
	SupabaseURL = os.Getenv("SUPABASE_URL")
	SupabaseJWTSecret = os.Getenv("SUPABASE_JWT_SECRET")
	SupabaseServiceRoleKey = os.Getenv("SUPABASE_SERVICE_ROLE_KEY")

	if SupabaseURL == "" || SupabaseJWTSecret == "" || SupabaseServiceRoleKey == "" {
		return nil // or return error if you want to enforce these
	}

	return nil
}
