package zero

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// Time is a nullable time.Time.
// JSON marshals to the zero value for time.Time if null.
// Considered to be null to SQL if zero.
type Time struct {
	sql.NullTime
}

// Value implements the driver Valuer interface.
func (t Time) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}
	return t.Time, nil
}

// NewTime creates a new Time.
func NewTime(t time.Time, valid bool) Time {
	return Time{
		NullTime: sql.NullTime{
			Time:  t,
			Valid: valid,
		},
	}
}

// TimeFrom creates a new Time that will
// be null if t is the zero value.
func TimeFrom(t time.Time) Time {
	return NewTime(t, !t.IsZero())
}

// TimeFromPtr creates a new Time that will
// be null if t is nil or *t is the zero value.
func TimeFromPtr(t *time.Time) Time {
	if t == nil {
		return NewTime(time.Time{}, false)
	}
	return TimeFrom(*t)
}

// ValueOrZero returns the inner value if valid, otherwise zero.
func (t Time) ValueOrZero() time.Time {
	if !t.Valid {
		return time.Time{}
	}
	return t.Time
}

// MarshalJSON implements json.Marshaler.
// It will encode the zero value of time.Time
// if this time is invalid.
func (t Time) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return (time.Time{}).MarshalJSON()
	}
	return t.Time.MarshalJSON()
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports string and null input.
func (t *Time) UnmarshalJSON(data []byte) error {
	switch string(data) {
	case "null", `""`:
		t.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &t.Time); err != nil {
		return fmt.Errorf("zero: couldn't unmarshal JSON: %w", err)
	}

	t.Valid = !t.Time.IsZero()
	return nil
}

// MarshalText implements encoding.TextMarshaler.
// It will encode to an empty time.Time if invalid.
func (t Time) MarshalText() ([]byte, error) {
	ti := t.Time
	if !t.Valid {
		ti = time.Time{}
	}
	return ti.MarshalText()
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It has compatibility with the null package in that it will accept empty strings as invalid values,
// which will be unmarshaled to an invalid zero value.
func (t *Time) UnmarshalText(text []byte) error {
	str := string(text)
	// allowing "null" is for backwards compatibility with v3
	if str == "" || str == "null" {
		t.Valid = false
		return nil
	}
	if err := t.Time.UnmarshalText(text); err != nil {
		return fmt.Errorf("zero: couldn't unmarshal text: %w", err)
	}
	t.Valid = !t.Time.IsZero()
	return nil
}

// SetValid changes this Time's value and
// sets it to be non-null.
func (t *Time) SetValid(v time.Time) {
	t.Time = v
	t.Valid = true
}

// Ptr returns a pointer to this Time's value,
// or a nil pointer if this Time is zero.
func (t Time) Ptr() *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}

// IsZero returns true for null or zero Times, for potential future omitempty support.
func (t Time) IsZero() bool {
	return !t.Valid || t.Time.IsZero()
}

// Equal returns true if both Time objects encode the same time or are both are either null or zero.
// Two times can be equal even if they are in different locations.
// For example, 6:00 +0200 CEST and 4:00 UTC are Equal.
func (t Time) Equal(other Time) bool {
	return t.ValueOrZero().Equal(other.ValueOrZero())
}

// ExactEqual returns true if both Time objects are equal or both are either null or zero.
// ExactEqual returns false for times that are in different locations or
// have a different monotonic clock reading.
func (t Time) ExactEqual(other Time) bool {
	return t.ValueOrZero() == other.ValueOrZero()
}
