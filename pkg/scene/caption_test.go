package scene

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testExts = []string{"mkv", "mp4"}

type testCase struct {
	captionPath        string
	expectedLang       string
	expectedCandidates []string
}

var testCases = []testCase{
	{
		captionPath:        "/stash/video.vtt",
		expectedLang:       LangUnknown,
		expectedCandidates: []string{"/stash/video.mkv", "/stash/video.mp4"},
	},
	{
		captionPath:        "/stash/video.en.vtt",
		expectedLang:       "en",
		expectedCandidates: []string{"/stash/video.mkv", "/stash/video.mp4"}, // lang code valid, remove en part
	},
	{
		captionPath:        "/stash/video.test.srt",
		expectedLang:       LangUnknown,
		expectedCandidates: []string{"/stash/video.test.mkv", "/stash/video.test.mp4"}, // no lang code/lang code invalid test should remain
	},
	{
		captionPath:        "C:\\videos\\video.fr.srt",
		expectedLang:       "fr",
		expectedCandidates: []string{"C:\\videos\\video.mkv", "C:\\videos\\video.mp4"},
	},
	{
		captionPath:        "C:\\videos\\video.xx.srt",
		expectedLang:       LangUnknown,
		expectedCandidates: []string{"C:\\videos\\video.xx.mkv", "C:\\videos\\video.xx.mp4"}, // no lang code/lang code invalid xx should remain
	},
}

func TestGenerateCaptionCandidates(t *testing.T) {
	for _, c := range testCases {
		assert.ElementsMatch(t, c.expectedCandidates, GenerateCaptionCandidates(c.captionPath, testExts))
	}
}

func TestGetCaptionsLangFromPath(t *testing.T) {
	for _, l := range testCases {
		assert.Equal(t, l.expectedLang, GetCaptionsLangFromPath(l.captionPath))
	}
}
