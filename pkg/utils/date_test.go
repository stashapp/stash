package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
			assert.Equal(t, tt.wantStart, start)
			assert.Equal(t, tt.wantEnd, end)
		})
	}
}

func TestFormatYearRange(t *testing.T) {
	intPtr := func(v int) *int { return &v }

	tests := []struct {
		name  string
		start *int
		end   *int
		want  string
	}{
		{"both nil", nil, nil, ""},
		{"only start", intPtr(2005), nil, "2005 -"},
		{"only end", nil, intPtr(2010), "- 2010"},
		{"start and end", intPtr(2005), intPtr(2010), "2005 - 2010"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatYearRange(tt.start, tt.end)
			assert.Equal(t, tt.want, got)
		})
	}
}
