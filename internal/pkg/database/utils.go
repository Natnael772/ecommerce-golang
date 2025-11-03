package database

import (
	"github.com/jackc/pgx/v5/pgtype"
)

func AssignIfProvided[T any](param *T, value *T) {
	if value != nil || param == nil {
		*param = *value
	}
}

func ToPGText(s *string) pgtype.Text {
	if s != nil {
		return pgtype.Text{String: *s, Valid: true}
	}
	return pgtype.Text{Valid: false}
}

func ToPGBool(b *bool) pgtype.Bool {
	if b != nil {
		return pgtype.Bool{Bool: *b, Valid: true}
	}
	return pgtype.Bool{Valid: false}
}

func ToPGString(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}
