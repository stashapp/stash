// Package ffmpeg provides a wrapper around the ffmpeg and ffprobe executables.
package ffmpeg

import "testing"

func TestFFMpegVersion_GreaterThan(t *testing.T) {
	tests := []struct {
		name  string
		this  FFMpegVersion
		other FFMpegVersion
		want  bool
	}{
		{
			"major greater, minor equal, patch equal",
			FFMpegVersion{2, 0, 0},
			FFMpegVersion{1, 0, 0},
			true,
		},
		{
			"major greater, minor less, patch less",
			FFMpegVersion{2, 1, 1},
			FFMpegVersion{1, 0, 0},
			true,
		},
		{
			"major equal, minor greater, patch equal",
			FFMpegVersion{1, 1, 0},
			FFMpegVersion{1, 0, 0},
			true,
		},
		{
			"major equal, minor equal, patch greater",
			FFMpegVersion{1, 0, 1},
			FFMpegVersion{1, 0, 0},
			true,
		},
		{
			"major equal, minor equal, patch equal",
			FFMpegVersion{1, 0, 0},
			FFMpegVersion{1, 0, 0},
			true,
		},
		{
			"major less, minor equal, patch equal",
			FFMpegVersion{1, 0, 0},
			FFMpegVersion{2, 0, 0},
			false,
		},
		{
			"major equal, minor less, patch equal",
			FFMpegVersion{1, 0, 0},
			FFMpegVersion{1, 1, 0},
			false,
		},
		{
			"major equal, minor equal, patch less",
			FFMpegVersion{1, 0, 0},
			FFMpegVersion{1, 0, 1},
			false,
		},
		{
			"major less, minor less, patch less",
			FFMpegVersion{1, 0, 0},
			FFMpegVersion{2, 1, 1},
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
