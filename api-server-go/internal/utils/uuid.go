package utils

import "github.com/google/uuid"

func IsValidUUID(id string) bool {
	_, err := uuid.Parse(id)
	return err == nil
}

func GenerateUUID() string {
	return uuid.New().String()
}
