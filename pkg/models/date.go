package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/stashapp/stash/pkg/utils"
)

type DatePrecision int

const (
	// default precision is day
	DatePrecisionDay DatePrecision = iota
	DatePrecisionMonth
	DatePrecisionYear
)

// Date wraps a time.Time with a format of "YYYY-MM-DD"
type Date struct {
	time.Time
	Precision DatePrecision
}

var dateFormatPrecision = []string{
	"2006-01-02",
	"2006-01",
	"2006",
}

func (d Date) String() string {
	return d.Format(dateFormatPrecision[d.Precision])
}

func (d Date) After(o Date) bool {
	return d.Time.After(o.Time)
}

// ParseDate tries to parse the input string into a date using utils.ParseDateStringAsTime.
// If that fails, it attempts to parse the string with decreasing precision (month, then year).
// It returns a Date struct with the appropriate precision set, or an error if all parsing attempts fail.
func ParseDate(s string) (Date, error) {
	var errs []error

	// default parse to day precision
	ret, err := utils.ParseDateStringAsTime(s)
	if err == nil {
		return Date{Time: ret, Precision: DatePrecisionDay}, nil
	}

	errs = append(errs, err)

	// try month and year precision
	for i, format := range dateFormatPrecision[1:] {
		ret, err := time.Parse(format, s)
		if err == nil {
			return Date{Time: ret, Precision: DatePrecision(i + 1)}, nil
		}
		errs = append(errs, err)
	}

	return Date{}, fmt.Errorf("failed to parse date %q: %v", s, errs)
}

func DateFromYear(year int) Date {
	return Date{
		Time:      time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC),
		Precision: DatePrecisionYear,
	}
}

func FormatYearRange(start *Date, end *Date) string {
	var (
		startStr, endStr string
	)

	if start != nil {
		startStr = start.Format(dateFormatPrecision[DatePrecisionYear])
	}

	if end != nil {
		endStr = end.Format(dateFormatPrecision[DatePrecisionYear])
	}

	switch {
	case startStr == "" && endStr == "":
		return ""
	case endStr == "":
		return fmt.Sprintf("%s -", startStr)
	case startStr == "":
		return fmt.Sprintf("- %s", endStr)
	default:
		return fmt.Sprintf("%s - %s", startStr, endStr)
	}
}

func FormatYearRangeString(start *string, end *string) string {
	switch {
	case start == nil && end == nil:
		return ""
	case end == nil:
		return fmt.Sprintf("%s -", *start)
	case start == nil:
		return fmt.Sprintf("- %s", *end)
	default:
		return fmt.Sprintf("%s - %s", *start, *end)
	}
}

// ParseYearRangeString parses a year range string into start and end year integers.
// Supported formats: "YYYY", "YYYY - YYYY", "YYYY-YYYY", "YYYY -", "- YYYY", "YYYY-present".
// Returns nil for start/end if not present in the string.
func ParseYearRangeString(s string) (start *Date, end *Date, err error) {
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
		return year, nil, nil
	}

	startStr := strings.TrimSpace(parts[0])
	endStr := strings.TrimSpace(parts[1])

	if startStr != "" {
		y, err := parseYear(startStr)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid start year in %q: %w", s, err)
		}
		start = y
	}

	if endStr != "" {
		y, err := parseYear(endStr)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid end year in %q: %w", s, err)
		}
		end = y
	}

	if start == nil && end == nil {
		return nil, nil, fmt.Errorf("could not parse year range %q", s)
	}

	return start, end, nil
}

func parseYear(s string) (*Date, error) {
	ret, err := ParseDate(s)
	if err != nil {
		return nil, fmt.Errorf("parsing year %q: %w", s, err)
	}

	year := ret.Time.Year()
	if year < 1900 || year > 2200 {
		return nil, fmt.Errorf("year %d out of reasonable range", year)
	}

	return &ret, nil
}
