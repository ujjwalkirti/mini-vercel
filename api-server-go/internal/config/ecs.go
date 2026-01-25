package config

import (
	"os"
	"strconv"
)

var (
	ECS_CLUSTER           = os.Getenv("ECS_CLUSTER")
	ECS_TASK_DEFINITION   = os.Getenv("ECS_TASK_DEFINITION")
	ECS_SUBNETS           = os.Getenv("ECS_SUBNETS")
	ECS_SECURITY_GROUP    = os.Getenv("ECS_SECURITY_GROUP")
	ECS_ASSIGN_PUBLIC_IP  = getEnvOrDefault("ECS_ASSIGN_PUBLIC_IP", "ENABLED")
	ECS_IMAGE_NAME        = os.Getenv("ECS_IMAGE_NAME")
	ECS_LAUNCH_TYPE       = getEnvOrDefault("ECS_LAUNCH_TYPE", "FARGATE")
	ECS_COUNT             = getEnvAsInt("ECS_COUNT", 1)
)

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}
