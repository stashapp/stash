// Package sliceutil provides utilities for working with slices.
package sliceutil

import (
	"slices"
)

// AppendUnique appends toAdd to the vs slice if toAdd does not already
// exist in the slice. It returns the new or unchanged slice.
func AppendUnique[T comparable](vs []T, toAdd T) []T {
	if slices.Contains(vs, toAdd) {
		return vs
	}

	return append(vs, toAdd)
}

// AppendUniques appends a slice of values to the vs slice. It only
// appends values that do not already exist in the slice.
// It returns the new or unchanged slice.
func AppendUniques[T comparable](vs []T, toAdd []T) []T {
	if len(toAdd) == 0 {
		return vs
	}

	// Extend the slice's capacity to avoid multiple re-allocations even in the worst case
	vs = slices.Grow(vs, len(toAdd))

	for _, v := range toAdd {
		vs = AppendUnique(vs, v)
	}

	return vs
}

// Exclude returns a copy of the vs slice, excluding all values
// that are also present in the toExclude slice.
func Exclude[T comparable](vs []T, toExclude []T) []T {
	ret := make([]T, 0, len(vs))
	for _, v := range vs {
		if !slices.Contains(toExclude, v) {
			ret = append(ret, v)
		}
	}

	return ret
}

// Unique returns a copy of the vs slice, with non-unique values removed.
func Unique[T comparable](vs []T) []T {
	distinctValues := make(map[T]struct{}, len(vs))
	ret := make([]T, 0, len(vs))
	for _, v := range vs {
		if _, exists := distinctValues[v]; !exists {
			distinctValues[v] = struct{}{}
			ret = append(ret, v)
		}
	}
	return ret
}

// Delete returns a copy of the vs slice with toDel values removed.
func Delete[T comparable](vs []T, toDel T) []T {
	ret := make([]T, 0, len(vs))
	for _, v := range vs {
		if v != toDel {
			ret = append(ret, v)
		}
	}
	return ret
}

// Intersect returns a slice containing values that exist in both provided slices.
func Intersect[T comparable](a []T, b []T) []T {
	var ret []T
	for _, v := range a {
		if slices.Contains(b, v) {
			ret = append(ret, v)
		}
	}

	return ret
}

// NotIntersect returns a slice containing values that do not exist in both provided slices.
func NotIntersect[T comparable](a []T, b []T) []T {
	var ret []T
	for _, v := range a {
		if !slices.Contains(b, v) {
			ret = append(ret, v)
		}
	}

	for _, v := range b {
		if !slices.Contains(a, v) {
			ret = append(ret, v)
		}
	}

	return ret
}

// SliceSame returns true if the two provided slices have equal elements,
// regardless of order.
func SliceSame[T comparable](a []T, b []T) bool {
	if len(a) != len(b) {
		return false
	}

	visited := make(map[int]struct{})
	for i := range a {
		found := false
		for j := range b {
			if _, exists := visited[j]; exists {
				continue
			}
			if a[i] == b[j] {
				found = true
				visited[j] = struct{}{}
				break
			}
		}

		if !found {
			return false
		}
	}

	return true
}

// Filter returns a slice containing the elements of the vs slice
// that meet the condition specified by f.
func Filter[T any](vs []T, f func(T) bool) []T {
	var ret []T
	for _, v := range vs {
		if f(v) {
			ret = append(ret, v)
		}
	}
	return ret
}

// Map returns the result of applying f to each element of the vs slice.
func Map[T any, V any](vs []T, f func(T) V) []V {
	ret := make([]V, len(vs))
	for i, v := range vs {
		ret[i] = f(v)
	}
	return ret
}

func PtrsToValues[T any](vs []*T) []T {
	ret := make([]T, len(vs))
	for i, v := range vs {
		ret[i] = *v
	}
	return ret
}

func ValuesToPtrs[T any](vs []T) []*T {
	ret := make([]*T, len(vs))
	for i, v := range vs {
		// We can do this safely because go.mod indicates Go 1.22
		// See: https://go.dev/blog/loopvar-preview
		ret[i] = &v
	}
	return ret
}
