package config

import (
	"os"
)

type AWSConfig struct {
	Region    string
	AccessKey string
	SecretKey string
}

func GetAWSConfig() AWSConfig {
	return AWSConfig{
		Region:    os.Getenv("AWS_REGION"),
		AccessKey: os.Getenv("AWS_ACCESS_KEY"),
		SecretKey: os.Getenv("AWS_SECRET_KEY"),
	}
}
