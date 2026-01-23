package config

import (
	"os"
)

var (
	SupabaseURL       string
	SupabaseJWTSecret string
)

func InitSupabase() error {
	SupabaseURL = os.Getenv("SUPABASE_URL")
	SupabaseJWTSecret = os.Getenv("SUPABASE_JWT_SECRET")

	if SupabaseURL == "" || SupabaseJWTSecret == "" {
		return nil // or return error if you want to enforce these
	}

	return nil
}
