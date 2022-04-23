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

func stringDblPtrToPtr(v **string) *string {
	if v == nil {
		return nil
	}

	return *v
}

func intDblPtrToPtr(v **int) *int {
	if v == nil {
		return nil
	}

	return *v
}

func intDblPtrToStringPtr(v **int) *string {
	if v == nil {
		return nil
	}

	var ret *string

	if *v != nil {
		vv := strconv.Itoa(**v)
		ret = &vv
	}

	return ret
}
