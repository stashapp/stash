package api

import "math"

// An enum https://golang.org/ref/spec#Iota
const (
	create = iota // 0
	update = iota // 1
)

// #1572 - Inf and NaN values cause the JSON marshaller to fail
// Return nil for these values
func handleFloat64(v float64) *float64 {
	if math.IsInf(v, 0) || math.IsNaN(v) {
		return nil
	}

	return &v
}
