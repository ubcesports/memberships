package util

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func GetValidatedUUID(uuid string) (pgtype.UUID, error) {
	var pgUUID pgtype.UUID
	if err := pgUUID.Scan(uuid); err != nil {
		return pgtype.UUID{}, err
	}
	return pgUUID, nil
}

func TextPointer(value pgtype.Text) *string {
	if !value.Valid {
		return nil
	}
	return &value.String
}

func TimestampPointer(value pgtype.Timestamptz) *time.Time {
	if !value.Valid {
		return nil
	}
	return &value.Time
}
