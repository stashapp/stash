package sqlite

import (
	"database/sql/driver"
	"time"
)

const TimestampFormat = time.RFC3339

// Timestamp represents a time stored in RFC3339 format.
type Timestamp struct {
	Timestamp time.Time
}

// Scan implements the Scanner interface.
func (t *Timestamp) Scan(value interface{}) error {
	t.Timestamp = value.(time.Time)
	return nil
}

// Value implements the driver Valuer interface.
func (t Timestamp) Value() (driver.Value, error) {
	return t.Timestamp.Format(TimestampFormat), nil
}

// UTCTimestamp stores a time in UTC.
// TODO - Timestamp should use UTC by default
type UTCTimestamp struct {
	Timestamp
}

// Value implements the driver Valuer interface.
func (t UTCTimestamp) Value() (driver.Value, error) {
	return t.Timestamp.Timestamp.UTC().Format(TimestampFormat), nil
}

// NullTimestamp represents a nullable time stored in RFC3339 format.
type NullTimestamp struct {
	Timestamp time.Time
	Valid     bool
}

// Scan implements the Scanner interface.
func (t *NullTimestamp) Scan(value interface{}) error {
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
func (t NullTimestamp) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}

	return t.Timestamp.Format(TimestampFormat), nil
}

func (t NullTimestamp) TimePtr() *time.Time {
	if !t.Valid {
		return nil
	}

	timestamp := t.Timestamp
	return &timestamp
}

func NullTimestampFromTimePtr(t *time.Time) NullTimestamp {
	if t == nil {
		return NullTimestamp{Valid: false}
	}
	return NullTimestamp{Timestamp: *t, Valid: true}
}
