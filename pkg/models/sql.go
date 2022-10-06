package models

import (
	"database/sql"
)

func NullString(v string) sql.NullString {
	return sql.NullString{
		String: v,
		Valid:  true,
	}
}

func NullInt64(v int64) sql.NullInt64 {
	return sql.NullInt64{
		Int64: v,
		Valid: true,
	}
}
