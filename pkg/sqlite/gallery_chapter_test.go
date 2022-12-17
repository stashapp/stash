//go:build integration
// +build integration

package sqlite_test

import (
	"context"
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sqlite"
	"github.com/stretchr/testify/assert"
)

func TestMarkerFindByGalleryID(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		mqb := sqlite.GalleryChapterReaderWriter

		galleryID := galleryIDs[galleryIdxWithChapters]
		chapters, err := mqb.FindByGalleryID(ctx, galleryID)

		if err != nil {
			t.Errorf("Error finding chapters: %s", err.Error())
		}

		assert.Greater(t, len(chapters), 0)
		for _, chapter := range chapters {
			assert.Equal(t, galleryIDs[galleryIdxWithChapters], int(chapter.GalleryID.Int64))
		}

		chapters, err = mqb.FindByGalleryID(ctx, 0)

		if err != nil {
			t.Errorf("Error finding chapter: %s", err.Error())
		}

		assert.Len(t, chapters, 0)

		return nil
	})
}

func queryChapters(ctx context.Context, t *testing.T, sqb models.GalleryChapterReader, markerFilter *models.GalleryChapterFilterType, findFilter *models.FindFilterType) []*models.GalleryChapter {
	t.Helper()
	result, _, err := sqb.Query(ctx, chapterFilter, findFilter)
	if err != nil {
		t.Errorf("Error querying chapters: %v", err)
	}

	return result
}

// TODO Update
// TODO Destroy
// TODO Find
// TODO Query
