package config

import (
	"os"
)

var (
	AWSAccessKey = os.Getenv("AWS_ACCESS_KEY_ID")
	AWSSecretKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
	AWSRegion    = os.Getenv("AWS_REGION")
)

func InitAWS() error {
	if AWSAccessKey == "" || AWSSecretKey == "" || AWSRegion == "" {
		return nil // or return error if you want to enforce these
	}

	return nil
}
