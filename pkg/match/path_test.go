package match

import "testing"

func Test_nameMatchesPath(t *testing.T) {
	const name = "first last"

	tests := []struct {
		name string
		path string
		want int
	}{
		{
			"exact",
			name,
			0,
		},
		{
			"partial",
			"first",
			-1,
		},
		{
			"separator",
			"first.last",
			0,
		},
		{
			"separator",
			"first-last",
			0,
		},
		{
			"separator",
			"first_last",
			0,
		},
		{
			"separators",
			"first.-_ last",
			0,
		},
		{
			"within string",
			"before_first last/after",
			6,
		},
		{
			"not within string",
			"beforefirst last/after",
			-1,
		},
		{
			"not within string",
			"before/first lastafter",
			-1,
		},
		{
			"not within string",
			"first last1",
			-1,
		},
		{
			"not within string",
			"1first last",
			-1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := nameMatchesPath(name, tt.path); got != tt.want {
				t.Errorf("nameMatchesPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
