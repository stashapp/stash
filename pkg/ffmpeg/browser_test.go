package ffmpeg

import "testing"

func TestIsStreamable(t *testing.T) {
	tests := []struct {
		name      string
		codec     string
		audio     ProbeAudioCodec
		container Container
		supported []string
		want      bool
	}{
		{"h264 mp4 default", H264, Aac, Mp4, DefaultSupportedCodecs, true},
		{"hevc mp4 default", Hevc, Aac, Mp4, DefaultSupportedCodecs, false},
		{"h265 mp4 default", H265, Aac, Mp4, DefaultSupportedCodecs, false},
		{"hevc mp4 whitelisted", Hevc, Aac, Mp4, []string{H264, Hevc}, true},
		{"hevc mkv", Hevc, Aac, Matroska, []string{Hevc}, false},
		{"h264 mkv", H264, Aac, Matroska, DefaultSupportedCodecs, false},
		{"av1 mkv even with mkv listed", "av1", Opus, Matroska, []string{"av1", Mkv}, false},
		{"av1 webm whitelisted", "av1", Opus, Webm, []string{"av1"}, true},
		{"av1 avi whitelisted", "av1", Mp3, Avi, []string{"av1"}, false},
		{"av1 webm bad audio", "av1", Aac, Webm, []string{"av1"}, false},
		{"av1 not whitelisted", "av1", Opus, Webm, DefaultSupportedCodecs, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := IsStreamable(tt.codec, tt.audio, tt.container, tt.supported)
			if got := err == nil; got != tt.want {
				t.Errorf("IsStreamable() = %v, want %v", err, tt.want)
			}
		})
	}
}
