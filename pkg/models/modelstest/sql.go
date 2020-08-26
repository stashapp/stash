package modelstest

import "database/sql"

func NullString(v string) sql.NullString {
	return sql.NullString{
		String: v,
		Valid:  true,
	}
}
