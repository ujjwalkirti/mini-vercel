package main

import (
	"github.com/joho/godotenv"
	"github.com/ujjwalkirti/mini-vercel-api-server/internal/app"
	"github.com/ujjwalkirti/mini-vercel-api-server/internal/config"
	"go.uber.org/zap"
)

func main() {
	// Load environment variables first
	envErr := godotenv.Load() // Ignore error, will use system env if .env not found

	// Initialize logger after env vars are loaded
	if err := config.InitLogger(); err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}

	if envErr != nil {
		config.GetLogger().Warn("Failed to load .env file", zap.Error(envErr))
	}

	defer config.SyncLogger()

	application := app.New()
	application.Run()
}
