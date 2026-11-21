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

type stubCaptionUpdater struct {
	captions []*models.VideoCaption
	updated  int
}

func (s *stubCaptionUpdater) GetCaptions(ctx context.Context, fileID models.FileID) ([]*models.VideoCaption, error) {
	return s.captions, nil
}

func (s *stubCaptionUpdater) UpdateCaptions(ctx context.Context, fileID models.FileID, captions []*models.VideoCaption) error {
	s.captions = captions
	s.updated++
	return nil
}

func TestNormalizeSubtitleLanguage(t *testing.T) {
	assert.Equal(t, "en", NormalizeSubtitleLanguage("eng"))
	assert.Equal(t, "en", NormalizeSubtitleLanguage("en"))
	assert.Equal(t, "ja", NormalizeSubtitleLanguage("jpn"))
	assert.Equal(t, "fr", NormalizeSubtitleLanguage("fre"))
	assert.Equal(t, LangUnknown, NormalizeSubtitleLanguage("und"))
	assert.Equal(t, LangUnknown, NormalizeSubtitleLanguage(""))
	assert.Equal(t, LangUnknown, NormalizeSubtitleLanguage("not-a-lang"))
}

func TestEmbeddedCaptionsFromStreams(t *testing.T) {
	s1 := &ffmpeg.FFProbeStream{Index: 2, CodecType: "subtitle", CodecName: "subrip"}
	s1.Tags.Language = "eng"
	s2 := &ffmpeg.FFProbeStream{Index: 3, CodecType: "subtitle", CodecName: "ass"}
	s2.Tags.Language = "eng"
	s3 := &ffmpeg.FFProbeStream{Index: 4, CodecType: "subtitle", CodecName: "hdmv_pgs_subtitle"}
	s3.Tags.Language = "jpn"
	s4 := &ffmpeg.FFProbeStream{Index: 5, CodecType: "subtitle", CodecName: "mov_text"}
	s4.Tags.Language = ""

	captions := EmbeddedCaptionsFromStreams([]*ffmpeg.FFProbeStream{s1, s2, s3, s4})
	require.Len(t, captions, 3)

	assert.Equal(t, "en", captions[0].LanguageCode)
	assert.Equal(t, "embedded:2", captions[0].Filename)
	assert.Equal(t, "embedded", captions[0].CaptionType)
	assert.True(t, captions[0].IsEmbedded())
	idx, ok := captions[0].EmbeddedStreamIndex()
	assert.True(t, ok)
	assert.Equal(t, 2, idx)

	assert.Equal(t, "en", captions[1].LanguageCode)
	assert.Equal(t, "embedded:3", captions[1].Filename)
	assert.Equal(t, "embedded-2", captions[1].CaptionType)

	assert.Equal(t, LangUnknown, captions[2].LanguageCode)
	assert.Equal(t, "embedded:5", captions[2].Filename)
	assert.Equal(t, "embedded", captions[2].CaptionType)
}

func TestCleanCaptionsPreservesEmbedded(t *testing.T) {
	dir := t.TempDir()
	videoPath := filepath.Join(dir, "movie.mkv")
	existingSidecar := filepath.Join(dir, "movie.en.vtt")
	require.NoError(t, os.WriteFile(existingSidecar, []byte("WEBVTT\n"), 0644))

	updater := &stubCaptionUpdater{
		captions: []*models.VideoCaption{
			{LanguageCode: "en", Filename: "movie.en.vtt", CaptionType: "vtt"},
			{LanguageCode: "fr", Filename: "movie.fr.srt", CaptionType: "srt"}, // missing sidecar
			models.NewEmbeddedVideoCaption("ja", 2, EmbeddedCaptionType),
		},
	}

	vf := &models.VideoFile{BaseFile: &models.BaseFile{ID: 1, Path: videoPath}}
	require.NoError(t, CleanCaptions(context.Background(), vf, nil, updater))
	assert.Equal(t, 1, updater.updated)
	require.Len(t, updater.captions, 2)
	assert.Equal(t, "movie.en.vtt", updater.captions[0].Filename)
	assert.Equal(t, "embedded:2", updater.captions[1].Filename)
}

func TestSyncEmbeddedCaptionsMergesWithSidecars(t *testing.T) {
	updater := &stubCaptionUpdater{
		captions: []*models.VideoCaption{
			{LanguageCode: "en", Filename: "movie.en.vtt", CaptionType: "vtt"},
		},
	}

	vf := &models.VideoFile{
		BaseFile: &models.BaseFile{ID: 1, Path: "/videos/movie.mkv"},
		EmbeddedCaptions: []*models.VideoCaption{
			models.NewEmbeddedVideoCaption("en", 2, EmbeddedCaptionType),
			models.NewEmbeddedVideoCaption("ja", 3, EmbeddedCaptionType),
		},
	}

	require.NoError(t, SyncEmbeddedCaptions(context.Background(), vf, nil, updater))
	assert.Equal(t, 1, updater.updated)
	require.Len(t, updater.captions, 3)

	// Calling again with identical embedded captions should be a no-op.
	require.NoError(t, SyncEmbeddedCaptions(context.Background(), vf, nil, updater))
	assert.Equal(t, 1, updater.updated)
}
