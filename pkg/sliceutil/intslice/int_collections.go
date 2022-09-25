package intslice

import "strconv"

// IntIndex returns the first index of the provided int value in the provided
// int slice. It returns -1 if it is not found.
func IntIndex(vs []int, t int) int {
	for i, v := range vs {
		if v == t {
			return i
		}
	}
	return -1
}

// IntInclude returns true if the provided int value exists in the provided int
// slice.
func IntInclude(vs []int, t int) bool {
	return IntIndex(vs, t) >= 0
}

// IntAppendUnique appends toAdd to the vs int slice if toAdd does not already
// exist in the slice. It returns the new or unchanged int slice.
func IntAppendUnique(vs []int, toAdd int) []int {
	if IntInclude(vs, toAdd) {
		return vs
	}

	return append(vs, toAdd)
}

// IntAppendUniques appends a slice of int values to the vs int slice. It only
// appends values that do not already exist in the slice. It returns the new or
// unchanged int slice.
func IntAppendUniques(vs []int, toAdd []int) []int {
	for _, v := range toAdd {
		vs = IntAppendUnique(vs, v)
	}

	return vs
}

// IntExclude removes all instances of any value in toExclude from the vs int
// slice. It returns the new or unchanged int slice.
func IntExclude(vs []int, toExclude []int) []int {
	var ret []int
	for _, v := range vs {
		if !IntInclude(toExclude, v) {
			ret = append(ret, v)
		}
	}

	return ret
}

// IntIntercect returns a slice of ints containing values that exist in both provided slices.
func IntIntercect(v1, v2 []int) []int {
	var ret []int
	for _, v := range v1 {
		if IntInclude(v2, v) {
			ret = append(ret, v)
		}
	}

	return ret
}

// IntSliceToStringSlice converts a slice of ints to a slice of strings.
func IntSliceToStringSlice(ss []int) []string {
	ret := make([]string, len(ss))
	for i, v := range ss {
		ret[i] = strconv.Itoa(v)
	}

	return ret
}
