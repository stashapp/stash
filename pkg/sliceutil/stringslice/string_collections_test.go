package stringslice

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUniqueExcludeFold(t *testing.T) {
	tests := []struct {
		name    string
		values  []string
		exclude string
		want    []string
	}{
		{"nil", nil, "name", []string{}},
		{"no duplicates", []string{"a", "b"}, "name", []string{"a", "b"}},
		{"duplicates removed case-insensitively", []string{"a", "A", "b"}, "name", []string{"a", "b"}},
		{"excluded value removed case-insensitively", []string{"a", "NAME"}, "name", []string{"a"}},
		{"first occurrence kept", []string{"B", "b", "a"}, "name", []string{"B", "a"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, UniqueExcludeFold(tt.values, tt.exclude))
		})
	}
}
