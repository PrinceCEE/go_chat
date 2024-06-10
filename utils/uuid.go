package utils

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func UUIDToString(id uuid.UUID) string {
	return id.String()
}

func StringToUUID(id string) uuid.UUID {
	uuid, _ := uuid.Parse(id)
	return uuid
}

func StringToText(str string) pgtype.Text {
	return pgtype.Text{String: str, Valid: true}
}

func TextToString(t pgtype.Text) string {
	return t.String
}
