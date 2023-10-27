package sliceutil

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
