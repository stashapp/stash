package graphql

import (
	"bytes"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshalID(t *testing.T) {
	marshalID := func(s string) string {
		var buf bytes.Buffer
		MarshalID(s).MarshalGQL(&buf)
		return buf.String()
	}

	assert.Equal(t, `"hello"`, marshalID("hello"))
	assert.Equal(t, `"he\tllo"`, marshalID("he\tllo"))
	assert.Equal(t, `"he\tllo"`, marshalID("he	llo"))
	assert.Equal(t, `"he\nllo"`, marshalID("he\nllo"))
	assert.Equal(t, `"he\r\nllo"`, marshalID("he\r\nllo"))
	assert.Equal(t, `"he\\llo"`, marshalID(`he\llo`))
	assert.Equal(t, `"quotes\"nested\"in\"quotes\""`, marshalID(`quotes"nested"in"quotes"`))
	assert.Equal(t, `"\u0000"`, marshalID("\u0000"))
	assert.Equal(t, "\"\U000fe4ed\"", marshalID("\U000fe4ed"))
	assert.Equal(t, "\"\\u001B\"", marshalID("\u001B"))
}

func TestUnmarshalID(t *testing.T) {
	tests := []struct {
		Name        string
		Input       interface{}
		Expected    string
		ShouldError bool
	}{
		{
			Name:        "int64",
			Input:       int64(12),
			Expected:    "12",
			ShouldError: false,
		},
		{
			Name:     "int64 max",
			Input:    math.MaxInt64,
			Expected: "9223372036854775807",
		},
		{
			Name:     "int64 min",
			Input:    math.MinInt64,
			Expected: "-9223372036854775808",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			id, err := UnmarshalID(tt.Input)

			assert.Equal(t, tt.Expected, id)
			if tt.ShouldError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
