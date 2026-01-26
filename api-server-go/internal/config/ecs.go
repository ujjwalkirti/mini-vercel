package config

import (
	"os"
	"strconv"
)

type ECSConfig struct {
	Cluster        string
	TaskDefinition string
	Subnets        string
	SecurityGroup  string
	AssignPublicIP string
	ImageName      string
	LaunchType     string
	Count          int
}

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

func GetECSConfig() ECSConfig {
	return ECSConfig{
		Cluster:        getEnvOrDefault("ECS_CLUSTER_NAME", ""),
		TaskDefinition: getEnvOrDefault("ECS_TASK_DEFINITION", ""),
		Subnets:        getEnvOrDefault("ECS_SUBNETS", ""),
		SecurityGroup:  getEnvOrDefault("ECS_SECURITY_GROUPS", ""),
		AssignPublicIP: getEnvOrDefault("ECS_ASSIGN_PUBLIC_IP", "ENABLED"),
		ImageName:      getEnvOrDefault("ECS_IMAGE_NAME", ""),
		LaunchType:     getEnvOrDefault("ECS_LAUNCH_TYPE", "FARGATE"),
		Count:          getEnvAsInt("ECS_COUNT", 1),
	}
}
