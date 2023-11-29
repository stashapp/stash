//go:build integration
// +build integration

package sqlite_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChapterFindByGalleryID(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		mqb := db.GalleryChapter

		galleryID := galleryIDs[galleryIdxWithChapters]
		chapters, err := mqb.FindByGalleryID(ctx, galleryID)

		if err != nil {
			t.Errorf("Error finding chapters: %s", err.Error())
		}

		assert.Greater(t, len(chapters), 0)
		for _, chapter := range chapters {
			assert.Equal(t, galleryIDs[galleryIdxWithChapters], chapter.GalleryID)
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
