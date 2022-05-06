package sqlite

import (
	"gopkg.in/guregu/null.v4"
)

// null package does not provide methods to convert null.Int to int pointer
func intFromPtr(i *int) null.Int {
	if i == nil {
		return null.NewInt(0, false)
	}

	return null.IntFrom(int64(*i))
}

func nullIntPtr(i null.Int) *int {
	if !i.Valid {
		return nil
	}

	v := int(i.Int64)
	return &v
}
