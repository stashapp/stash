package utils

import (
	"testing"
	"time"
)

func TestParseDateStringAsTime(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
	}{
		// Full date formats (existing support)
		{"RFC3339", "2014-01-02T15:04:05Z", false},
		{"Date only", "2014-01-02", false},
		{"Date with time", "2014-01-02 15:04:05", false},

		// Partial date formats (new support)
		{"Year-Month", "2006-08", false},
		{"Year only", "2014", false},

		// Invalid formats
		{"Invalid format", "not-a-date", true},
		{"Empty string", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseDateStringAsTime(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for input %q, but got none", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for input %q: %v", tt.input, err)
				}
				if result.IsZero() {
					t.Errorf("Expected non-zero time for input %q", tt.input)
				}
			}
		})
	}
}

func TestParseDateStringAsTime_YearOnly(t *testing.T) {
	result, err := ParseDateStringAsTime("2014")
	if err != nil {
		t.Fatalf("Failed to parse year-only date: %v", err)
	}

	if result.Year() != 2014 {
		t.Errorf("Expected year 2014, got %d", result.Year())
	}
	if result.Month() != time.January {
		t.Errorf("Expected month January, got %s", result.Month())
	}
	if result.Day() != 1 {
		t.Errorf("Expected day 1, got %d", result.Day())
	}
}

func TestParseDateStringAsTime_YearMonth(t *testing.T) {
	result, err := ParseDateStringAsTime("2006-08")
	if err != nil {
		t.Fatalf("Failed to parse year-month date: %v", err)
	}

	if result.Year() != 2006 {
		t.Errorf("Expected year 2006, got %d", result.Year())
	}
	if result.Month() != time.August {
		t.Errorf("Expected month August, got %s", result.Month())
	}
	if result.Day() != 1 {
		t.Errorf("Expected day 1, got %d", result.Day())
	}
}