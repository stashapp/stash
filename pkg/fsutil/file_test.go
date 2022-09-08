package fsutil

import "testing"

func TestSanitiseBasename(t *testing.T) {
	tests := []struct {
		name string
		v    string
		want string
	}{
		{"basic", "basic", "basic"},
		{"spaces", `spaced name`, "spaced-name"},
		{"leading/trailing spaces", `  spaced name  `, "spaced-name"},
		{"hyphen name", `hyphened-name`, "hyphened-name"},
		{"multi-hyphen", `hyphened--name`, "hyphened-name"},
		{"replaced characters", `a&b=c\d/:e*"f?_ g`, "a-b-c-d-e-f-g"},
		{"removed characters", `foo!!bar@@and, more`, "foobarand-more"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SanitiseBasename(tt.v); got != tt.want {
				t.Errorf("SanitiseBasename() = %v, want %v", got, tt.want)
			}
		})
	}
}
