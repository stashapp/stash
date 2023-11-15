package sliceutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSliceSame(t *testing.T) {
	tests := []struct {
		name string
		a    []int
		b    []int
		want bool
	}{
		{"nil values", nil, nil, true},
		{"empty", []int{}, []int{}, true},
		{"nil and empty", nil, []int{}, true},
		{
			"different length",
			[]int{1, 2, 3},
			[]int{1, 2},
			false,
		},
		{
			"equal",
			[]int{1, 2, 3, 4, 5},
			[]int{1, 2, 3, 4, 5},
			true,
		},
		{
			"different order",
			[]int{5, 4, 3, 2, 1},
			[]int{1, 2, 3, 4, 5},
			true,
		},
		{
			"different",
			[]int{5, 4, 3, 2, 6},
			[]int{1, 2, 3, 4, 5},
			false,
		},
		{
			"same with duplicates",
			[]int{1, 1, 2, 3, 4},
			[]int{1, 2, 3, 4, 1},
			true,
		},
		{
			"subset",
			[]int{1, 1, 2, 2, 3},
			[]int{1, 2, 3, 4, 5},
			false,
		},
		{
			"superset",
			[]int{1, 2, 3, 4, 5},
			[]int{1, 1, 2, 2, 3},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SliceSame(tt.a, tt.b)
			assert.Equal(t, tt.want, got)
		})
	}
}
