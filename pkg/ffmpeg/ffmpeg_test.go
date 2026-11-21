// Package ffmpeg provides a wrapper around the ffmpeg and ffprobe executables.
package ffmpeg

import "testing"

func TestParseEncoders(t *testing.T) {
	out := `Encoders:
 V..... = Video
 ------
 V....D libx264              libx264 H.264 / AVC / MPEG-4 AVC / MPEG-4 part 10 (codec h264)
 V....D libsvtav1            SVT-AV1(Scalable Video Technology for AV1) encoder (codec av1)
 A....D aac                  AAC (Advanced Audio Coding)
`
	got := parseEncoders(out)
	for _, name := range []string{"libx264", "libsvtav1", "aac"} {
		if !got[name] {
			t.Errorf("missing encoder %s", name)
		}
	}
	if got["libx265"] {
		t.Error("libx265 should not be detected")
	}
}

func TestFFMpegVersion_GreaterThan(t *testing.T) {
	tests := []struct {
		name  string
		this  Version
		other Version
		want  bool
	}{
		{
			"major greater, minor equal, patch equal",
			Version{2, 0, 0},
			Version{1, 0, 0},
			true,
		},
		{
			"major greater, minor less, patch less",
			Version{2, 1, 1},
			Version{1, 0, 0},
			true,
		},
		{
			"major equal, minor greater, patch equal",
			Version{1, 1, 0},
			Version{1, 0, 0},
			true,
		},
		{
			"major equal, minor equal, patch greater",
			Version{1, 0, 1},
			Version{1, 0, 0},
			true,
		},
		{
			"major equal, minor equal, patch equal",
			Version{1, 0, 0},
			Version{1, 0, 0},
			true,
		},
		{
			"major less, minor equal, patch equal",
			Version{1, 0, 0},
			Version{2, 0, 0},
			false,
		},
		{
			"major equal, minor less, patch equal",
			Version{1, 0, 0},
			Version{1, 1, 0},
			false,
		},
		{
			"major equal, minor equal, patch less",
			Version{1, 0, 0},
			Version{1, 0, 1},
			false,
		},
		{
			"major less, minor less, patch less",
			Version{1, 0, 0},
			Version{2, 1, 1},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.this.Gteq(tt.other); got != tt.want {
				t.Errorf("FFMpegVersion.GreaterThan() = %v, want %v", got, tt.want)
			}
		})
	}
}
