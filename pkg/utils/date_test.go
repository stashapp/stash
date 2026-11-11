package utils

import (
	"testing"
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
		{"Unpadded date", "2014-1-2", false},
		{"Compact date", "20140102", false},
		{"Dotted date", "2014.01.02", false},
		{"Two-digit year date", "14-01-02", false},
		{"Two-digit year dotted date", "14.01.02", false},

		// Invalid formats
		{"Invalid format", "not-a-date", true},
		{"Empty string", "", true},
		{"Year-Month", "2006-08", true},
		{"Year only", "2014", true},
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
