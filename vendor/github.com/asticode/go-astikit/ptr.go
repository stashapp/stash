package astikit

import "time"

// BoolPtr transforms a bool into a *bool
func BoolPtr(i bool) *bool {
	return &i
}

// BytePtr transforms a byte into a *byte
func BytePtr(i byte) *byte {
	return &i
}

// DurationPtr transforms a time.Duration into a *time.Duration
func DurationPtr(i time.Duration) *time.Duration {
	return &i
}

// Float64Ptr transforms a float64 into a *float64
func Float64Ptr(i float64) *float64 {
	return &i
}

// IntPtr transforms an int into an *int
func IntPtr(i int) *int {
	return &i
}

// Int64Ptr transforms an int64 into an *int64
func Int64Ptr(i int64) *int64 {
	return &i
}

// StrSlicePtr transforms a []string into a *[]string
func StrSlicePtr(i []string) *[]string {
	return &i
}

// StrPtr transforms a string into a *string
func StrPtr(i string) *string {
	return &i
}

// TimePtr transforms a time.Time into a *time.Time
func TimePtr(i time.Time) *time.Time {
	return &i
}

// UInt8Ptr transforms a uint8 into a *uint8
func UInt8Ptr(i uint8) *uint8 {
	return &i
}

// UInt32Ptr transforms a uint32 into a *uint32
func UInt32Ptr(i uint32) *uint32 {
	return &i
}
