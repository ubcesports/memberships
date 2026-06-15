package utils

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func ParseUUID(value string) (pgtype.UUID, error) {
	parsed, err := uuid.Parse(value)
	if err != nil {
		return pgtype.UUID{}, err
	}
	return pgtype.UUID{Bytes: parsed, Valid: true}, nil
}

func UUIDToString(value pgtype.UUID) string {
	if !value.Valid {
		return ""
	}
	return uuid.UUID(value.Bytes).String()
}

func TextPtr(value pgtype.Text) *string {
	if !value.Valid {
		return nil
	}
	return &value.String
}

func TimestamptzPtr(value pgtype.Timestamptz) *time.Time {
	if !value.Valid {
		return nil
	}
	return &value.Time
}

func NumericString(value pgtype.Numeric) string {
	driverValue, err := value.Value()
	if err != nil || driverValue == nil {
		return ""
	}
	return fmt.Sprint(driverValue)
}

func NumericStringPtr(value pgtype.Numeric) *string {
	if !value.Valid {
		return nil
	}
	formatted := NumericString(value)
	return &formatted
}

func NumericIsZero(value pgtype.Numeric) bool {
	return !value.Valid || value.Int == nil || value.Int.Sign() == 0
}
