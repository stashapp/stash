package stringslice

import (
	"strconv"
	"strings"
)

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

// FromString converts a string to a slice of strings, splitting on the sep character.
// Unlike strings.Split, this function will also trim whitespace from the resulting strings.
func FromString(s string, sep string) []string {
	v := strings.Split(s, ",")
	for i, vv := range v {
		v[i] = strings.TrimSpace(vv)
	}
	return v
}

// Unique returns a slice containing only unique values from the provided slice.
// The comparison is case-insensitive.
func UniqueFold(s []string) []string {
	seen := make(map[string]struct{}, len(s))
	ret := make([]string, 0, len(s))
	for _, v := range s {
		if _, exists := seen[strings.ToLower(v)]; exists {
			continue
		}
		seen[strings.ToLower(v)] = struct{}{}
		ret = append(ret, v)
	}
	return ret
}

// UniqueExcludeFold returns a deduplicated slice of strings with the excluded string removed.
// The comparison is case-insensitive.
func UniqueExcludeFold(values []string, exclude string) []string {
	seen := make(map[string]struct{}, len(values))
	seen[strings.ToLower(exclude)] = struct{}{}
	ret := make([]string, 0, len(values))
	for _, v := range values {
		vLower := strings.ToLower(v)
		if _, exists := seen[vLower]; exists {
			continue
		}
		seen[vLower] = struct{}{}
		ret = append(ret, v)
	}
	return ret
}

// TrimSpace trims whitespace from each string in a slice.
func TrimSpace(s []string) []string {
	for i, v := range s {
		s[i] = strings.TrimSpace(v)
	}
	return s
}
