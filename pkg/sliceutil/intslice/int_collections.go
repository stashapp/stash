package intslice

import "strconv"

// IntSliceToStringSlice converts a slice of ints to a slice of strings.
func IntSliceToStringSlice(ss []int) []string {
	ret := make([]string, len(ss))
	for i, v := range ss {
		ret[i] = strconv.Itoa(v)
	}

	return ret
}
