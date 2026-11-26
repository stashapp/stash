package migrations

import (
	"reflect"
	"testing"
)

func TestSplitGroupAliases(t *testing.T) {
	tests := []struct {
		in   string
		want []string
	}{
		{"abc", []string{"abc"}},
		{" abc ", []string{"abc"}},
		{"a, b,c", []string{"a", "b", "c"}},
		{"a, , b,", []string{"a", "b"}},
		{"a, a, b", []string{"a", "b"}},
		{" , ", nil},
	}

	for _, tt := range tests {
		if got := splitGroupAliases(tt.in); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("splitGroupAliases(%q) = %v, want %v", tt.in, got, tt.want)
		}
	}
}
