package performer

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseCareerLength parses a career_length string into start and end year integers.
// Supported formats: "YYYY", "YYYY - YYYY", "YYYY-YYYY", "YYYY -", "- YYYY", "YYYY-present".
// Returns nil for start/end if not present in the string.
func ParseCareerLength(s string) (start *int, end *int, err error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, nil, fmt.Errorf("empty career length string")
	}

	// normalize "present" to empty end
	lower := strings.ToLower(s)
	lower = strings.ReplaceAll(lower, "present", "")

	// split on " - ", "-", or " -"
	var parts []string
	switch {
	case strings.Contains(lower, " - "):
		parts = strings.SplitN(lower, " - ", 2)
	case strings.Contains(lower, "-"):
		parts = strings.SplitN(lower, "-", 2)
	default:
		// single value, treat as start year
		year, err := parseYear(lower)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid career length %q: %w", s, err)
		}
		return &year, nil, nil
	}

	startStr := strings.TrimSpace(parts[0])
	endStr := strings.TrimSpace(parts[1])

	if startStr != "" {
		y, err := parseYear(startStr)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid career start in %q: %w", s, err)
		}
		start = &y
	}

	if endStr != "" {
		y, err := parseYear(endStr)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid career end in %q: %w", s, err)
		}
		end = &y
	}

	if start == nil && end == nil {
		return nil, nil, fmt.Errorf("could not parse career length %q", s)
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
