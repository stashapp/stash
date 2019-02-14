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
