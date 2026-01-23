package config

import (
	"os"
)

type Config struct {
	DatabaseURL  string
	R2PublicURL  string
	Port         string
}

func Load() *Config {
	return &Config{
		DatabaseURL:  getEnv("DATABASE_URL", ""),
		R2PublicURL:  getEnv("R2_PUBLIC_URL", ""),
		Port:         getEnv("PORT", "8001"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
