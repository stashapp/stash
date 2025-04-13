package fsutil

import "testing"

func TestSanitiseBasename(t *testing.T) {
	tests := []struct {
		name string
		v    string
		want string
	}{
		{"basic", "basic", "basic-61a7508e"},
		{"spaces", `spaced name`, "spaced-name-b297cf60"},
		{"leading/trailing spaces", `  spaced name  `, "spaced-name-175433e9"},
		{"hyphen name", `hyphened-name`, "hyphened-name-789c55f2"},
		{"multi-hyphen", `hyphened--name`, "hyphened-name-2da2a58f"},
		{"replaced characters", `a&b=c\d/:e*"f?_ g`, "a-b-c-d-e-f-g-ffca6fb0"},
		{"removed characters", `foo!!bar@@and, more`, "foobarand-more-7cee02ab"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SanitiseBasename(tt.v); got != tt.want {
				t.Errorf("SanitiseBasename() = %v, want %v", got, tt.want)
			}
		})
	}
}
