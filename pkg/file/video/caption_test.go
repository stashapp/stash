package video

import (
	"testing"

	"github.com/stashapp/stash/pkg/models"
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

func TestMergeGeneratedCaptionsPreservesExternalCaption(t *testing.T) {
	external := &models.VideoCaption{
		LanguageCode: "en",
		Filename:     "video.en.srt",
		CaptionType:  "srt",
	}
	generated := &models.VideoCaption{
		LanguageCode: "en",
		Filename:     "123.en.srt",
		CaptionType:  "srt",
		Generated:    true,
	}

	assert.Equal(t, []*models.VideoCaption{external}, MergeGeneratedCaptions([]*models.VideoCaption{external}, []*models.VideoCaption{generated}))
}

func TestMergeGeneratedCaptionsReplacesGeneratedCaption(t *testing.T) {
	oldGenerated := &models.VideoCaption{
		LanguageCode: "en",
		Filename:     "123.en.srt",
		CaptionType:  "srt",
		Generated:    true,
	}
	newGenerated := &models.VideoCaption{
		LanguageCode: "en",
		Filename:     "456.en.srt",
		CaptionType:  "srt",
		Generated:    true,
	}

	assert.Equal(t, []*models.VideoCaption{newGenerated}, MergeGeneratedCaptions([]*models.VideoCaption{oldGenerated}, []*models.VideoCaption{newGenerated}))
}
