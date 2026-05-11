package video

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type stubEmbeddedCaptionProbe struct {
	file *ffmpeg.VideoFile
	err  error
}

func (p stubEmbeddedCaptionProbe) NewVideoFile(string) (*ffmpeg.VideoFile, error) {
	return p.file, p.err
}

type stubEmbeddedCaptionEncoder struct {
	args []ffmpeg.Args
	err  error
}

func (e *stubEmbeddedCaptionEncoder) Generate(_ context.Context, args ffmpeg.Args) error {
	e.args = append(e.args, args)
	return e.err
}

func TestEmbeddedCaptionTracks(t *testing.T) {
	streams := []ffmpeg.FFProbeStream{
		subtitleStream(1, "subrip", "eng"),
		subtitleStream(2, "ass", "eng"),
		subtitleStream(3, "hdmv_pgs_subtitle", "eng"),
		subtitleStream(4, "webvtt", "und"),
		subtitleStream(5, "mov_text", "fr"),
		{
			Index:     6,
			CodecType: "audio",
			CodecName: "subrip",
		},
	}

	tracks := embeddedCaptionTracks(streams)

	require.Len(t, tracks, 3)
	assert.Equal(t, embeddedCaptionTrack{streamIndex: 1, language: "en", captionType: "srt"}, tracks[0])
	assert.Equal(t, embeddedCaptionTrack{streamIndex: 4, language: LangUnknown, captionType: "srt"}, tracks[1])
	assert.Equal(t, embeddedCaptionTrack{streamIndex: 5, language: "fr", captionType: "srt"}, tracks[2])
}

func TestEmbeddedCaptionExtractorExtractsSupportedTracks(t *testing.T) {
	dir := t.TempDir()
	videoPath := filepath.Join(dir, "video.mkv")
	probe := stubEmbeddedCaptionProbe{
		file: &ffmpeg.VideoFile{
			JSON: ffmpeg.FFProbeJSON{
				Streams: []ffmpeg.FFProbeStream{
					subtitleStream(1, "subrip", "eng"),
					subtitleStream(2, "hdmv_pgs_subtitle", "eng"),
				},
			},
		},
	}
	encoder := &stubEmbeddedCaptionEncoder{}
	extractor := EmbeddedCaptionExtractor{
		FFProbe: probe,
		FFMpeg:  encoder,
	}

	paths, err := extractor.Extract(context.Background(), &models.VideoFile{
		BaseFile: &models.BaseFile{Path: videoPath},
	})

	require.NoError(t, err)
	assert.Equal(t, []string{filepath.Join(dir, "video.en.srt")}, paths)
	require.Len(t, encoder.args, 1)
	assert.Equal(t, ffmpeg.Args{
		"-v", "error",
		"-nostdin",
		"-i", videoPath,
		"-map", "0:1",
		"-vn",
		"-an",
		"-c:s", "srt",
		filepath.Join(dir, "video.en.srt"),
	}, encoder.args[0])
}

func TestEmbeddedCaptionExtractorSkipsExistingSidecar(t *testing.T) {
	dir := t.TempDir()
	videoPath := filepath.Join(dir, "video.mkv")
	captionPath := filepath.Join(dir, "video.en.srt")
	require.NoError(t, os.WriteFile(captionPath, []byte("synthetic caption"), 0o600))

	encoder := &stubEmbeddedCaptionEncoder{}
	extractor := EmbeddedCaptionExtractor{
		FFProbe: stubEmbeddedCaptionProbe{
			file: &ffmpeg.VideoFile{
				JSON: ffmpeg.FFProbeJSON{
					Streams: []ffmpeg.FFProbeStream{subtitleStream(1, "subrip", "eng")},
				},
			},
		},
		FFMpeg: encoder,
	}

	paths, err := extractor.Extract(context.Background(), &models.VideoFile{
		BaseFile: &models.BaseFile{Path: videoPath},
	})

	require.NoError(t, err)
	assert.Equal(t, []string{captionPath}, paths)
	assert.Empty(t, encoder.args)
}

func TestEmbeddedCaptionExtractorReturnsEncoderErrors(t *testing.T) {
	expectedErr := errors.New("ffmpeg failed")
	extractor := EmbeddedCaptionExtractor{
		FFProbe: stubEmbeddedCaptionProbe{
			file: &ffmpeg.VideoFile{
				JSON: ffmpeg.FFProbeJSON{
					Streams: []ffmpeg.FFProbeStream{subtitleStream(1, "subrip", "eng")},
				},
			},
		},
		FFMpeg: &stubEmbeddedCaptionEncoder{err: expectedErr},
	}

	_, err := extractor.Extract(context.Background(), &models.VideoFile{
		BaseFile: &models.BaseFile{Path: filepath.Join(t.TempDir(), "video.mkv")},
	})

	require.Error(t, err)
	assert.ErrorIs(t, err, expectedErr)
}

func TestEmbeddedCaptionExtractorWithSyntheticMedia(t *testing.T) {
	if os.Getenv("STASH_FFMPEG_INTEGRATION") != "1" {
		t.Skip("set STASH_FFMPEG_INTEGRATION=1 to run the synthetic ffmpeg fixture test")
	}

	ffmpegPath, err := exec.LookPath("ffmpeg")
	require.NoError(t, err)
	ffprobePath, err := exec.LookPath("ffprobe")
	require.NoError(t, err)

	dir := t.TempDir()
	subtitlePath := filepath.Join(dir, "input.srt")
	require.NoError(t, os.WriteFile(subtitlePath, []byte("1\n00:00:00,000 --> 00:00:01,000\nSynthetic embedded caption\n"), 0o600))

	videoPath := filepath.Join(dir, "video.mkv")
	cmd := exec.Command(
		ffmpegPath,
		"-y",
		"-v", "error",
		"-f", "lavfi",
		"-i", "color=c=black:s=16x16:d=1",
		"-i", subtitlePath,
		"-map", "0:v:0",
		"-map", "1:0",
		"-metadata:s:s:0", "language=eng",
		"-c:v", "ffv1",
		"-c:s", "srt",
		videoPath,
	)
	require.NoError(t, cmd.Run())

	extractor := EmbeddedCaptionExtractor{
		FFProbe: ffmpeg.NewFFProbe(ffprobePath),
		FFMpeg:  ffmpeg.NewEncoder(ffmpegPath),
	}

	paths, err := extractor.Extract(context.Background(), &models.VideoFile{
		BaseFile: &models.BaseFile{Path: videoPath},
	})

	captionPath := filepath.Join(dir, "video.en.srt")
	require.NoError(t, err)
	assert.Equal(t, []string{captionPath}, paths)

	extracted, err := os.ReadFile(captionPath)
	require.NoError(t, err)
	assert.Contains(t, string(extracted), "Synthetic embedded caption")
}

func subtitleStream(index int, codecName string, lang string) ffmpeg.FFProbeStream {
	ret := ffmpeg.FFProbeStream{
		Index:     index,
		CodecType: "subtitle",
		CodecName: codecName,
	}
	ret.Tags.Language = lang
	return ret
}
