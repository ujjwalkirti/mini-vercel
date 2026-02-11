package config

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

// InitLogger initializes the global zap logger based on environment
func InitLogger() error {
	var err error
	env := os.Getenv("APP_ENV")

	if env == "production" {
		// Production: JSON format, INFO level
		config := zap.NewProductionConfig()
		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		logger, err = config.Build()
	} else {
		// Development: Console format, DEBUG level
		config := zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		logger, err = config.Build()
	}

	if err != nil {
		return err
	}

	// Replace global logger
	zap.ReplaceGlobals(logger)
	return nil
}

// GetLogger returns the global logger instance
func GetLogger() *zap.Logger {
	if logger == nil {
		// Fallback: create a default logger if not initialized
		logger, _ = zap.NewProduction()
		zap.ReplaceGlobals(logger)
	}
	return logger
}

// Sync flushes any buffered log entries
func SyncLogger() error {
	if logger != nil {
		return logger.Sync()
	}
	return nil
}
