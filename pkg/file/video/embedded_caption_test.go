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
	outputDir := filepath.Join(dir, "generated", "captions")
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

	captions, err := extractor.Extract(context.Background(), &models.VideoFile{
		BaseFile: &models.BaseFile{Path: videoPath},
	}, outputDir, "123", nil, false)

	require.NoError(t, err)
	assert.Equal(t, []*models.VideoCaption{
		{
			LanguageCode: "en",
			Filename:     "123.en.srt",
			CaptionType:  "srt",
			Generated:    true,
		},
	}, captions)
	require.Len(t, encoder.args, 1)
	assert.Equal(t, ffmpeg.Args{
		"-v", "error",
		"-nostdin",
		"-i", videoPath,
		"-map", "0:1",
		"-vn",
		"-an",
		"-c:s", "srt",
		filepath.Join(outputDir, "123.en.srt"),
	}, encoder.args[0])
}

func TestEmbeddedCaptionExtractorSkipsExistingExternalCaption(t *testing.T) {
	dir := t.TempDir()
	outputDir := filepath.Join(dir, "generated", "captions")
	videoPath := filepath.Join(dir, "video.mkv")

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

	captions, err := extractor.Extract(context.Background(), &models.VideoFile{
		BaseFile: &models.BaseFile{Path: videoPath},
	}, outputDir, "123", []*models.VideoCaption{
		{
			LanguageCode: "en",
			Filename:     "video.en.srt",
			CaptionType:  "srt",
		},
	}, false)

	require.NoError(t, err)
	assert.Empty(t, captions)
	assert.Empty(t, encoder.args)
}

func TestEmbeddedCaptionExtractorSkipsExistingGeneratedCaptionFile(t *testing.T) {
	dir := t.TempDir()
	outputDir := filepath.Join(dir, "generated", "captions")
	captionPath := filepath.Join(outputDir, "123.en.srt")
	require.NoError(t, os.MkdirAll(outputDir, 0755))
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

	captions, err := extractor.Extract(context.Background(), &models.VideoFile{
		BaseFile: &models.BaseFile{Path: filepath.Join(dir, "video.mkv")},
	}, outputDir, "123", []*models.VideoCaption{
		{
			LanguageCode: "en",
			Filename:     "123.en.srt",
			CaptionType:  "srt",
			Generated:    true,
		},
	}, false)

	require.NoError(t, err)
	assert.Equal(t, []*models.VideoCaption{
		{
			LanguageCode: "en",
			Filename:     "123.en.srt",
			CaptionType:  "srt",
			Generated:    true,
		},
	}, captions)
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
	}, t.TempDir(), "123", nil, false)

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
	outputDir := filepath.Join(dir, "generated", "captions")
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

	captions, err := extractor.Extract(context.Background(), &models.VideoFile{
		BaseFile: &models.BaseFile{Path: videoPath},
	}, outputDir, "123", nil, false)

	captionPath := filepath.Join(outputDir, "123.en.srt")
	require.NoError(t, err)
	assert.Equal(t, []*models.VideoCaption{
		{
			LanguageCode: "en",
			Filename:     "123.en.srt",
			CaptionType:  "srt",
			Generated:    true,
		},
	}, captions)

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
