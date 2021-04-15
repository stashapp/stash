package utils

import (
	"fmt"
	"time"
)

const railsTimeLayout = "2006-01-02 15:04:05 MST"

func GetYMDFromDatabaseDate(dateString string) string {
	result, _ := ParseDateStringAsFormat(dateString, "2006-01-02")
	return result
}

func ParseDateStringAsFormat(dateString string, format string) (string, error) {
	t, e := ParseDateStringAsTime(dateString)
	if e == nil {
		return t.Format(format), e
	}
	return "", fmt.Errorf("ParseDateStringAsFormat failed: dateString <%s>, format <%s>", dateString, format)
}

func ParseDateStringAsTime(dateString string) (time.Time, error) {
	// https://stackoverflow.com/a/20234207 WTF?

	t, e := time.Parse(time.RFC3339, dateString)
	if e == nil {
		return t, nil
	}

	t, e = time.Parse("2006-01-02", dateString)
	if e == nil {
		return t, nil
	}

	t, e = time.Parse("2006-01-02 15:04:05", dateString)
	if e == nil {
		return t, nil
	}

	t, e = time.Parse(railsTimeLayout, dateString)
	if e == nil {
		return t, nil
	}

	return time.Time{}, fmt.Errorf("ParseDateStringAsTime failed: dateString <%s>", dateString)
}
