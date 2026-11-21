package utils

import (
	"fmt"
	"time"
)

func ParseDateStringAsTime(dateString string) (time.Time, error) {
	// https://stackoverflow.com/a/20234207 WTF?

	t, e := time.Parse(time.RFC3339, dateString)
	if e == nil {
		return t, nil
	}

	for _, format := range []string{
		"2006-01-02",
		"2006-1-2",
		"20060102",
		"2006.01.02",
		"2006.1.2",
		"06-01-02",
		"06-1-2",
		"06.01.02",
		"06.1.2",
	} {
		t, e = time.Parse(format, dateString)
		if e == nil {
			return t, nil
		}
	}

	t, e = time.Parse("2006-01-02 15:04:05", dateString)
	if e == nil {
		return t, nil
	}

	return time.Time{}, fmt.Errorf("ParseDateStringAsTime failed: dateString <%s>", dateString)
}
