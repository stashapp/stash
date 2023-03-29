package sqlite

import (
	"github.com/stashapp/stash/pkg/file"

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

func nullFloatPtr(i null.Float) *float64 {
	if !i.Valid {
		return nil
	}

	v := float64(i.Float64)
	return &v
}

func nullIntFolderIDPtr(i null.Int) *file.FolderID {
	if !i.Valid {
		return nil
	}

	v := file.FolderID(i.Int64)

	return &v
}

func nullIntFileIDPtr(i null.Int) *file.ID {
	if !i.Valid {
		return nil
	}

	v := file.ID(i.Int64)

	return &v
}

func nullIntFromFileIDPtr(i *file.ID) null.Int {
	if i == nil {
		return null.NewInt(0, false)
	}

	return null.IntFrom(int64(*i))
}

func nullIntFromFolderIDPtr(i *file.FolderID) null.Int {
	if i == nil {
		return null.NewInt(0, false)
	}

	return null.IntFrom(int64(*i))
}
