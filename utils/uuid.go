package utils

import (
	"log"

	"github.com/google/uuid"
)

func UUIDToString(id uuid.UUID) string {
	return id.String()
}

func StringToUUID(id string) uuid.UUID {
	uuid, err := uuid.Parse(id)
	if err != nil {
		log.Fatal(err)
	}

	return uuid
}
