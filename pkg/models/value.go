package models

import (
	"strconv"
	"time"
)

// OptionalString represents an optional string argument that may be null.
// A value is only considered null if both Set and Null is true.
type OptionalString struct {
	Value string
	Null  bool
	Set   bool
}

// Ptr returns a pointer to the underlying value. Returns nil if Set is false or Null is true.
func (o *OptionalString) Ptr() *string {
	if !o.Set || o.Null {
		return nil
	}

	v := o.Value
	return &v
}

// NewOptionalString returns a new OptionalString with the given value.
func NewOptionalString(v string) OptionalString {
	return OptionalString{v, false, true}
}

// NewOptionalStringPtr returns a new OptionalString with the given value.
// If the value is nil, the returned OptionalString will be set and null.
func NewOptionalStringPtr(v *string) OptionalString {
	if v == nil {
		return OptionalString{
			Null: true,
			Set:  true,
		}
	}

	return OptionalString{*v, false, true}
}

// OptionalInt represents an optional int argument that may be null. See OptionalString.
type OptionalInt struct {
	Value int
	Null  bool
	Set   bool
}

// Ptr returns a pointer to the underlying value. Returns nil if Set is false or Null is true.
func (o *OptionalInt) Ptr() *int {
	if !o.Set || o.Null {
		return nil
	}

	v := o.Value
	return &v
}

// NewOptionalInt returns a new OptionalInt with the given value.
func NewOptionalInt(v int) OptionalInt {
	return OptionalInt{v, false, true}
}

// NewOptionalIntPtr returns a new OptionalInt with the given value.
// If the value is nil, the returned OptionalInt will be set and null.
func NewOptionalIntPtr(v *int) OptionalInt {
	if v == nil {
		return OptionalInt{
			Null: true,
			Set:  true,
		}
	}

	return OptionalInt{*v, false, true}
}

// StringPtr returns a pointer to a string representation of the value.
// Returns nil if Set is false or null is true.
func (o *OptionalInt) StringPtr() *string {
	if !o.Set || o.Null {
		return nil
	}

	v := strconv.Itoa(o.Value)
	return &v
}

// OptionalInt64 represents an optional int64 argument that may be null. See OptionalString.
type OptionalInt64 struct {
	Value int64
	Null  bool
	Set   bool
}

// Ptr returns a pointer to the underlying value. Returns nil if Set is false or Null is true.
func (o *OptionalInt64) Ptr() *int64 {
	if !o.Set || o.Null {
		return nil
	}

	v := o.Value
	return &v
}

// NewOptionalInt64 returns a new OptionalInt64 with the given value.
func NewOptionalInt64(v int64) OptionalInt64 {
	return OptionalInt64{v, false, true}
}

// NewOptionalInt64Ptr returns a new OptionalInt64 with the given value.
// If the value is nil, the returned OptionalInt64 will be set and null.
func NewOptionalInt64Ptr(v *int64) OptionalInt64 {
	if v == nil {
		return OptionalInt64{
			Null: true,
			Set:  true,
		}
	}

	return OptionalInt64{*v, false, true}
}

// OptionalBool represents an optional int64 argument that may be null. See OptionalString.
type OptionalBool struct {
	Value bool
	Null  bool
	Set   bool
}

func (o *OptionalBool) Ptr() *bool {
	if !o.Set || o.Null {
		return nil
	}

	v := o.Value
	return &v
}

// NewOptionalBool returns a new OptionalBool with the given value.
func NewOptionalBool(v bool) OptionalBool {
	return OptionalBool{v, false, true}
}

// NewOptionalBoolPtr returns a new OptionalBool with the given value.
// If the value is nil, the returned OptionalBool will be set and null.
func NewOptionalBoolPtr(v *bool) OptionalBool {
	if v == nil {
		return OptionalBool{
			Null: true,
			Set:  true,
		}
	}

	return OptionalBool{*v, false, true}
}

// OptionalBool represents an optional float64 argument that may be null. See OptionalString.
type OptionalFloat64 struct {
	Value float64
	Null  bool
	Set   bool
}

// Ptr returns a pointer to the underlying value. Returns nil if Set is false or Null is true.
func (o *OptionalFloat64) Ptr() *float64 {
	if !o.Set || o.Null {
		return nil
	}

	v := o.Value
	return &v
}

// NewOptionalFloat64 returns a new OptionalFloat64 with the given value.
func NewOptionalFloat64(v float64) OptionalFloat64 {
	return OptionalFloat64{v, false, true}
}

// OptionalDate represents an optional date argument that may be null. See OptionalString.
type OptionalDate struct {
	Value Date
	Null  bool
	Set   bool
}

// Ptr returns a pointer to the underlying value. Returns nil if Set is false or Null is true.
func (o *OptionalDate) Ptr() *Date {
	if !o.Set || o.Null {
		return nil
	}

	v := o.Value
	return &v
}

// NewOptionalDate returns a new OptionalDate with the given value.
func NewOptionalDate(v Date) OptionalDate {
	return OptionalDate{v, false, true}
}

// NewOptionalBoolPtr returns a new OptionalDate with the given value.
// If the value is nil, the returned OptionalDate will be set and null.
func NewOptionalDatePtr(v *Date) OptionalDate {
	if v == nil {
		return OptionalDate{
			Null: true,
			Set:  true,
		}
	}

	return OptionalDate{*v, false, true}
}

// OptionalTime represents an optional time argument that may be null. See OptionalString.
type OptionalTime struct {
	Value time.Time
	Null  bool
	Set   bool
}

// NewOptionalTime returns a new OptionalTime with the given value.
func NewOptionalTime(v time.Time) OptionalTime {
	return OptionalTime{v, false, true}
}

// NewOptionalTimePtr returns a new OptionalTime with the given value.
// If the value is nil, the returned OptionalTime will be set and null.
func NewOptionalTimePtr(v *time.Time) OptionalTime {
	if v == nil {
		return OptionalTime{
			Null: true,
			Set:  true,
		}
	}

	return OptionalTime{*v, false, true}
}

// Ptr returns a pointer to the underlying value. Returns nil if Set is false or Null is true.
func (o *OptionalTime) Ptr() *time.Time {
	if !o.Set || o.Null {
		return nil
	}

	v := o.Value
	return &v
}
