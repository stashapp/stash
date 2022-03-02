package sliceutil

import "testing"

func TestSliceSame(t *testing.T) {
	objs := []struct {
		a string
		b int
	}{
		{"1", 2},
		{"1", 2},
		{"2", 1},
	}

	tests := []struct {
		name string
		a    interface{}
		b    interface{}
		want bool
	}{
		{"nil values", nil, nil, true},
		{"empty", []int{}, []int{}, true},
		{"nil and empty", nil, []int{}, true},
		{
			"different type",
			[]string{"1"},
			[]int{1},
			false,
		},
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
		{
			"structs equal",
			objs[0:1],
			objs[0:1],
			true,
		},
		{
			"structs not equal",
			objs[0:2],
			objs[1:3],
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SliceSame(tt.a, tt.b); got != tt.want {
				t.Errorf("SliceSame() = %v, want %v", got, tt.want)
			}
		})
	}
}
