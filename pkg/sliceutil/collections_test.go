package sliceutil

import (
	"reflect"
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

func TestAppendUniques(t *testing.T) {
	type args struct {
		vs    []int
		toAdd []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "append to empty slice",
			args: args{
				vs:    []int{},
				toAdd: []int{1, 2, 3},
			},
			want: []int{1, 2, 3},
		},
		{
			name: "append all unique values",
			args: args{
				vs:    []int{1, 2, 3},
				toAdd: []int{4, 5, 6},
			},
			want: []int{1, 2, 3, 4, 5, 6},
		},
		{
			name: "append with some duplicates",
			args: args{
				vs:    []int{1, 2, 3},
				toAdd: []int{3, 4, 5},
			},
			want: []int{1, 2, 3, 4, 5},
		},
		{
			name: "append all duplicates",
			args: args{
				vs:    []int{1, 2, 3},
				toAdd: []int{1, 2, 3},
			},
			want: []int{1, 2, 3},
		},
		{
			name: "append to nil slice",
			args: args{
				vs:    nil,
				toAdd: []int{1, 2, 3},
			},
			want: []int{1, 2, 3},
		},
		{
			name: "append empty slice",
			args: args{
				vs:    []int{1, 2, 3},
				toAdd: []int{},
			},
			want: []int{1, 2, 3},
		},
		{
			name: "append nil to slice",
			args: args{
				vs:    []int{1, 2, 3},
				toAdd: nil,
			},
			want: []int{1, 2, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AppendUniques(tt.args.vs, tt.args.toAdd); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AppendUniques() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkAppendUniques(b *testing.B) {
	for i := 0; i < b.N; i++ {
		AppendUniques([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, []int{3, 4, 4, 11, 12, 13, 14, 15, 16, 17, 18})
	}
}
