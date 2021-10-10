package api

import (
	"golang.org/x/text/collate"
	"golang.org/x/text/language"
)

// matcher defines a matcher for the languages we support
var matcher = language.NewMatcher([]language.Tag{
	language.MustParse("en-US"), // The first language is used as fallback.
	language.MustParse("de-DE"),
	language.MustParse("en-GB"),
	language.MustParse("en-US"),
	language.MustParse("pt-BR"),
	language.MustParse("sv-SE"),
	language.MustParse("zh-CN"),
	language.MustParse("zh-TW"),
})

// newCollator parses a locale into a collator
// Go through the available matches and return a valid match, in practice the first is a fallback
// If passed a nil-locale string, return nil
func newCollator(locale *string) *collate.Collator {
	if locale == nil {
		return nil
	}

	tag, _ := language.MatchStrings(matcher, *locale)
	return collate.New(tag)
}
