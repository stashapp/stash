package video

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testCase struct {
	captionPath    string
	expectedLang   string
	expectedResult string
}

var testCases = []testCase{
	{
		captionPath:    "/stash/video.vtt",
		expectedLang:   LangUnknown,
		expectedResult: "/stash/video.",
	},
	{
		captionPath:    "/stash/video.en.vtt",
		expectedLang:   "en",
		expectedResult: "/stash/video.", // lang code valid, remove en part
	},
	{
		captionPath:    "/stash/video.test.srt",
		expectedLang:   LangUnknown,
		expectedResult: "/stash/video.test.", // no lang code/lang code invalid test should remain
	},
	{
		captionPath:    "C:\\videos\\video.fr.srt",
		expectedLang:   "fr",
		expectedResult: "C:\\videos\\video.",
	},
	{
		captionPath:    "C:\\videos\\video.xx.srt",
		expectedLang:   LangUnknown,
		expectedResult: "C:\\videos\\video.xx.", // no lang code/lang code invalid xx should remain
	},
}

func TestGenerateCaptionCandidates(t *testing.T) {
	for _, c := range testCases {
		assert.Equal(t, c.expectedResult, getCaptionPrefix(c.captionPath))
	}
}

func TestGetCaptionsLangFromPath(t *testing.T) {
	for _, l := range testCases {
		assert.Equal(t, l.expectedLang, getCaptionsLangFromPath(l.captionPath))
	}
}

func TestEmbeddedCaptionCandidates(t *testing.T) {
	streams := []models.VideoSubtitleStream{
		{Index: 2, LanguageCode: "ENG"},
		{Index: 3, LanguageCode: "fr"},
		{Index: 4, LanguageCode: ""},
		{Index: 5, LanguageCode: "eng"},
	}
	captions := []*models.VideoCaption{
		{LanguageCode: "fr", CaptionType: embeddedCaptionType},
	}
	existingFiles := map[string]bool{
		"/stash/video.en.vtt": true,
	}

	candidates := embeddedCaptionCandidates("/stash/video.mkv", streams, captions, func(path string) bool {
		return existingFiles[path]
	})

	assert.Len(t, candidates, 2)
	assert.Equal(t, 2, candidates[0].streamIndex)
	assert.Equal(t, "/stash/video.en.vtt", candidates[0].captionPath)
	assert.Equal(t, "video.en.vtt", candidates[0].caption.Filename)
	assert.Equal(t, "en", candidates[0].caption.LanguageCode)
	assert.False(t, candidates[0].needsExtraction)

	assert.Equal(t, 4, candidates[1].streamIndex)
	assert.Equal(t, "/stash/video.vtt", candidates[1].captionPath)
	assert.Equal(t, "video.vtt", candidates[1].caption.Filename)
	assert.Equal(t, LangUnknown, candidates[1].caption.LanguageCode)
	assert.True(t, candidates[1].needsExtraction)
}

type fakeEmbeddedCaptionEncoder struct {
	args ffmpeg.Args
}

func (e *fakeEmbeddedCaptionEncoder) Generate(ctx context.Context, args ffmpeg.Args) error {
	e.args = args
	return os.WriteFile(args[len(args)-1], []byte("WEBVTT\n\n00:00:00.000 --> 00:00:01.000\nTest\n"), 0600)
}

func TestExtractEmbeddedCaption(t *testing.T) {
	dir := t.TempDir()
	videoPath := filepath.Join(dir, "video.mkv")
	captionPath := filepath.Join(dir, "video.en.vtt")
	encoder := &fakeEmbeddedCaptionEncoder{}

	err := extractEmbeddedCaption(context.Background(), encoder, videoPath, 2, captionPath)

	require.NoError(t, err)
	got, err := os.ReadFile(captionPath)
	require.NoError(t, err)
	assert.Contains(t, string(got), "WEBVTT")
	assert.Equal(t, ffmpeg.Args{
		"-v", "error",
		"-y",
		"-i", videoPath,
		"-map", "0:2",
		"-c:s", "webvtt",
		"-f", "webvtt",
		captionPath + ".2.tmp",
	}, encoder.args)
}

func TestExtractEmbeddedCaptionDoesNotOverwriteExistingCaption(t *testing.T) {
	dir := t.TempDir()
	videoPath := filepath.Join(dir, "video.mkv")
	captionPath := filepath.Join(dir, "video.en.vtt")
	require.NoError(t, os.WriteFile(captionPath, []byte("existing"), 0600))

	err := extractEmbeddedCaption(context.Background(), &fakeEmbeddedCaptionEncoder{}, videoPath, 2, captionPath)

	require.NoError(t, err)
	got, err := os.ReadFile(captionPath)
	require.NoError(t, err)
	assert.Equal(t, "existing", string(got))
}
