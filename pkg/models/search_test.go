package models

import (
	"reflect"
	"testing"
)

func TestParseSearchString(t *testing.T) {
	tests := []struct {
		name string
		q    string
		want SearchSpecs
	}{
		{
			"basic",
			"a b c",
			SearchSpecs{
				MustHave: []string{"a", "b", "c"},
			},
		},
		{
			"empty",
			"",
			SearchSpecs{},
		},
		{
			"whitespace",
			" ",
			SearchSpecs{},
		},
		{
			"single",
			"a",
			SearchSpecs{
				MustHave: []string{"a"},
			},
		},
		{
			"quoted",
			`"a b" c`,
			SearchSpecs{
				MustHave: []string{"a b", "c"},
			},
		},
		{
			"quoted double space",
			`"a  b" c`,
			SearchSpecs{
				MustHave: []string{"a  b", "c"},
			},
		},
		{
			"quoted end space",
			`"a  b " c`,
			SearchSpecs{
				MustHave: []string{"a  b ", "c"},
			},
		},
		{
			"no matching end quote",
			`"a b c`,
			SearchSpecs{
				MustHave: []string{`"a`, "b", "c"},
			},
		},
		{
			"no matching start quote",
			`a b c"`,
			SearchSpecs{
				MustHave: []string{"a", "b", `c"`},
			},
		},
		{
			"or",
			"a OR b",
			SearchSpecs{
				AnySets: [][]string{
					{"a", "b"},
				},
			},
		},
		{
			"multi or",
			"a OR b c OR d",
			SearchSpecs{
				AnySets: [][]string{
					{"a", "b"},
					{"c", "d"},
				},
			},
		},
		{
			"lowercase or",
			"a or b",
			SearchSpecs{
				AnySets: [][]string{
					{"a", "b"},
				},
			},
		},
		{
			"or symbol",
			"a | b",
			SearchSpecs{
				AnySets: [][]string{
					{"a", "b"},
				},
			},
		},
		{
			"quoted or",
			`a "OR" b`,
			SearchSpecs{
				MustHave: []string{"a", "OR", "b"},
			},
		},
		{
			"quoted or symbol",
			`a "|" b`,
			SearchSpecs{
				MustHave: []string{"a", "|", "b"},
			},
		},
		{
			"or phrases",
			`"a b" OR "c d"`,
			SearchSpecs{
				AnySets: [][]string{
					{"a b", "c d"},
				},
			},
		},
		{
			"or at start",
			"OR a",
			SearchSpecs{
				MustHave: []string{"OR", "a"},
			},
		},
		{
			"or at end",
			"a OR",
			SearchSpecs{
				MustHave: []string{"a", "OR"},
			},
		},
		{
			"or symbol at start",
			"| a",
			SearchSpecs{
				MustHave: []string{"|", "a"},
			},
		},
		{
			"or symbol at end",
			"a |",
			SearchSpecs{
				MustHave: []string{"a", "|"},
			},
		},
		{
			"nots",
			"-a -b",
			SearchSpecs{
				MustNot: []string{"a", "b"},
			},
		},
		{
			"not or",
			"-a OR b",
			SearchSpecs{
				AnySets: [][]string{
					{"-a", "b"},
				},
			},
		},
		{
			"not phrase",
			`-"a b"`,
			SearchSpecs{
				MustNot: []string{"a b"},
			},
		},
		{
			"not in phrase",
			`"-a b"`,
			SearchSpecs{
				MustHave: []string{"-a b"},
			},
		},
		{
			"double not",
			"--a",
			SearchSpecs{
				MustNot: []string{"-a"},
			},
		},
		{
			"empty quote",
			`"" a`,
			SearchSpecs{
				MustHave: []string{"a"},
			},
		},
		{
			"not empty quote",
			`-"" a`,
			SearchSpecs{
				MustHave: []string{"a"},
			},
		},
		{
			"quote in word",
			`ab"cd"`,
			SearchSpecs{
				MustHave: []string{`ab"cd"`},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseSearchString(tt.q); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindFilterType.ParseSearchString() = %v, want %v", got, tt.want)
			}
		})
	}
}
