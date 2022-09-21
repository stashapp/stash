// Package zero contains SQL types that consider zero input and null input to be equivalent
// with convenient support for JSON and text marshaling.
// Types in this package will JSON marshal to their zero value, even if null.
// Use the null parent package if you don't want this.
package zero

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
)

// nullBytes is a JSON null literal
var nullBytes = []byte("null")

// String is a nullable string.
// JSON marshals to a blank string if null.
// Considered null to SQL if zero.
type String struct {
	sql.NullString
}

// NewString creates a new String
func NewString(s string, valid bool) String {
	return String{
		NullString: sql.NullString{
			String: s,
			Valid:  valid,
		},
	}
}

// StringFrom creates a new String that will be null if s is blank.
func StringFrom(s string) String {
	return NewString(s, s != "")
}

// StringFromPtr creates a new String that be null if s is nil or blank.
// It will make s point to the String's value.
func StringFromPtr(s *string) String {
	if s == nil {
		return NewString("", false)
	}
	return NewString(*s, *s != "")
}

// ValueOrZero returns the inner value if valid, otherwise zero.
func (s String) ValueOrZero() string {
	if !s.Valid {
		return ""
	}
	return s.String
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports string and null input. Blank string input produces a null String.
func (s *String) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, nullBytes) {
		s.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &s.String); err != nil {
		return fmt.Errorf("zero: couldn't unmarshal JSON: %w", err)
	}

	s.Valid = s.String != ""
	return nil
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a blank string when this String is null.
func (s String) MarshalText() ([]byte, error) {
	if !s.Valid {
		return []byte{}, nil
	}
	return []byte(s.String), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null String if the input is a blank string.
func (s *String) UnmarshalText(text []byte) error {
	s.String = string(text)
	s.Valid = s.String != ""
	return nil
}

// SetValid changes this String's value and also sets it to be non-null.
func (s *String) SetValid(v string) {
	s.String = v
	s.Valid = true
}

// Ptr returns a pointer to this String's value, or a nil pointer if this String is null.
func (s String) Ptr() *string {
	if !s.Valid {
		return nil
	}
	return &s.String
}

// IsZero returns true for null or empty strings, for potential future omitempty support.
func (s String) IsZero() bool {
	return !s.Valid || s.String == ""
}

// Equal returns true if both strings have the same value or are both either null or empty.
func (s String) Equal(other String) bool {
	return s.ValueOrZero() == other.ValueOrZero()
}
