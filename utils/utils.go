package utils

import "github.com/google/uuid"

func BuildID(existingID *string) string {
	if existingID == nil || *existingID == "" {
		return uuid.New().String()
	}
	return *existingID
}
