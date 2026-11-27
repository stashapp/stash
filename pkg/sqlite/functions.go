package sqlite

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
)

func durationToTinyIntFn(str string) (int64, error) {
	splits := strings.Split(str, ":")

	if len(splits) > 3 {
		return 0, nil
	}

	seconds := 0
	factor := 1
	for len(splits) > 0 {
		// pop the last split
		var thisSplit string
		thisSplit, splits = splits[len(splits)-1], splits[:len(splits)-1]

		thisInt, err := strconv.Atoi(thisSplit)
		if err != nil {
			return 0, nil
		}

		seconds += factor * thisInt
		factor *= 60
	}

	return int64(seconds), nil
}

func basenameFn(str string) (string, error) {
	return filepath.Base(str), nil
}

// custom SQLite function to enable case-insensitive searches
// that properly handle unicode characters
func lowerUnicodeFn(str interface{}) (string, error) {
	// handle NULL values
	if str == nil {
		return "", nil
	}

	// handle different types
	switch v := str.(type) {
	case string:
		return strings.ToLower(v), nil
	case int64:
		// convert int64 to string (for phash fingerprints)
		return strings.ToLower(strconv.FormatInt(v, 10)), nil
	case []byte:
		// handle BLOB type if needed
		return strings.ToLower(string(v)), nil
	default:
		// for any other type, try converting to string
		return strings.ToLower(fmt.Sprintf("%v", v)), nil
	}
}
