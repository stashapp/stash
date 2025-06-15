package sqlite

import (
	"path/filepath"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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

func titleCaseFn(str string) string {
	return cases.Title(language.Und, cases.NoLower).String(str)
}

func trimTextFn(str string) string {
	return strings.TrimSpace(str)
}

func generateRandomUUIDFn() (string, error) {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return uuid.String(), nil
}
