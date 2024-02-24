package sliceutil

// Index returns the first index of the provided value in the provided
// slice. It returns -1 if it is not found.
func Index[T comparable](vs []T, t T) int {
	for i, v := range vs {
		if v == t {
			return i
		}
	}
	return -1
}

// Contains returns whether the vs slice contains t.
func Contains[T comparable](vs []T, t T) bool {
	return Index(vs, t) >= 0
}

// AppendUnique appends toAdd to the vs slice if toAdd does not already
// exist in the slice. It returns the new or unchanged slice.
func AppendUnique[T comparable](vs []T, toAdd T) []T {
	if Contains(vs, toAdd) {
		return vs
	}

	return append(vs, toAdd)
}

// AppendUniques appends a slice of values to the vs slice. It only
// appends values that do not already exist in the slice.
// It returns the new or unchanged slice.
func AppendUniques[T comparable](vs []T, toAdd []T) []T {
	for _, v := range toAdd {
		vs = AppendUnique(vs, v)
	}

	return vs
}

// Exclude returns a copy of the vs slice, excluding all values
// that are also present in the toExclude slice.
func Exclude[T comparable](vs []T, toExclude []T) []T {
	var ret []T
	for _, v := range vs {
		if !Contains(toExclude, v) {
			ret = append(ret, v)
		}
	}

	return ret
}

// Unique returns a copy of the vs slice, with non-unique values removed.
func Unique[T comparable](vs []T) []T {
	distinctValues := make(map[T]struct{})
	var ret []T
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
	var ret []T
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
		if Contains(b, v) {
			ret = append(ret, v)
		}
	}

	return ret
}

// NotIntersect returns a slice containing values that do not exist in both provided slices.
func NotIntersect[T comparable](a []T, b []T) []T {
	var ret []T
	for _, v := range a {
		if !Contains(b, v) {
			ret = append(ret, v)
		}
	}

	for _, v := range b {
		if !Contains(a, v) {
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

// Filter returns the result of applying f to each element of the vs slice.
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
		vv := v
		ret[i] = &vv
	}
	return ret
}
