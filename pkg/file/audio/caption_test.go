// TODO(audio): Can this file be deleted if we utilize audioCaptions?
package audio

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
		captionPath:    "/stash/audio.vtt",
		expectedLang:   LangUnknown,
		expectedResult: "/stash/audio.",
	},
	{
		captionPath:    "/stash/audio.en.vtt",
		expectedLang:   "en",
		expectedResult: "/stash/audio.", // lang code valid, remove en part
	},
	{
		captionPath:    "/stash/audio.test.srt",
		expectedLang:   LangUnknown,
		expectedResult: "/stash/audio.test.", // no lang code/lang code invalid test should remain
	},
	{
		captionPath:    "C:\\audios\\audio.fr.srt",
		expectedLang:   "fr",
		expectedResult: "C:\\audios\\audio.",
	},
	{
		captionPath:    "C:\\audios\\audio.xx.srt",
		expectedLang:   LangUnknown,
		expectedResult: "C:\\audios\\audio.xx.", // no lang code/lang code invalid xx should remain
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
