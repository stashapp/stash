package scanner

import "testing"

func TestIsSupportedVideoPath(t *testing.T) {
	tests := []struct {
		path string
		want bool
	}{
		{"scene.mp4", true},
		{"SCENE.MKV", true},
		{"clip.webm", true},
		{"cover.jpg", false},
		{"notes.txt", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			if got := isSupportedVideoPath(tt.path); got != tt.want {
				t.Fatalf("isSupportedVideoPath(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}
