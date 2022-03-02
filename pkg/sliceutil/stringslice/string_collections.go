package stringslice

import "strconv"

// https://gobyexample.com/collection-functions

func StrIndex(vs []string, t string) int {
	for i, v := range vs {
		if v == t {
			return i
		}
	}
	return -1
}

func StrInclude(vs []string, t string) bool {
	return StrIndex(vs, t) >= 0
}

func StrFilter(vs []string, f func(string) bool) []string {
	vsf := make([]string, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func StrMap(vs []string, f func(string) string) []string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

// StrAppendUnique appends toAdd to the vs string slice if toAdd does not already
// exist in the slice. It returns the new or unchanged string slice.
func StrAppendUnique(vs []string, toAdd string) []string {
	if StrInclude(vs, toAdd) {
		return vs
	}

	return append(vs, toAdd)
}

// StrAppendUniques appends a slice of string values to the vs string slice. It only
// appends values that do not already exist in the slice. It returns the new or
// unchanged string slice.
func StrAppendUniques(vs []string, toAdd []string) []string {
	for _, v := range toAdd {
		vs = StrAppendUnique(vs, v)
	}

	return vs
}

// StrUnique returns the vs string slice with non-unique values removed.
func StrUnique(vs []string) []string {
	distinctValues := make(map[string]struct{})
	var ret []string
	for _, v := range vs {
		if _, exists := distinctValues[v]; !exists {
			distinctValues[v] = struct{}{}
			ret = append(ret, v)
		}
	}
	return ret
}

// StrDelete returns the vs string slice with toDel values removed.
func StrDelete(vs []string, toDel string) []string {
	var ret []string
	for _, v := range vs {
		if v != toDel {
			ret = append(ret, v)
		}
	}
	return ret
}

// StringSliceToIntSlice converts a slice of strings to a slice of ints.
// Returns an error if any values cannot be parsed.
func StringSliceToIntSlice(ss []string) ([]int, error) {
	ret := make([]int, len(ss))
	for i, v := range ss {
		var err error
		ret[i], err = strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
	}

	return ret, nil
}
