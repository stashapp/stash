package sliceutil

import "reflect"

// Exclude removes all instances of any value in toExclude from the vs
// slice. It returns the new or unchanged slice.
func Exclude[T comparable](vs []T, toExclude []T) []T {
	var ret []T
	for _, v := range vs {
		if !Include(toExclude, v) {
			ret = append(ret, v)
		}
	}

	return ret
}

func Index[T comparable](vs []T, t T) int {
	for i, v := range vs {
		if v == t {
			return i
		}
	}
	return -1
}

func Include[T comparable](vs []T, t T) bool {
	return Index(vs, t) >= 0
}

// IntAppendUnique appends toAdd to the vs int slice if toAdd does not already
// exist in the slice. It returns the new or unchanged int slice.
func AppendUnique[T comparable](vs []T, toAdd T) []T {
	if Include(vs, toAdd) {
		return vs
	}

	return append(vs, toAdd)
}

// IntAppendUniques appends a slice of values to the vs slice. It only
// appends values that do not already exist in the slice. It returns the new or
// unchanged slice.
func AppendUniques[T comparable](vs []T, toAdd []T) []T {
	for _, v := range toAdd {
		vs = AppendUnique(vs, v)
	}

	return vs
}

// SliceSame returns true if the two provided lists have the same elements,
// regardless of order. Panics if either parameter is not a slice.
func SliceSame(a, b interface{}) bool {
	v1 := reflect.ValueOf(a)
	v2 := reflect.ValueOf(b)

	if (v1.IsValid() && v1.Kind() != reflect.Slice) || (v2.IsValid() && v2.Kind() != reflect.Slice) {
		panic("not a slice")
	}

	v1Len := 0
	v2Len := 0

	v1Valid := v1.IsValid()
	v2Valid := v2.IsValid()

	if v1Valid {
		v1Len = v1.Len()
	}
	if v2Valid {
		v2Len = v2.Len()
	}

	if !v1Valid || !v2Valid {
		return v1Len == v2Len
	}

	if v1Len != v2Len {
		return false
	}

	if v1.Type() != v2.Type() {
		return false
	}

	visited := make(map[int]bool)
	for i := 0; i < v1.Len(); i++ {
		found := false
		for j := 0; j < v2.Len(); j++ {
			if visited[j] {
				continue
			}
			if reflect.DeepEqual(v1.Index(i).Interface(), v2.Index(j).Interface()) {
				found = true
				visited[j] = true
				break
			}
		}

		if !found {
			return false
		}
	}

	return true
}
