package config

import (
	"os"
	"strings"
)

type KafkaConfig struct {
	Brokers  []string
	ClientID string
	Username string
	Password string
}

func LoadKafkaConfig() KafkaConfig {
	return KafkaConfig{
		Brokers:  strings.Split(os.Getenv("KAFKA_BROKERS"), ","),
		ClientID: os.Getenv("KAFKA_CLIENT_ID"),
		Username: os.Getenv("KAFKA_USERNAME"),
		Password: os.Getenv("KAFKA_PASSWORD"),
	}
}
