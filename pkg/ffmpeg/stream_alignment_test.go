package ffmpeg

import (
	"context"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stashapp/stash/pkg/models"
)

func TestValidHLSAlignment(t *testing.T) {
	tests := []struct {
		name      string
		stream    *StreamType
		alignment string
		want      bool
	}{
		{name: "empty HLS", stream: StreamTypeHLS, alignment: "", want: true},
		{name: "aligned HLS", stream: StreamTypeHLS, alignment: StreamAlignment16x8, want: true},
		{name: "invalid HLS", stream: StreamTypeHLS, alignment: "2x2", want: false},
		{name: "alignment on DASH", stream: StreamTypeDASHVideo, alignment: StreamAlignment16x8, want: false},
		{name: "empty DASH", stream: StreamTypeDASHVideo, alignment: "", want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidHLSAlignment(tt.stream, tt.alignment); got != tt.want {
				t.Fatalf("ValidHLSAlignment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStreamDirWithAlignment(t *testing.T) {
	const base = "scene_hls_720"
	if got := streamDirWithAlignment(base, ""); got != base {
		t.Fatalf("ordinary cache dir = %q", got)
	}
	if got := streamDirWithAlignment(base, StreamAlignment16x8); got != base+"_align16x8" {
		t.Fatalf("aligned cache dir = %q", got)
	}
}

func TestVideoFilterWithAlignment(t *testing.T) {
	ordinary := VideoFilter("scale=-2:720")
	if got := videoFilterWithAlignment(ordinary, ""); got != ordinary {
		t.Fatalf("ordinary filter = %q", got)
	}
	want := "scale=-2:720,pad=ceil(iw/16)*16:ceil(ih/8)*8:x=0:y=0:color=black"
	if got := videoFilterWithAlignment(ordinary, StreamAlignment16x8); string(got) != want {
		t.Fatalf("aligned filter = %q, want %q", got, want)
	}
}

func TestAlignedHLSForcesCPUAndDisablesFullHardware(t *testing.T) {
	codec := hlsCodecWithAlignment(nil, "hls", StreamAlignment16x8)
	if codec.CodeName != VideoCodecLibX264.CodeName {
		t.Fatalf("codec = %q, want libx264", codec.CodeName)
	}
	if fullHardwareWithAlignment(true, StreamAlignment16x8) {
		t.Fatal("aligned HLS must not use full hardware filters")
	}
	if !fullHardwareWithAlignment(true, "") {
		t.Fatal("ordinary HLS must preserve the hardware setting")
	}
}

func TestHLSSegmentQuery(t *testing.T) {
	req := httptest.NewRequest(
		http.MethodGet,
		"http://stash.example/scene/7/stream.m3u8?resolution=STANDARD&align=16x8&apikey=test-only",
		nil,
	)
	got := hlsSegmentQuery(req.URL, "STANDARD")
	want := "?align=16x8&apikey=test-only&resolution=STANDARD"
	if got != want {
		t.Fatalf("segment query = %q, want %q", got, want)
	}
}

func TestHLSSegmentQueryWithoutAlignment(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://stash.example/scene/7/stream.m3u8", nil)
	if got := hlsSegmentQuery(req.URL, ""); got != "" {
		t.Fatalf("ordinary segment query = %q, want empty", got)
	}
}

func TestHLSSegmentQueryAuthentication(t *testing.T) {
	tests := []struct {
		name   string
		target string
		want   string
	}{
		{
			name:   "signed URL",
			target: "http://stash.example/scene/7/stream.m3u8?align=16x8&cid=test-credential&expires=2000000000&signature=test-signature",
			want:   "?align=16x8&cid=test-credential&expires=2000000000&signature=test-signature",
		},
		{
			name:   "API key precedence",
			target: "http://stash.example/scene/7/stream.m3u8?align=16x8&apikey=test-only&cid=test-credential&expires=2000000000&signature=test-signature",
			want:   "?align=16x8&apikey=test-only",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.target, nil)
			if got := hlsSegmentQuery(req.URL, ""); got != tt.want {
				t.Fatalf("segment query = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestAlignedHLSMakeStreamArgs(t *testing.T) {
	sm := &StreamManager{
		encoder: &FFMpeg{hwCodecSupport: []VideoCodec{VideoCodecN264}},
		config:  alignmentStreamConfig{},
	}
	stream := &runningStream{
		streamType: StreamTypeHLS,
		vf: &models.VideoFile{
			BaseFile: &models.BaseFile{Path: "input.mp4"},
			Width:    404,
			Height:   720,
		},
		alignment: StreamAlignment16x8,
		outputDir: "output",
	}

	args := strings.Join(stream.makeStreamArgs(sm, 0), " ")
	if !strings.Contains(args, "-c:v libx264") {
		t.Fatalf("transcode args do not select libx264: %s", args)
	}
	wantFilter := "pad=ceil(iw/16)*16:ceil(ih/8)*8:x=0:y=0:color=black"
	if !strings.Contains(args, wantFilter) {
		t.Fatalf("transcode args do not contain alignment filter %q: %s", wantFilter, args)
	}
}

func TestServeSegmentUsesAlignedCache(t *testing.T) {
	sm := &StreamManager{
		cacheDir:       t.TempDir(),
		config:         alignmentStreamConfig{},
		runningStreams: make(map[string]*runningStream),
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	req := httptest.NewRequest(http.MethodGet, "/scene/7/0.ts?align=16x8", nil).WithContext(ctx)
	options := StreamOptions{
		StreamType: StreamTypeHLS,
		VideoFile: &models.VideoFile{
			BaseFile: &models.BaseFile{Path: "input.mp4"},
			Duration: 10,
		},
		Alignment: StreamAlignment16x8,
		Hash:      "example",
		Segment:   "0",
	}

	sm.ServeSegment(httptest.NewRecorder(), req, options)

	const dir = "example_hls_align16x8"
	stream := sm.runningStreams[dir]
	if stream == nil {
		t.Fatalf("aligned stream was not registered under %q", dir)
	}
	if stream.alignment != StreamAlignment16x8 {
		t.Fatalf("stream alignment = %q, want %q", stream.alignment, StreamAlignment16x8)
	}
	if want := filepath.Join(sm.cacheDir, dir); stream.outputDir != want {
		t.Fatalf("stream output dir = %q, want %q", stream.outputDir, want)
	}
}

type alignmentStreamConfig struct{}

func (alignmentStreamConfig) GetMaxStreamingTranscodeSize() models.StreamingResolutionEnum {
	return models.StreamingResolutionEnumOriginal
}

func (alignmentStreamConfig) GetLiveTranscodeInputArgs() []string {
	return nil
}

func (alignmentStreamConfig) GetLiveTranscodeOutputArgs() []string {
	return nil
}

func (alignmentStreamConfig) GetTranscodeHardwareAcceleration() bool {
	return true
}
