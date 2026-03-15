package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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

func TestFormatYearRange(t *testing.T) {
	datePtr := func(v int) *Date {
		date := DateFromYear(v)
		return &date
	}

	tests := []struct {
		name  string
		start *Date
		end   *Date
		want  string
	}{
		{"both nil", nil, nil, ""},
		{"only start", datePtr(2005), nil, "2005 -"},
		{"only end", nil, datePtr(2010), "- 2010"},
		{"start and end", datePtr(2005), datePtr(2010), "2005 - 2010"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatYearRange(tt.start, tt.end)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFormatYearRangeString(t *testing.T) {
	stringPtr := func(v string) *string { return &v }

	tests := []struct {
		name  string
		start *string
		end   *string
		want  string
	}{
		{"both nil", nil, nil, ""},
		{"only start", stringPtr("2005"), nil, "2005 -"},
		{"only end", nil, stringPtr("2010"), "- 2010"},
		{"start and end", stringPtr("2005"), stringPtr("2010"), "2005 - 2010"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatYearRangeString(tt.start, tt.end)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParseYearRangeString(t *testing.T) {
	intPtr := func(v int) *int { return &v }

	tests := []struct {
		name      string
		input     string
		wantStart *int
		wantEnd   *int
		wantErr   bool
	}{
		{"single year", "2005", intPtr(2005), nil, false},
		{"year range with spaces", "2005 - 2010", intPtr(2005), intPtr(2010), false},
		{"year range no spaces", "2005-2010", intPtr(2005), intPtr(2010), false},
		{"year dash open", "2005 -", intPtr(2005), nil, false},
		{"year dash open no space", "2005-", intPtr(2005), nil, false},
		{"dash year", "- 2010", nil, intPtr(2010), false},
		{"year present", "2005-present", intPtr(2005), nil, false},
		{"year Present caps", "2005 - Present", intPtr(2005), nil, false},
		{"whitespace padding", "  2005 - 2010  ", intPtr(2005), intPtr(2010), false},
		{"empty string", "", nil, nil, true},
		{"garbage", "not a year", nil, nil, true},
		{"partial garbage start", "abc - 2010", nil, nil, true},
		{"partial garbage end", "2005 - abc", nil, nil, true},
		{"year out of range", "1800", nil, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start, end, err := ParseYearRangeString(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			if tt.wantStart != nil {
				assert.NotNil(t, start)
				assert.Equal(t, *tt.wantStart, start.Time.Year())
			} else {
				assert.Nil(t, start)
			}
			if tt.wantEnd != nil {
				assert.NotNil(t, end)
				assert.Equal(t, *tt.wantEnd, end.Time.Year())
			} else {
				assert.Nil(t, end)
			}
		})
	}
}
