//go:build integration
// +build integration

package sqlite_test

import (
	"context"
	"testing"

	"github.com/stashapp/stash/pkg/sqlite"
	"github.com/stretchr/testify/assert"
)

func TestChapterFindByGalleryID(t *testing.T) {
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

// TODO Update
// TODO Destroy
// TODO Find
