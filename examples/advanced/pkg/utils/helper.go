package utils

import (
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func StringToPgUUID(s string) pgtype.UUID {
	id, err := uuid.Parse(s)
	if err != nil {
		log.Fatal("Not UUID Format", err.Error())
	}

	return pgtype.UUID{
		Bytes: id,
		Valid: true,
	}
}

func PgUUIDToStringPtr(id pgtype.UUID) *string {
	if !id.Valid {
		return nil
	}

	s := id.String()
	return &s
}

func PgxStringPtr(s *string) pgtype.Text {
	return pgtype.Text{
		String: *s,
		Valid:  true,
	}
}
