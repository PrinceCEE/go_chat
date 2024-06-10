package utils

import (
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
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

func StringToText(str string) pgtype.Text {
	return pgtype.Text{String: str, Valid: true}
}

func TextToString(t pgtype.Text) string {
	return t.String
}
