package models

import (
	"testing"
	"time"
)

func TestParseDateStringAsTime(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		output      Date
		expectError bool
	}{
		// Full date formats (existing support)
		{"RFC3339", "2014-01-02T15:04:05Z", Date{Time: time.Date(2014, 1, 2, 15, 4, 5, 0, time.UTC), Precision: DatePrecisionDay}, false},
		{"Date only", "2014-01-02", Date{Time: time.Date(2014, 1, 2, 0, 0, 0, 0, time.UTC), Precision: DatePrecisionDay}, false},
		{"Date with time", "2014-01-02 15:04:05", Date{Time: time.Date(2014, 1, 2, 15, 4, 5, 0, time.UTC), Precision: DatePrecisionDay}, false},

		// Partial date formats (new support)
		{"Year-Month", "2006-08", Date{Time: time.Date(2006, 8, 1, 0, 0, 0, 0, time.UTC), Precision: DatePrecisionMonth}, false},
		{"Year only", "2014", Date{Time: time.Date(2014, 1, 1, 0, 0, 0, 0, time.UTC), Precision: DatePrecisionYear}, false},

		// Invalid formats
		{"Invalid format", "not-a-date", Date{}, true},
		{"Empty string", "", Date{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseDate(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for input %q, but got none", tt.input)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error for input %q: %v", tt.input, err)
				return
			}

			if !result.Time.Equal(tt.output.Time) || result.Precision != tt.output.Precision {
				t.Errorf("For input %q, expected output %+v, got %+v", tt.input, tt.output, result)
			}
		})
	}
}
