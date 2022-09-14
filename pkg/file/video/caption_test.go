package video

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
