package api

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertMapJSONNumbers(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]interface{}
		expected map[string]interface{}
	}{
		{
			name: "Convert JSON numbers to numbers",
			input: map[string]interface{}{
				"int":    json.Number("12"),
				"float":  json.Number("12.34"),
				"string": "foo",
			},
			expected: map[string]interface{}{
				"int":    int64(12),
				"float":  12.34,
				"string": "foo",
			},
		},
		{
			name: "Convert JSON numbers to numbers in nested maps",
			input: map[string]interface{}{
				"foo": map[string]interface{}{
					"int":           json.Number("56"),
					"float":         json.Number("56.78"),
					"nested-string": "bar",
				},
				"int":    json.Number("12"),
				"float":  json.Number("12.34"),
				"string": "foo",
			},
			expected: map[string]interface{}{
				"foo": map[string]interface{}{
					"int":           int64(56),
					"float":         56.78,
					"nested-string": "bar",
				},
				"int":    int64(12),
				"float":  12.34,
				"string": "foo",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertMapJSONNumbers(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
