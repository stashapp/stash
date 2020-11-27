package models

import (
	"database/sql/driver"
	"time"
)

type SQLiteTimestamp struct {
	Timestamp time.Time
}

// Scan implements the Scanner interface.
func (t *SQLiteTimestamp) Scan(value interface{}) error {
	t.Timestamp = value.(time.Time)
	return nil
}

// Value implements the driver Valuer interface.
func (t SQLiteTimestamp) Value() (driver.Value, error) {
	return t.Timestamp.Format(time.RFC3339), nil
}

type NullSQLiteTimestamp struct {
	Timestamp time.Time
	Valid     bool
}

// Scan implements the Scanner interface.
func (t *NullSQLiteTimestamp) Scan(value interface{}) error {
	var ok bool
	t.Timestamp, ok = value.(time.Time)
	if !ok {
		t.Timestamp = time.Time{}
		t.Valid = false
		return nil
	}

	t.Valid = true
	return nil
}

// Value implements the driver Valuer interface.
func (t NullSQLiteTimestamp) Value() (driver.Value, error) {
	if t.Timestamp.IsZero() {
		return nil, nil
	}

	return t.Timestamp.Format(time.RFC3339), nil
}
