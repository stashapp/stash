package zero

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
)

// Bool is a nullable bool. False input is considered null.
// JSON marshals to false if null.
// Considered null to SQL unmarshaled from a false value.
type Bool struct {
	sql.NullBool
}

// NewBool creates a new Bool
func NewBool(b bool, valid bool) Bool {
	return Bool{
		NullBool: sql.NullBool{
			Bool:  b,
			Valid: valid,
		},
	}
}

// BoolFrom creates a new Bool that will be null if false.
func BoolFrom(b bool) Bool {
	return NewBool(b, b)
}

// BoolFromPtr creates a new Bool that be null if b is nil.
func BoolFromPtr(b *bool) Bool {
	if b == nil {
		return NewBool(false, false)
	}
	return NewBool(*b, true)
}

// ValueOrZero returns the inner value if valid, otherwise false.
func (b Bool) ValueOrZero() bool {
	return b.Valid && b.Bool
}

// UnmarshalJSON implements json.Unmarshaler.
// "false" will be considered a null Bool.
func (b *Bool) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, nullBytes) {
		b.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &b.Bool); err != nil {
		return fmt.Errorf("zero: couldn't unmarshal JSON: %w", err)
	}

	b.Valid = b.Bool
	return nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null Bool if the input is false or blank.
// It will return an error if the input is not a float, blank, or "null".
func (b *Bool) UnmarshalText(text []byte) error {
	str := string(text)
	switch str {
	case "", "null":
		b.Valid = false
		return nil
	case "true":
		b.Bool = true
		b.Valid = true
		return nil
	case "false":
		b.Bool = false
		b.Valid = false
		return nil
	}
	return errors.New("invalid input:" + str)
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this Bool is null.
func (b Bool) MarshalJSON() ([]byte, error) {
	if !b.Valid || !b.Bool {
		return []byte("false"), nil
	}
	return []byte("true"), nil
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a zero if this Bool is null.
func (b Bool) MarshalText() ([]byte, error) {
	if !b.Valid || !b.Bool {
		return []byte("false"), nil
	}
	return []byte("true"), nil
}

// SetValid changes this Bool's value and also sets it to be non-null.
func (b *Bool) SetValid(v bool) {
	b.Bool = v
	b.Valid = true
}

// Ptr returns a poBooler to this Bool's value, or a nil poBooler if this Bool is null.
func (b Bool) Ptr() *bool {
	if !b.Valid {
		return nil
	}
	return &b.Bool
}

// IsZero returns true for null or zero Bools, for future omitempty support (Go 1.4?)
func (b Bool) IsZero() bool {
	return !b.Valid || !b.Bool
}

// Equal returns true if both booleans are true and valid, or if both booleans are either false or invalid.
func (b Bool) Equal(other Bool) bool {
	return b.ValueOrZero() == other.ValueOrZero()
}
