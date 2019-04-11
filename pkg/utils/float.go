package utils

import "math"

// IsValidFloat64 ensures the given value is a valid number (not NaN) which is not equal to 0
func IsValidFloat64(value float64) bool {
	return !math.IsNaN(value) && value != 0
}
