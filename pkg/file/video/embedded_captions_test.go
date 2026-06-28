package video

import (
	"testing"

	"github.com/stashapp/stash/pkg/ffmpeg"
)

func subStream(index int, codec, lang string, def, forced int) *ffmpeg.FFProbeStream {
	s := &ffmpeg.FFProbeStream{CodecType: "subtitle", CodecName: codec, Index: index}
	s.Tags.Language = lang
	s.Disposition.Default = def
	s.Disposition.Forced = forced
	return s
}

func TestNormalizeSubtitleLang(t *testing.T) {
	cases := map[string]string{
		"eng":     "eng",
		"en":      "en",
		"ENG":     "eng",
		" fr ":    "fr",
		"":        LangUnknown,
		"und":     LangUnknown,
		"unknown": LangUnknown,
		"zxx":     LangUnknown,
		"engx":    LangUnknown, // too long / invalid
		"xx":      LangUnknown, // not a real language
	}
	for in, want := range cases {
		if got := normalizeSubtitleLang(in); got != want {
			t.Errorf("normalizeSubtitleLang(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestSelectEmbeddedSubtitles(t *testing.T) {
	streams := []*ffmpeg.FFProbeStream{
		subStream(2, "subrip", "eng", 0, 0),
		subStream(3, "ass", "eng", 1, 0), // default eng -> wins over index 2
		subStream(4, "subrip", "spa", 0, 0),
		subStream(5, "hdmv_pgs_subtitle", "fre", 0, 0), // image-based -> skipped
		subStream(6, "mov_text", "", 0, 0),             // unknown language
	}

	got := selectEmbeddedSubtitles(streams)

	want := []EmbeddedSubtitle{
		{Index: 3, LanguageCode: "eng"}, // default eng track preferred
		{Index: 4, LanguageCode: "spa"},
		{Index: 6, LanguageCode: LangUnknown}, // untagged stream
	}

	if len(got) != len(want) {
		t.Fatalf("got %d subtitles %+v, want %d", len(got), got, len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("subtitle[%d] = %+v, want %+v", i, got[i], want[i])
		}
	}
}

func TestSelectEmbeddedSubtitlesAllImageBased(t *testing.T) {
	streams := []*ffmpeg.FFProbeStream{
		subStream(2, "hdmv_pgs_subtitle", "eng", 1, 0),
		subStream(3, "dvd_subtitle", "spa", 0, 0),
	}
	if got := selectEmbeddedSubtitles(streams); len(got) != 0 {
		t.Errorf("expected no extractable subtitles, got %+v", got)
	}
}
