package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

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

	return time.Time{}, fmt.Errorf("ParseDateStringAsTime failed: dateString <%s>", dateString)
}

// ParseYearRangeString parses a year range string into start and end year integers.
// Supported formats: "YYYY", "YYYY - YYYY", "YYYY-YYYY", "YYYY -", "- YYYY", "YYYY-present".
// Returns nil for start/end if not present in the string.
func ParseYearRangeString(s string) (start *int, end *int, err error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, nil, fmt.Errorf("empty year range string")
	}

	// normalize "present" to empty end
	lower := strings.ToLower(s)
	lower = strings.ReplaceAll(lower, "present", "")

	// split on "-" if it contains one
	var parts []string
	if strings.Contains(lower, "-") {
		parts = strings.SplitN(lower, "-", 2)
	} else {
		// single value, treat as start year
		year, err := parseYear(lower)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid year range %q: %w", s, err)
		}
		return &year, nil, nil
	}

	startStr := strings.TrimSpace(parts[0])
	endStr := strings.TrimSpace(parts[1])

	if startStr != "" {
		y, err := parseYear(startStr)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid start year in %q: %w", s, err)
		}
		start = &y
	}

	if endStr != "" {
		y, err := parseYear(endStr)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid end year in %q: %w", s, err)
		}
		end = &y
	}

	if start == nil && end == nil {
		return nil, nil, fmt.Errorf("could not parse year range %q", s)
	}

	return start, end, nil
}

func parseYear(s string) (int, error) {
	s = strings.TrimSpace(s)
	year, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("invalid year %q: %w", s, err)
	}
	if year < 1900 || year > 2200 {
		return 0, fmt.Errorf("year %d out of reasonable range", year)
	}
	return year, nil
}

func FormatYearRange(start *int, end *int) string {
	switch {
	case start == nil && end == nil:
		return ""
	case end == nil:
		return fmt.Sprintf("%d -", *start)
	case start == nil:
		return fmt.Sprintf("- %d", *end)
	default:
		return fmt.Sprintf("%d - %d", *start, *end)
	}
}
