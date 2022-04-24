package models

import (
	"strconv"
	"time"
)

type OptionalString struct {
	Value string
	Null  bool
	Set   bool
}

func (o *OptionalString) Ptr() *string {
	if !o.Set || o.Null {
		return nil
	}

	v := o.Value
	return &v
}

func NewOptionalString(v string) OptionalString {
	return OptionalString{v, false, true}
}

func NewOptionalStringPtr(v *string) OptionalString {
	if v == nil {
		return OptionalString{
			Null: true,
			Set:  true,
		}
	}

	return OptionalString{*v, false, true}
}

type OptionalInt struct {
	Value int
	Null  bool
	Set   bool
}

func (o *OptionalInt) Ptr() *int {
	if !o.Set || o.Null {
		return nil
	}

	v := o.Value
	return &v
}

func NewOptionalInt(v int) OptionalInt {
	return OptionalInt{v, false, true}
}

func NewOptionalIntPtr(v *int) OptionalInt {
	if v == nil {
		return OptionalInt{
			Null: true,
			Set:  true,
		}
	}

	return OptionalInt{*v, false, true}
}

func (o *OptionalInt) StringPtr() *string {
	if !o.Set || o.Null {
		return nil
	}

	v := strconv.Itoa(o.Value)
	return &v
}

type OptionalInt64 struct {
	Value int64
	Null  bool
	Set   bool
}

func (o *OptionalInt64) Ptr() *int64 {
	if !o.Set || o.Null {
		return nil
	}

	v := o.Value
	return &v
}

func NewOptionalInt64(v int64) OptionalInt64 {
	return OptionalInt64{v, false, true}
}

func NewOptionalInt64Ptr(v *int64) OptionalInt64 {
	if v == nil {
		return OptionalInt64{
			Null: true,
			Set:  true,
		}
	}

	return OptionalInt64{*v, false, true}
}

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

func NewOptionalBool(v bool) OptionalBool {
	return OptionalBool{v, false, true}
}

func NewOptionalBoolPtr(v *bool) OptionalBool {
	if v == nil {
		return OptionalBool{
			Null: true,
			Set:  true,
		}
	}

	return OptionalBool{*v, false, true}
}

type OptionalFloat64 struct {
	Value float64
	Null  bool
	Set   bool
}

func (o *OptionalFloat64) Ptr() *float64 {
	if !o.Set || o.Null {
		return nil
	}

	v := o.Value
	return &v
}

func NewOptionalFloat64(v float64) OptionalFloat64 {
	return OptionalFloat64{v, false, true}
}

type OptionalDate struct {
	Value Date
	Null  bool
	Set   bool
}

func (o *OptionalDate) Ptr() *Date {
	if !o.Set || o.Null {
		return nil
	}

	v := o.Value
	return &v
}

func NewOptionalDate(v Date) OptionalDate {
	return OptionalDate{v, false, true}
}

func NewOptionalDatePtr(v *Date) OptionalDate {
	if v == nil {
		return OptionalDate{
			Null: true,
			Set:  true,
		}
	}

	return OptionalDate{*v, false, true}
}

type OptionalTime struct {
	Value time.Time
	Null  bool
	Set   bool
}

func NewOptionalTime(v time.Time) OptionalTime {
	return OptionalTime{v, false, true}
}

func NewOptionalTimePtr(v *time.Time) OptionalTime {
	if v == nil {
		return OptionalTime{
			Null: true,
			Set:  true,
		}
	}

	return OptionalTime{*v, false, true}
}

func (o *OptionalTime) Ptr() *time.Time {
	if !o.Set || o.Null {
		return nil
	}

	v := o.Value
	return &v
}
