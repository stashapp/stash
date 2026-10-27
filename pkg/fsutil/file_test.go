package fsutil

import (
	"runtime"
	"testing"
)

// Unicode forms of the Japanese character "が" used to exercise NFC/NFD handling.
const (
	gaNFC = "が"  // composed: が
	gaNFD = "が" // decomposed: か + combining dakuten
)

func TestNormalizePath(t *testing.T) {
	// NFD input is normalized to NFC on macOS, left unchanged elsewhere.
	wantNFD := gaNFD
	if runtime.GOOS == "darwin" {
		wantNFD = gaNFC
	}

	tests := []struct {
		name string
		in   string
		want string
	}{
		{"nfd", gaNFD, wantNFD},
		{"nfc unchanged", gaNFC, gaNFC},
		{"ascii unchanged", "/a/b/c.mp4", "/a/b/c.mp4"},
		{"empty", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NormalizePath(tt.in); got != tt.want {
				t.Errorf("NormalizePath(%q) = %q, want %q (GOOS=%s)", tt.in, got, tt.want, runtime.GOOS)
			}
		})
	}
}

func TestNormalizePaths(t *testing.T) {
	in := []string{gaNFD, "/ascii"}
	got := NormalizePaths(in)

	// must return a new slice, leaving the input untouched
	if in[0] != gaNFD {
		t.Errorf("NormalizePaths mutated its input: %q", in[0])
	}

	want0 := gaNFD
	if runtime.GOOS == "darwin" {
		want0 = gaNFC
	}
	if got[0] != want0 || got[1] != "/ascii" {
		t.Errorf("NormalizePaths() = %q, want [%q /ascii]", got, want0)
	}
}

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
		{"unicode cjk", `テスト`, "テスト-63b560db"},
		{"unicode korean", `시험`, "시험-3fcc7beb"},
		{"mixed unicode", `Test テスト`, "Test-テスト-366aff1e"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SanitiseBasename(tt.v); got != tt.want {
				t.Errorf("SanitiseBasename() = %v, want %v", got, tt.want)
			}
		})
	}
}
