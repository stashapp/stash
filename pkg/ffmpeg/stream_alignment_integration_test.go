//go:build integration

package ffmpeg

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestAlignedHLSOutputIntegration(t *testing.T) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("ffmpeg is not installed")
	}
	if _, err := exec.LookPath("ffprobe"); err != nil {
		t.Skip("ffprobe is not installed")
	}

	tests := []struct {
		name  string
		size  string
		scale VideoFilter
		want  string
	}{
		{name: "small unaligned", size: "404x720", scale: "", want: "h264,416,720"},
		{name: "scaled unaligned", size: "1080x1920", scale: "scale=404:-2", want: "h264,416,720"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			input := filepath.Join(dir, "input.mp4")
			runIntegrationCommand(t, "ffmpeg",
				"-y", "-f", "lavfi", "-i", "testsrc2=size="+tt.size+":rate=30",
				"-t", "2", "-c:v", "libx264", "-pix_fmt", "yuv420p", input,
			)

			filter := videoFilterWithAlignment(tt.scale, StreamAlignment16x8)
			args := StreamTypeHLS.Args(VideoCodecLibX264, 0, filter, true, dir)
			commandArgs := append([]string{"-y", "-i", input}, args...)
			runIntegrationCommand(t, "ffmpeg", commandArgs...)

			output, err := exec.Command(
				"ffprobe", "-v", "error", "-select_streams", "v:0",
				"-show_entries", "stream=codec_name,width,height", "-of", "csv=p=0",
				filepath.Join(dir, "manifest.m3u8"),
			).CombinedOutput()
			if err != nil {
				t.Fatalf("ffprobe failed: %v\n%s", err, output)
			}
			reported := false
			for _, line := range strings.Split(strings.TrimSpace(string(output)), "\n") {
				got := strings.TrimSpace(line)
				if got == "" {
					continue
				}
				reported = true
				if got != tt.want {
					t.Fatalf("probed stream = %q, want %q", got, tt.want)
				}
			}
			if !reported {
				t.Fatal("ffprobe reported no video stream")
			}
		})
	}
}

func runIntegrationCommand(t *testing.T, name string, args ...string) {
	t.Helper()
	output, err := exec.Command(name, args...).CombinedOutput()
	if err != nil {
		t.Fatalf("%s %s failed: %v\n%s", name, fmt.Sprint(args), err, output)
	}
}
