package models

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"

	"github.com/stashapp/stash/pkg/utils"
)

type SQLiteDate struct {
	String string
	Valid  bool
}

// Scan implements the Scanner interface.
func (t *SQLiteDate) Scan(value interface{}) error {
	dateTime, ok := value.(time.Time)
	if !ok {
		t.String = ""
		t.Valid = false
		return nil
	}

	t.String = dateTime.Format("2006-01-02")
	if t.String != "" && t.String != "0001-01-01" {
		t.Valid = true
	} else {
		t.Valid = false
	}
	return nil
}

// Value implements the driver Valuer interface.
func (t SQLiteDate) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}

	s := strings.TrimSpace(t.String)
	// handle empty string
	if s == "" {
		return "", nil
	}

	result, err := utils.ParseDateStringAsFormat(s, "2006-01-02")
	if err != nil {
		return nil, fmt.Errorf("converting sqlite date %q: %w", s, err)
	}
	return result, nil
}

func (t *SQLiteDate) StringPtr() *string {
	if t == nil || !t.Valid {
		return nil
	}

	vv := t.String
	return &vv
}
