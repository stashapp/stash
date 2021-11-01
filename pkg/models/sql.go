package models

import (
	"database/sql"
	"strconv"
)

func NullString(v string) sql.NullString {
	return sql.NullString{
		String: v,
		Valid:  true,
	}
}

func NullStringPtr(v string) *sql.NullString {
	return &sql.NullString{
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

func nullStringPtrToStringPtr(v *sql.NullString) *string {
	if v == nil || !v.Valid {
		return nil
	}

	vv := v.String
	return &vv
}

func nullInt64PtrToIntPtr(v *sql.NullInt64) *int {
	if v == nil || !v.Valid {
		return nil
	}

	vv := int(v.Int64)
	return &vv
}

func nullInt64PtrToStringPtr(v *sql.NullInt64) *string {
	if v == nil || !v.Valid {
		return nil
	}

	vv := strconv.FormatInt(v.Int64, 10)
	return &vv
}
